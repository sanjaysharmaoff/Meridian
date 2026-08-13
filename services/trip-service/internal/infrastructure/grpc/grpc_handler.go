package grpc

import (
	"context"
	"meridian/services/trip-service/internal/domain"
	pb "meridian/shared/proto/trip"
	"meridian/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGrpcHandler(server *grpc.Server, service domain.TripService) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *grpcHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	droppoint := req.GetEndLocation()
	pick := types.Coordinate{
		Longitude: pickup.Longitude,
		Latitude:  pickup.Latitude,
	}
	drop := types.Coordinate{
		Longitude: droppoint.Longitude,
		Latitude:  droppoint.Latitude,
	}
	route, err := h.service.GetRoute(ctx, &pick, &drop)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	estimatedFares := h.service.EstimateFare(route)
	fares, err := h.service.GenerateFare(ctx, req.GetUserID(), estimatedFares, route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "problem generating the fare: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}, nil

}

func (h *grpcHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareid := req.GetRideFareID()
	userid := req.GetUserID()
	fare, err := h.service.GetAndValidateFare(ctx, fareid, userid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fare is invalid error : %v", err)
	}
	trip, err := h.service.CreateTrip(ctx, fare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trip : %v", err)
	}
	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
