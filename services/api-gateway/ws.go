package main

import (
	"log"
	"meridian/services/api-gateway/grpc_clients"
	"meridian/shared/contracts"
	"meridian/shared/proto/driver"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error in establishing connection :", err)
		return
	}
	defer conn.Close()
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Print("No userID found :", err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("error parsing message from websocker", err)
		}
		log.Print(string(message))
	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error in establishing connection :", err)
		return
	}
	defer conn.Close()
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Print("No userID found :")
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Print("No package slug found :", err)
	}

	ctx := r.Context()
	driverService, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		driverService.Client.UnregisterDriver(ctx, &driver.RegisterDriverRequest{
			DriverID:    userID,
			PackageSlug: packageSlug,
		})
		driverService.Close()
		log.Println("Driver unregistered: ", userID)
	}()

	driverData, err := driverService.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverID:    userID,
		PackageSlug: packageSlug,
	})
	if err != nil {
		log.Printf("Error Registering Driver : %v", err)
		return
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: driverData.Driver,
	}

	if err = conn.WriteJSON(msg); err != nil {
		log.Print("error sending message", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("error parsing message from websocker", err)
			break
		}
		log.Print(message)
	}
}
