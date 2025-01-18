package main

import (
	"log"
	"net/http"

	receiptsProcessorServer "github.com/drew-loukusa/drew-fetch-receipts-processor/server/openapi"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func fooBar(w http.ResponseWriter, r *http.Request) {
	log.Printf("AHahahahahahahahahahah")
}

func main() {

	log.Printf("Server Started")

	ReceiptsProcessorService := receiptsProcessorServer.NewMyService()
	ReceiptsProcessorController := receiptsProcessorServer.NewDefaultAPIController(ReceiptsProcessorService)

	router := receiptsProcessorServer.NewRouter(ReceiptsProcessorController)

	router.HandleFunc("/fooBar", fooBar).Methods("GET")

	router.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
