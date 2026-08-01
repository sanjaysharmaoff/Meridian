package main

import (
	"encoding/json"
	"log"
	"meridian/services/api-gateway/grpc_clients"
	"meridian/shared/contracts"
	"net/http"
)

func welcomewala(w http.ResponseWriter, r *http.Request) {
	log.Print("welcome to my server bhai")
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if reqBody.UserID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()
	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.ToProto())

	response := contracts.APIResponse{
		Data:  tripPreview,
		Error: nil,
	}
	writeJSON(w, http.StatusCreated, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if reqBody.UserID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	tripStart, err := tripService.Client.CreateTrip(r.Context(), reqBody.ToProto())
	resp := contracts.APIResponse{
		Data:  tripStart,
		Error: nil,
	}
	writeJSON(w, http.StatusCreated, resp)
}
