package grpc_clients

import (
	pb "meridian/shared/proto/trip"
	"meridian/shared/tracing"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tripServiceClient struct {
	Client pb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient() (*tripServiceClient, error) {
	tripServiceURL := os.Getenv("TRIP_SERVICE_URL")
	if tripServiceURL == "" {
		tripServiceURL = "trip-service:9093"
	}

	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(tripServiceURL, dialOptions...)
	// conn, err := grpc.NewClient(tripServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	Client := pb.NewTripServiceClient(conn)

	return &tripServiceClient{
		Client: Client,
		conn:   conn,
	}, nil
}

func (s *tripServiceClient) Close() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return
		}
	}

}
