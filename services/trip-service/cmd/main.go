package main

import (
	"context"
	"log"
	"meridian/services/trip-service/internal/infrastructure/grpc"
	"meridian/services/trip-service/internal/infrastructure/repository"
	"meridian/services/trip-service/internal/service"
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

	grpcServer := grpcServer.NewServer()
	log.Printf("the grpc server is starting at %v", lis.Addr().String())

	grpc.NewGrpcHandler(grpcServer, svc)

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
