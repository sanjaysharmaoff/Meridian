package domain

import (
	"meridian/services/trip-service/pkg/types"
	pb "meridian/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string
	TotalPriceInCents float64
	Route             *types.OsrmApiResponse
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(r []*RideFareModel) []*pb.RideFare {
	arr := make([]*pb.RideFare, 0, len(r))

	for _, f := range r {
		arr = append(arr, f.ToProto())
	}
	return arr
}
