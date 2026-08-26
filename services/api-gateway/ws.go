package main

import (
	"encoding/json"
	"log"
	"meridian/services/api-gateway/grpc_clients"
	"meridian/shared/contracts"
	"meridian/shared/messaging"
	"meridian/shared/proto/driver"
	"net/http"
)

var (
	connManager = messaging.NewConnectionManager()
)

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request,
	rb *messaging.RabbitMQ) {
	conn, err := connManager.Upgrade(w, r)
	if err != nil {
		log.Print("error in establishing connection :", err)
		return
	}
	defer conn.Close()
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Print("No userID found :", err)
	}
	connManager.Add(userID, conn)
	defer connManager.Remove(userID)
	queues := []string{
		messaging.NotifyDriverNoDriversFoundQueue,
		messaging.NotifyDriverAssignQueue,
		messaging.NotifyPaymentSessionCreatedQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(rb, connManager, q)

		if err := consumer.Start(); err != nil {
			log.Printf("Failed to start consumer for queue: %s: err: %v", q, err)
		}
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("error parsing message from websocker", err)
			break
		}
		log.Print(string(message))
	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request, rb *messaging.RabbitMQ) {
	conn, err := connManager.Upgrade(w, r)
	if err != nil {
		log.Print("error in establishing connection :", err)
		return
	}
	defer conn.Close()
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Print("No userID found :")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Print("No package slug found :", err)
	}
	connManager.Add(userID, conn)

	ctx := r.Context()
	driverService, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		connManager.Remove(userID)
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

	if err := connManager.SendMessage(userID, contracts.WSMessage{
		Type: contracts.DriverCmdRegister,
		Data: driverData.Driver,
	}); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	queues := []string{
		messaging.DriverCmdTripRequestQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(rb, connManager, q)

		if err := consumer.Start(); err != nil {
			log.Printf("Failed to start consumer for queue: %s: err: %v", q, err)
		}
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Print("error parsing message from websocker", err)
			break
		}
		log.Print(message)

		type driverMessage struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}

		var driverMsg driverMessage
		if err := json.Unmarshal(message, &driverMsg); err != nil {
			log.Printf("Error unmarshaling driver message: %v", err)
			continue
		}

		// Handle the different message type
		switch driverMsg.Type {
		case contracts.DriverCmdLocation:
			// Handle driver location update in the future
			continue
		case contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline:
			// Forward the message to RabbitMQ
			if err := rb.PublishMessage(ctx, driverMsg.Type, contracts.AmqpMessage{
				OwnerID: userID,
				Data:    driverMsg.Data,
			}); err != nil {
				log.Printf("Error publishing message to RabbitMQ: %v", err)
			}
		default:
			log.Printf("Unknown message type: %s", driverMsg.Type)
		}
	}

}
