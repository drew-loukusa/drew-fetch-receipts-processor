// main_test.go

package main

import (
	"bytes"
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

var a App
var client = Client{app: &a}

func TestMain(m *testing.M) {
	a.Initialize()
	code := m.Run()
	os.Exit(code)
}

func checkResponseCode(t *testing.T, expectedResponseCode int, actualResponse *httptest.ResponseRecorder) {
	if expectedResponseCode != actualResponse.Code {
		t.Errorf("Expected response code %d. Got %d\n", expectedResponseCode, actualResponse.Code)
		t.Error(actualResponse.Body)
	}
}

// ===========================================================================
// 												Integration-ish Testing
// ===========================================================================

func TestProcessReceipt(t *testing.T) {
	bodyJsonStr := []byte(validReceipt)
	req, _ := http.NewRequest("POST", PROCESS_RECEIPTS_URL, bytes.NewBuffer(bodyJsonStr))
	response := client.executeRequest(req)
	checkResponseCode(t, http.StatusOK, response)
}

// Parameterized tests for receipt processing
func TestProcessManyReceipts(t *testing.T) {
	testCases := []struct {
		receipt        string
		expectedPoints int
	}{
		{targetReceipt, 28},
		{mAndMCornerMarketReceipt, 109},
	}
	for _, tc := range testCases {
		receiptId := client.processReceipt(tc.receipt)
		points := client.getReceiptPoints(receiptId)
		if points != tc.expectedPoints {
			t.Errorf("Expected %d points but got %d points", tc.expectedPoints, points)
		}
	}
}

func TestGetReceipt(t *testing.T) {
	// Given
	receiptId := client.processReceipt(validReceipt)
	// When
	response := client.getReceiptPointsHttp(receiptId)
	// Then
	checkResponseCode(t, http.StatusOK, response)
}

func TestGetNonExistantReceipt(t *testing.T) {
	// Given
	receiptId := uuid.NewString()
	// When
	response := client.getReceiptPointsHttp(receiptId)
	// Then
	checkResponseCode(t, http.StatusNotFound, response)
}

// ===========================================================================
// 												Input Validation Testing
// ===========================================================================

func TestProcessReceipt_MissingRequiredFields(t *testing.T) {
	// Missing purchaseDate and purchaseTime
	bodyJsonStr := `{
    "retailer": "Strosin Inc",
    "total": "10.40",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
	}`
	// When
	response := client.processReceiptHttp(bodyJsonStr)
	// Then
	checkResponseCode(t, http.StatusUnprocessableEntity, response)
}

func TestProcessReceipt_ItemMissingRequiredFields(t *testing.T) {
	// An item is missing price
	bodyJsonStr := `{
    "retailer": "Strosin Inc",
    "total": "10.40",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani"}
    ]
	}`
	// When
	response := client.processReceiptHttp(bodyJsonStr)
	// Then
	checkResponseCode(t, http.StatusUnprocessableEntity, response)
}

func TestProcessReceipt_IncorrectlyFormattedTotal(t *testing.T) {
	// Given total has too many digits after the '.'
	bodyJsonStr := `{
    "retailer": "Strosin Inc",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "08:13",
    "total": "10.409",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
	}`
	// When
	response := client.processReceiptHttp(bodyJsonStr)
	// Then
	checkResponseCode(t, http.StatusBadRequest, response)
}

func TestProcessReceipt_IncorrectlyFormattedDate(t *testing.T) {
	// Given total has too many digits after the '.'
	bodyJsonStr := `{
    "retailer": "Strosin Inc",
    "purchaseDate": "022-1-02",
    "purchaseTime": "08:13",
    "total": "10.40",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
	}`
	// When
	response := client.processReceiptHttp(bodyJsonStr)
	// Then
	checkResponseCode(t, http.StatusBadRequest, response)
}

func TestProcessReceipt_IncorrectlyFormattedTime(t *testing.T) {
	// Given total has too many digits after the '.'
	bodyJsonStr := `{
    "retailer": "Strosin Inc",
    "purchaseDate": "2025-01-18",
    "purchaseTime": "08:131",
    "total": "10.40",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
	}`
	// When
	response := client.processReceiptHttp(bodyJsonStr)
	// Then
	checkResponseCode(t, http.StatusBadRequest, response)
}

func TestProcessReceipt_IncorrectlyFormattedRetailer(t *testing.T) {
	// Given total has too many digits after the '.'
	bodyJsonStr := `{
    "retailer": "Strosin Inc !$){()}",
    "purchaseDate": "2025-01-18",
    "purchaseTime": "08:11",
    "total": "10.40",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
	}`
	// When
	response := client.processReceiptHttp(bodyJsonStr)
	// Then
	checkResponseCode(t, http.StatusBadRequest, response)
}

func TestProcessReceipt_IncorrectlyFormattedItemPrice(t *testing.T) {
	// Given total has too many digits after the '.'
	bodyJsonStr := `{
    "retailer": "Strosin Inc",
    "purchaseDate": "2025-01-18",
    "purchaseTime": "08:11",
    "total": "10.40",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.20005"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
	}`
	// When
	response := client.processReceiptHttp(bodyJsonStr)
	// Then
	checkResponseCode(t, http.StatusBadRequest, response)
}

func TestProcessReceipt_IncorrectlyFormattedItemDescription(t *testing.T) {
	// Given total has too many digits after the '.'
	bodyJsonStr := `{
    "retailer": "Strosin Inc",
    "purchaseDate": "2025-01-18",
    "purchaseTime": "08:11",
    "total": "10.40",
    "items": [
        {"shortDescription": "Pepsi - 12-oz !@#!@#!#!@#!@#!#", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
	}`
	// When
	response := client.processReceiptHttp(bodyJsonStr)
	// Then
	checkResponseCode(t, http.StatusBadRequest, response)
}
