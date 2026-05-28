package main

import (
	"bytes"
	"encoding/json"
	"log"
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

	jsonbody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonbody)

	resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		log.Print(err)
		return
	}

	defer resp.Body.Close()

	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, "failed to parse JSON data from trip service", http.StatusBadRequest)
		return
	}

	response := contracts.APIResponse{
		Data:  respBody,
		Error: nil,
	}
	writeJSON(w, http.StatusCreated, response)
}
