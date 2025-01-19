package main

import (
	"log"
	"net/http"

	"github.com/drew-loukusa/drew-fetch-receipts-processor/server/openapi"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {

	log.Printf("Server Started")

	ReceiptsProcessorService := NewReceiptsService()
	ReceiptsProcessorController := openapi.NewDefaultAPIController(ReceiptsProcessorService)
	router := openapi.NewRouter(ReceiptsProcessorController)

	router.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
