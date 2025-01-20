package main

var targetReceipt = `{
	"retailer": "Target",
	"purchaseDate": "2022-01-01",
	"purchaseTime": "13:01",
	"items": [
		{
			"shortDescription": "Mountain Dew 12PK",
			"price": "6.49"
		},{
			"shortDescription": "Emils Cheese Pizza",
			"price": "12.25"
		},{
			"shortDescription": "Knorr Creamy Chicken",
			"price": "1.26"
		},{
			"shortDescription": "Doritos Nacho Cheese",
			"price": "3.35"
		},{
			"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
			"price": "12.00"
		}
	],
	"total": "35.35"
}`

var mAndMCornerMarketReceipt = `{
	"retailer": "M&M Corner Market",
	"purchaseDate": "2022-03-20",
	"purchaseTime": "14:33",
	"items": [
		{
			"shortDescription": "Gatorade",
			"price": "2.25"
		},{
			"shortDescription": "Gatorade",
			"price": "2.25"
		},{
			"shortDescription": "Gatorade",
			"price": "2.25"
		},{
			"shortDescription": "Gatorade",
			"price": "2.25"
		}
	],
	"total": "9.00"
}`

var validReceipt = `
{
	"retailer": "Walgreens",
	"purchaseDate": "2022-01-02",
	"purchaseTime": "08:13",
	"total": "2.65",
	"items": [
			{"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
			{"shortDescription": "Dasani", "price": "1.40"}
	]
}`
