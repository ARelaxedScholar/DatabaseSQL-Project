package mocks

import (
	"errors"

	"github.com/sql-project-backend/internal/ports"
)

type MockQueryRepository struct {
	// For simplicity, use in-memory fixed data.
	availableRoomsByZone map[string]int
	hotelRoomCapacities  map[int]int
}

func NewMockQueryRepository() ports.QueryRepository {
	return &MockQueryRepository{
		availableRoomsByZone: map[string]int{
			"North": 10,
			"South": 8,
			"East":  12,
			"West":  5,
		},
		hotelRoomCapacities: map[int]int{
			1: 100,
			2: 80,
			3: 120,
		},
	}
}

func (r *MockQueryRepository) GetAvailableRoomsByZone() (map[string]int, error) {
	if r.availableRoomsByZone == nil {
		return nil, errors.New("no data available")
	}
	return r.availableRoomsByZone, nil
}

func (r *MockQueryRepository) GetHotelRoomCapacity(hotelId int) (int, error) {
	if capacity, exists := r.hotelRoomCapacities[hotelId]; exists {
		return capacity, nil
	}
	return 0, errors.New("hotel not found")
}
