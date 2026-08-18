package main

import (
	"math/rand"
	pb "meridian/shared/proto/driver"
	"meridian/shared/util"
	"sync"

	"github.com/mmcloughlin/geohash"
)

type driverInMap struct {
	Driver            *pb.Driver
	CurrentRoute      [][]float64
	CurrentRouteIndex int
}

type Service struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

func NewService() *Service {
	return &Service{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *Service) RegisterDriver(driverId string, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	randomIdx := rand.Intn(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIdx]
	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])
	latitude := randomRoute[0][0]
	longitude := randomRoute[0][1]
	randomPlate := GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIdx)
	driver := &pb.Driver{
		Geohash:     geohash,
		Id:          driverId,
		Name:        "Sanjay Verstappen",
		PackageSlug: packageSlug,
		Location: &pb.Location{
			Latitude:  latitude,
			Longitude: longitude,
		},
		CarPlate:       randomPlate,
		ProfilePicture: randomAvatar,
	}
	driverInMap := &driverInMap{
		Driver:            driver,
		CurrentRoute:      randomRoute,
		CurrentRouteIndex: 0,
	}

	s.drivers = append(s.drivers, driverInMap)
	return driver, nil
}

func (s *Service) UnregisterDriver(driverId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.drivers); i++ {
		if s.drivers[i].Driver.Id == driverId {
			s.drivers = append(s.drivers[:i], s.drivers[i+1:]...)
			return
		}
	}
}

func (s *Service) FindAvailableDrivers(packagetype string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var matchingDrivers []string
	for _, j := range s.drivers {
		if j.Driver.PackageSlug == packagetype {
			matchingDrivers = append(matchingDrivers, j.Driver.Id)
		}
	}

	if len(matchingDrivers) == 0 {
		return []string{}
	}
	return matchingDrivers
}
