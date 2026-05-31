package main

import (
	"context"
	"log"
	httppkg "meridian/services/trip-service/internal/infrastructure/http"
	"meridian/services/trip-service/internal/infrastructure/repository"
	"meridian/services/trip-service/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server listening on %s", addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)

	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the server gracefully: %v", err)
			server.Close()
		}

	}
}
