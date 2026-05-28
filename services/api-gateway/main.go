package main

import (
	"log"
	"net/http"

	"meridian/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway")
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	mux.HandleFunc("/", welcomewala)
	mux.HandleFunc("POST /trip/preview", handleTripPreview)

	if err := server.ListenAndServe(); err != nil {
		log.Print("http server error", err)
	}
}
