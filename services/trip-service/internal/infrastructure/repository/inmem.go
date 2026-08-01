package repository

import (
	"context"
	"fmt"
	"meridian/services/trip-service/internal/domain"
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
