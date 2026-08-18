package domain

import (
	"context"
	triptype "meridian/services/trip-service/pkg/types"
	pb "meridian/shared/proto/trip"
	"meridian/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

func (t *TripModel) ToProto() *pb.Trip {
	return &pb.Trip{
		Id:           t.ID.Hex(),
		UserID:       t.UserID,
		SelectedFare: t.RideFare.ToProto(),
		Status:       t.Status,
		Driver:       t.Driver,
		Route:        t.RideFare.Route.ToProto(),
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, ridefare *RideFareModel) error
	GetRideFareByID(context context.Context, id string) (*RideFareModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*triptype.OsrmApiResponse, error)
	GenerateFare(ctx context.Context, UserID string, Fare []*RideFareModel, Route *triptype.OsrmApiResponse) ([]*RideFareModel, error)
	EstimateFare(osrm *triptype.OsrmApiResponse) []*RideFareModel
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
}
