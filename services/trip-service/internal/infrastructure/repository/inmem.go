package repository

import (
	"context"
	"fmt"
	"meridian/services/trip-service/internal/domain"
	pbd "meridian/shared/proto/driver"
	pb "meridian/shared/proto/trip"
)

type inmemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func (r *inmemRepository) CreateTrip(context context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inmemRepository) SaveRideFare(context context.Context, d *domain.RideFareModel) error {
	r.rideFares[d.ID.Hex()] = d
	return nil
}

func (r *inmemRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	fare, found := r.rideFares[id]
	if found == false {
		return nil, fmt.Errorf("fare not found for id: %v", id)
	}
	return fare, nil
}

func (r *inmemRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	trip, ok := r.trips[id]
	if !ok {
		return nil, nil
	}
	return trip, nil
}

func (r *inmemRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error {
	trip, ok := r.trips[tripID]
	if !ok {
		return fmt.Errorf("trip not found with ID %s", tripID)
	}
	trip.Status = status
	if driver != nil {
		trip.Driver = &pb.TripDriver{
			Id:             driver.Id,
			Name:           driver.Name,
			CarPlate:       driver.CarPlate,
			ProfilePicture: driver.ProfilePicture,
		}
	}
	return nil
}
