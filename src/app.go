package main

import (
	"log"
	"net/http"

	oapi "github.com/drew-loukusa/drew-fetch-receipts-processor/server/openapi"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	log.Printf("Server Started")

	ReceiptsProcessorService := NewReceiptsService()
	ReceiptsProcessorController := oapi.NewDefaultAPIController(ReceiptsProcessorService)
	router := oapi.NewRouter(ReceiptsProcessorController)

	a.Router = router

	router.Use(loggingMiddleware)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
