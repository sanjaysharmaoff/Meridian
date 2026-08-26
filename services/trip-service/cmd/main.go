package main

import (
	"context"
	"log"
	"meridian/services/trip-service/internal/infrastructure/events"
	"meridian/services/trip-service/internal/infrastructure/grpc"
	"meridian/services/trip-service/internal/infrastructure/repository"
	"meridian/services/trip-service/internal/service"
	"meridian/shared/env"
	"meridian/shared/messaging"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcServer "google.golang.org/grpc"
)

var GrpcAddr = ":9093"

func main() {
	inmem := repository.NewInmemRepository()
	svc := service.NewService(inmem)
	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
		<-sigchan
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to Listen %v", err)
	}

	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	publisher := events.NewTripEventPublisher(rabbitmq)
	log.Println("Starting RabbitMQ connection")
	driverConsumer := events.NewDriverConsumer(rabbitmq, svc)
	go driverConsumer.Listen()

	paymentConsumer := events.NewPaymentConsumer(rabbitmq, svc)
	go paymentConsumer.Listen()

	grpcServer := grpcServer.NewServer()
	log.Printf("the grpc server is starting at %v", lis.Addr().String())

	grpc.NewGrpcHandler(grpcServer, svc, publisher)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("grpc server has encountered an error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Print("shutting down the server ....")
	grpcServer.GracefulStop()
}
