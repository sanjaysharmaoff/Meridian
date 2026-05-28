package main

import (
	"log"
	httppkg "meridian/services/trip-service/internal/infrastructure/http"
	"meridian/services/trip-service/internal/infrastructure/repository"
	"meridian/services/trip-service/internal/service"
	"net/http"
)

func main() {
	inmem := repository.NewInmemRepository()
	svc := service.NewService(inmem)

	addr := ":8083"
	mux := http.NewServeMux()
	h := &httppkg.Httphandler{
		Service: svc,
	}
	mux.HandleFunc("POST /preview", h.HandleTripPreview)
	server := http.Server{
		Handler: mux,
		Addr:    addr,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Print("an error has occured in server ", err)
	}
}
