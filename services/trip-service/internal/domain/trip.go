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

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, ridefare *RideFareModel) error
	GetRideFareByID(context context.Context, id string) (*RideFareModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*triptype.OsrmApiResponse, error)
	GenerateFare(ctx context.Context, UserID string, Fare []*RideFareModel) ([]*RideFareModel, error)
	EstimateFare(osrm *triptype.OsrmApiResponse) []*RideFareModel
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
}
