package mocks

import (
	"errors"
	"sync"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type MockHotelRepository struct {
	mu      sync.Mutex
	hotels  map[int]*models.Hotel
	nextID  int
}

func NewMockHotelRepository() ports.HotelRepository {
	return &MockHotelRepository{
		hotels: make(map[int]*models.Hotel),
		nextID: 1,
	}
}

func (r *MockHotelRepository) Save(hotel *models.Hotel) (*models.Hotel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	hotel.ID = r.nextID
	r.nextID++
	r.hotels[hotel.ID] = hotel
	return hotel, nil
}

func (r *MockHotelRepository) FindByID(id int) (*models.Hotel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	hotel, exists := r.hotels[id]
	if !exists {
		return nil, errors.New("hotel not found")
	}
	return hotel, nil
}

func (r *MockHotelRepository) Update(hotel *models.Hotel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.hotels[hotel.ID]; !exists {
		return errors.New("hotel not found")
	}
	r.hotels[hotel.ID] = hotel
	return nil
}

func (r *MockHotelRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.hotels[id]; !exists {
		return errors.New("hotel not found")
	}
	delete(r.hotels, id)
	return nil
}
