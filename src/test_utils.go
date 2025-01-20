package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

var PROCESS_RECEIPTS_URL = "/receipts/process"

func mkGetReceiptUrl(receiptId string) (url string) {
	return fmt.Sprintf("/receipts/%s/points", receiptId)
}

type Client struct {
	app *App
}

func (c *Client) executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	c.app.Router.ServeHTTP(rr, req)

	return rr
}

func (c *Client) getReceiptPointsHttp(receiptId string) (response *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("GET", mkGetReceiptUrl(receiptId), nil)
	return c.executeRequest(req)
}

func (c *Client) getReceiptPoints(receiptId string) (points int) {
	response := c.getReceiptPointsHttp(receiptId)
	var m map[string]int
	json.Unmarshal(response.Body.Bytes(), &m)
	points, ok := m["points"]
	if !ok {
		panic(fmt.Sprintf("'points' was not found in response body: %s", response.Body))
	}
	return points
}

func (c *Client) processReceiptHttp(receipt string) (response *httptest.ResponseRecorder) {
	bodyJsonStr := []byte(receipt)
	req, _ := http.NewRequest("POST", PROCESS_RECEIPTS_URL, bytes.NewBuffer(bodyJsonStr))
	return c.executeRequest(req)
}

func (c *Client) processReceipt(receipt string) (receiptId string) {
	response := c.processReceiptHttp(receipt)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	id, ok := m["id"]
	if !ok {
		panic(fmt.Sprintf("id was not found in response body: %s", response.Body))
	}
	return id
}
