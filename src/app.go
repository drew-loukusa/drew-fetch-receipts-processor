package main

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"

	oapi "github.com/drew-loukusa/drew-fetch-receipts-processor/server/openapi"
)

type ReceiptsService struct {
	receiptsRepo map[string]int64
}

func NewReceiptsService() *ReceiptsService {
	return &ReceiptsService{receiptsRepo: make(map[string]int64)}
}

func (s *ReceiptsService) ProcessReceipt(ctx context.Context, receipt oapi.Receipt) (oapi.ImplResponse, error) {
	// thing, err := s.store.GetThing(uuid)
	// if err == helpers.ErrNotFound {
	// 	return oapi.Response(http.StatusNotFound, nil), nil
	// }
	// if err != nil {
	// 	return oapi.Response(http.StatusInternalServerError, nil), err
	// }
	log.Println("Got request to process receipt")
	log.Printf("%#v\n", receipt)
	receiptId := uuid.NewString()
	response := oapi.ProcessReceipt200Response{Id: receiptId}

	s.receiptsRepo[receiptId] = 10

	return oapi.Response(http.StatusOK, response), nil
}

func (s *ReceiptsService) GetReceiptPoints(ctx context.Context, id string) (oapi.ImplResponse, error) {
	log.Printf("Got request to get receipt points for id %s", id)
	points, ok := s.receiptsRepo[id]

	if !ok {
		return oapi.Response(http.StatusNotFound, nil), nil
	}

	log.Printf("Points for receipt are %d", points)
	response := oapi.GetReceiptPoints200Response{Points: points}
	return oapi.Response(http.StatusOK, response), nil
}
