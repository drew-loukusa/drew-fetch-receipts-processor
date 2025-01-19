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

func TotalIsMultipleOf(multipleOf float64, total string) bool {
	totalFloat, ok := strconv.ParseFloat(total, 64)
	if ok != nil {
		panic(fmt.Sprintf("Could not parse float from %s", total))
	}
	return totalFloat/float64(multipleOf) == 0
}

func CountAlphaNumericChars(src string) int {
	alphanumericRe := regexp.MustCompile("[a-zA-Z0-9]")
	matches := alphanumericRe.FindAllStringIndex(src, -1)
	return len(matches)
}

func DateIsOdd(src string) bool {
	oddDateRe := regexp.MustCompile(`[13579]$`)
	return oddDateRe.MatchString(src)
}

func MustParseTime(layout, value string) time.Time {
	// valueTrimmed := strings.TrimLeft(value, "0")
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

	if TotalIsMultipleOf(0.25, receipt.Total) {
		points += 25
	}

	// 5 points for every 2 items
	points += len(receipt.Items) / 2

	// For each item, check if trimmed desc len is multiple of 3
	// and do some math do the item price and add that to the total
	for _, item := range receipt.Items {
		descTrimmed := strings.TrimSpace(item.ShortDescription)
		if len(descTrimmed)/3 == 0 {
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
		return oapi.Response(http.StatusNotFound, nil), nil
	}

	log.Printf("Points for receipt are %d", points)
	response := oapi.GetReceiptPoints200Response{Points: points}
	return oapi.Response(http.StatusOK, response), nil
}
