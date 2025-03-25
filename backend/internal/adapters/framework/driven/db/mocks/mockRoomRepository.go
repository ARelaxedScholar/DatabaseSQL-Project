package mocks

import (
	"errors"
	"sync"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type MockRoomRepository struct {
	mu     sync.Mutex
	rooms  map[int]*models.Room
	nextID int
}

func NewMockRoomRepository() ports.RoomRepository {
	return &MockRoomRepository{
		rooms:  make(map[int]*models.Room),
		nextID: 1,
	}
}

func (r *MockRoomRepository) Save(room *models.Room) (*models.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	room.ID = r.nextID
	r.nextID++
	r.rooms[room.ID] = room
	return room, nil
}

func (r *MockRoomRepository) FindByID(id int) (*models.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	room, exists := r.rooms[id]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (r *MockRoomRepository) Update(room *models.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.rooms[room.ID]; !exists {
		return errors.New("room not found")
	}
	r.rooms[room.ID] = room
	return nil
}

func (r *MockRoomRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.rooms[id]; !exists {
		return errors.New("room not found")
	}
	delete(r.rooms, id)
	return nil
}

func (r *MockRoomRepository) FindAvailableRooms(hotelID int, startDate time.Time, endDate time.Time) ([]*models.Room, error) {
	// In this mock, we simply return all rooms for the hotel.
	r.mu.Lock()
	defer r.mu.Unlock()
	var available []*models.Room
	for _, room := range r.rooms {
		if room.HotelID == hotelID {
			available = append(available, room)
		}
	}
	if len(available) == 0 {
		return nil, errors.New("no available rooms")
	}
	return available, nil
}

func (r *MockRoomRepository) SearchRooms(startDate time.Time, endDate time.Time, capacity int, priceMin, priceMax float64, hotelChainID int, category string) ([]*models.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var results []*models.Room
	for _, room := range r.rooms {
		include := true

		// For mock only check on price and capacity
		// Price filter: Only apply if a non-zero min or max is provided.
		if priceMin > 0 && room.Price < priceMin {
			include = false
		}
		if priceMax > 0 && room.Price > priceMax {
			include = false
		}

		// Capacity filter: Only apply if capacity is provided.
		if capacity > 0 && room.Capacity < capacity {
			include = false
		}

		if include {
			results = append(results, room)
		}
	}

	return results, nil
}
