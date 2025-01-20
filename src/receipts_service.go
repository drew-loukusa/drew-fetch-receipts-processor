package main

import (
	"context"
	"log"
	"net/http"

	oapi "github.com/drew-loukusa/drew-fetch-receipts-processor/server/openapi"
	"github.com/google/uuid"
)

type ReceiptsService struct {
	ReceiptsRepo map[string]int64 // Just store receipts in memory as per instructions
}

func NewReceiptsService() *ReceiptsService {
	return &ReceiptsService{ReceiptsRepo: make(map[string]int64)}
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
