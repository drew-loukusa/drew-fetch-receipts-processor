package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	oapi "github.com/drew-loukusa/drew-fetch-receipts-processor/server/openapi"
	"github.com/google/uuid"
)

type ReceiptsService struct {
	ReceiptsRepo map[string]int64 // Just store receipts in memory as per instructions
}

func NewReceiptsService() *ReceiptsService {
	return &ReceiptsService{ReceiptsRepo: make(map[string]int64)}
}

// Returns true if total is a round number (no cents)
func TotalIsRound(total string) bool {
	roundTotal := regexp.MustCompile(`\d+\.00`)
	return roundTotal.MatchString(total)
}

// Should probably make this generic (work for any multiple), but this works for now
func TotalIsMultipleOf25(total string) bool {
	multipleOf25Re := regexp.MustCompile(`^\d+\.(00|25|50|75)`)
	return multipleOf25Re.MatchString(total)
}

func CountAlphaNumericChars(src string) int {
	alphanumericRe := regexp.MustCompile("[a-zA-Z0-9]")
	matches := alphanumericRe.FindAllStringIndex(src, -1)
	return len(matches)
}

func stringMatches(src, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(src)
}

var STRING_ENDS_WITH_ODD_PATTERN = `[13579]$` // Pattern for checking if string ends with odd number
func DateIsOdd(src string) bool {
	return stringMatches(src, STRING_ENDS_WITH_ODD_PATTERN)
}

func MustParseTime(layout, value string) time.Time {
	result, err := time.Parse(layout, value)

	if err != nil {
		panic(fmt.Sprintf("Failed to parse time from %s", value))
	}

	return result
}

// Returns true if targetTime between startTime and endTime
func TimeBetween(targetTime, startTime, endTime string) bool {
	layout := "15:04"
	target := MustParseTime(layout, targetTime)
	start := MustParseTime(layout, startTime)
	end := MustParseTime(layout, endTime)
	return target.After(start) && target.Before(end)
}

// Count how many points a receipt is worth
func CountPoints(receipt oapi.Receipt) int64 {
	points := 0

	points += CountAlphaNumericChars(receipt.Retailer)

	if TotalIsRound(receipt.Total) {
		points += 50
	}

	if TotalIsMultipleOf25(receipt.Total) {
		points += 25
	}

	// 5 points for every 2 items
	points += (len(receipt.Items) / 2) * 5

	// For each item, check if trimmed desc len is multiple of 3
	// and do some math do the item price and add that to the total
	for _, item := range receipt.Items {
		descTrimmed := strings.TrimSpace(item.ShortDescription)
		descLen := len(descTrimmed)
		descLenMultipleOf3Res := descLen % 3
		if descLenMultipleOf3Res == 0 {
			newItemPrice, ok := strconv.ParseFloat(item.Price, 64)
			if ok != nil {
				panic(fmt.Sprintf("Failed to parse float from %s", item.Price))
			}
			points += int(math.Ceil(newItemPrice * 0.2))
		}
	}

	if DateIsOdd(receipt.PurchaseDate) {
		points += 6
	}

	if TimeBetween(receipt.PurchaseTime, "14:00", "16:00") {
		points += 10
	}

	return int64(points)
}

func (s *ReceiptsService) ProcessReceipt(ctx context.Context, receipt oapi.Receipt) (oapi.ImplResponse, error) {
	log.Println("Got request to process receipt")
	log.Printf("%#v\n", receipt)
	receiptId := uuid.NewString()
	response := oapi.ProcessReceipt200Response{Id: receiptId}

	s.ReceiptsRepo[receiptId] = CountPoints(receipt)

	return oapi.Response(http.StatusOK, response), nil
}

func (s *ReceiptsService) GetReceiptPoints(ctx context.Context, id string) (oapi.ImplResponse, error) {
	log.Printf("Got request to get receipt points for id %s", id)
	points, ok := s.ReceiptsRepo[id]

	if !ok {
		return oapi.Response(http.StatusNotFound, "No receipt found for that ID."), nil
	}

	log.Printf("Points for receipt are %d", points)
	response := oapi.GetReceiptPoints200Response{Points: points}
	return oapi.Response(http.StatusOK, response), nil
}
