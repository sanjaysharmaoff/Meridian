package main

import (
	"log"
	"meridian/shared/contracts"
	"meridian/shared/util"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error in establishing connection :", err)
		return
	}
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Print("No userID found :", err)
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("error parsing message from websocker", err)
		}
		log.Print(string(message))
	}
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error in establishing connection :", err)
		return
	}
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Print("No userID found :")
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Print("No package slug found :", err)
	}

	carPlate := r.URL.Query().Get("carPlate")
	if carPlate == "" {
		carPlate = "abc123"
	}

	type Driver struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profile_picture"`
		CarPlate       string `json:"car_plate"`
		PackageSlug    string `json:"package_slug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			Id:             userID,
			Name:           "Sanjay",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       carPlate,
			PackageSlug:    packageSlug,
		},
	}

	if err = conn.WriteJSON(msg); err != nil {
		log.Print("error sending message", err)
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("error parsing message from websocker", err)
			break
		}
		log.Print(message)
	}
}
