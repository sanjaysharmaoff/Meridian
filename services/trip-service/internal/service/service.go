package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meridian/services/trip-service/internal/domain"
	triptype "meridian/services/trip-service/pkg/types"
	pb "meridian/shared/proto/trip"
	"meridian/shared/types"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	trip := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
		Driver:   &pb.TripDriver{},
	}

	return s.repo.CreateTrip(ctx, trip)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*triptype.OsrmApiResponse, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}
	defer resp.Body.Close()
	var osrm triptype.OsrmApiResponse
	err = json.NewDecoder(resp.Body).Decode(&osrm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode  from OSRM API: %v", err)
	}
	return &osrm, nil
}

func (s *service) EstimateFare(osrm *triptype.OsrmApiResponse) []*domain.RideFareModel {
	basefare := getBaseFares()
	for i, it := range basefare {
		basefare[i] = estimateFareRoute(osrm, it)
	}
	return basefare
}

func estimateFareRoute(osrm *triptype.OsrmApiResponse, f *domain.RideFareModel) *domain.RideFareModel {
	basePackagePrice := f.TotalPriceInCents
	priceConfig := triptype.DefaultPriceConfig()
	distancePrice := osrm.Routes[0].Distance * priceConfig.PricePerUnitOfDistance
	durationMinutes := osrm.Routes[0].Duration / 60
	timePrice := durationMinutes * priceConfig.PricingPerMinute
	totalPrice := basePackagePrice + distancePrice + timePrice
	return &domain.RideFareModel{
		PackageSlug:       f.PackageSlug,
		TotalPriceInCents: totalPrice,
	}
}

func getBaseFares() []*domain.RideFareModel {
	var basefare []*domain.RideFareModel
	basefare = append(basefare, &domain.RideFareModel{
		PackageSlug:       "suv",
		TotalPriceInCents: 200,
	})
	basefare = append(basefare, &domain.RideFareModel{
		PackageSlug:       "sedan",
		TotalPriceInCents: 350,
	})
	basefare = append(basefare, &domain.RideFareModel{
		PackageSlug:       "van",
		TotalPriceInCents: 400,
	})
	basefare = append(basefare, &domain.RideFareModel{
		PackageSlug:       "luxury",
		TotalPriceInCents: 1000,
	})
	return basefare
}

func (h *service) GenerateFare(ctx context.Context, UserID string, Fare []*domain.RideFareModel) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(Fare))

	for i := range fares {
		id := primitive.NewObjectID()
		fare := &domain.RideFareModel{
			ID:                id,
			PackageSlug:       Fare[i].PackageSlug,
			UserID:            UserID,
			TotalPriceInCents: Fare[i].TotalPriceInCents,
		}
		if err := h.repo.SaveRideFare(ctx, fare); err != nil {
			log.Print("error saving the fare :", err)
			return nil, fmt.Errorf("failed to save trip fare: %w", err)
		}
		fares[i] = fare
	}
	return fares, nil
}

func (s *service) GetAndValidateFare(ctx context.Context, fareID, userID string) (*domain.RideFareModel, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fare: %w", err)
	}
	if fare == nil {
		return nil, fmt.Errorf("fare does not exist")
	}

	if fare.UserID != userID {
		return nil, fmt.Errorf("userID does not match the fair")
	}
	return fare, nil
}
