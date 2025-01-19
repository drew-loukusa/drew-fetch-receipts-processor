// main_test.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize()
	code := m.Run()
	os.Exit(code)
}

var validReciept = `{
	"retailer": "Strosin Inc",
	"purchaseDate": "2022-01-02",
	"purchaseTime": "08:13",
	"total": "10.40",
	"items": [
			{"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
			{"shortDescription": "Dasani", "price": "1.40"}
	]
}`
var PROCESS_RECEIPTS_URL = "/receipts/process"

func mkGetReceiptUrl(receiptId string) (url string) {
	return fmt.Sprintf("/receipts/%s/points", receiptId)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expectedResponseCode int, actualResponse *httptest.ResponseRecorder) {
	if expectedResponseCode != actualResponse.Code {
		t.Errorf("Expected response code %d. Got %d\n", expectedResponseCode, actualResponse.Code)
		t.Error(actualResponse.Body)
	}
}

func TestProcessReceipt(t *testing.T) {
	bodyJsonStr := []byte(validReciept)
	req, _ := http.NewRequest("POST", PROCESS_RECEIPTS_URL, bytes.NewBuffer(bodyJsonStr))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response)
}

func TestProcessReceipt_IncorrectlyFormattedTotal(t *testing.T) {
	// Total has too many digits after the '.'
	bodyJsonStr := []byte(`{
    "retailer": "Strosin Inc",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "08:13",
    "total": "10.409",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
		}`)
	req, _ := http.NewRequest("POST", PROCESS_RECEIPTS_URL, bytes.NewBuffer(bodyJsonStr))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response)
}

func makeReceipt() (id string) {
	bodyJsonStr := []byte(validReciept)
	req, _ := http.NewRequest("POST", PROCESS_RECEIPTS_URL, bytes.NewBuffer(bodyJsonStr))
	response := executeRequest(req)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	return m["id"]
}

func TestGetReceipt(t *testing.T) {
	receiptId := makeReceipt()
	req, _ := http.NewRequest("GET", mkGetReceiptUrl(receiptId), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response)
}

func TestGetNonExistantReceipt(t *testing.T) {
	receiptId := uuid.NewString()
	req, _ := http.NewRequest("GET", mkGetReceiptUrl(receiptId), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response)
}
