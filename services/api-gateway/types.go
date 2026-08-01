package main

import (
	pb "meridian/shared/proto/trip"
	"meridian/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (t *previewTripRequest) ToProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: t.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  t.Pickup.Latitude,
			Longitude: t.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  t.Destination.Latitude,
			Longitude: t.Destination.Longitude,
		},
	}
}

type startTripRequest struct {
	RideFareID string `json:"rideFareID"`
	UserID     string `json:"userID"`
}

func (s *startTripRequest) ToProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		RideFareID: s.RideFareID,
		UserID:     s.UserID,
	}
}
