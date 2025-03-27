package mocks

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports" // Import ports package
)

type MockRoomRepository struct {
	mu               sync.Mutex
	rooms            map[int]*models.Room
	nextID           int
	findByIDError    error
	deleteError      error
	findAvailError   error
	searchRoomsError error
	updateError      error
	saveError        error
}

func NewMockRoomRepository() *MockRoomRepository {
	return &MockRoomRepository{
		rooms:  make(map[int]*models.Room),
		nextID: 1,
	}
}

func (r *MockRoomRepository) SetFindByIDError(err error) {
	r.mu.Lock()
	r.findByIDError = err
	r.mu.Unlock()
}
func (r *MockRoomRepository) SetDeleteError(err error) {
	r.mu.Lock()
	r.deleteError = err
	r.mu.Unlock()
}
func (r *MockRoomRepository) SetUpdateError(err error) {
	r.mu.Lock()
	r.updateError = err
	r.mu.Unlock()
}
func (r *MockRoomRepository) SetSaveError(err error) { r.mu.Lock(); r.saveError = err; r.mu.Unlock() }

func (r *MockRoomRepository) Save(room *models.Room) (*models.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.saveError != nil {
		return nil, r.saveError
	}
	if room == nil {
		return nil, errors.New("Cannot save nil room.")
	}
	for _, existing := range r.rooms {
		if existing.HotelID == room.HotelID && existing.Number == room.Number && room.ID == 0 {
			// Simulate unique constraint
			return nil, fmt.Errorf("Mock Error: Room number %s already exists in hotel %d.", room.Number, room.HotelID)
		}
	}
	if room.ID == 0 {
		room.ID = r.nextID
		r.nextID++
	}
	savedRoom := *room
	r.rooms[savedRoom.ID] = &savedRoom
	return &savedRoom, nil
}

func (r *MockRoomRepository) FindByID(id int) (*models.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.findByIDError != nil {
		return nil, r.findByIDError
	}
	room, exists := r.rooms[id]
	if !exists {
		return nil, models.ErrNotFound
	}
	foundRoom := *room
	return &foundRoom, nil
}

func (r *MockRoomRepository) Update(room *models.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.updateError != nil {
		return r.updateError
	}
	if room == nil {
		return errors.New("Cannot update with nil room.")
	}
	_, exists := r.rooms[room.ID]
	if !exists {
		return models.ErrNotFound
	}
	for _, rCheck := range r.rooms {
		if rCheck.ID != room.ID && rCheck.HotelID == room.HotelID && rCheck.Number == room.Number {
			// Simulate unique constraint
			return fmt.Errorf("Mock Error: Cannot update, room number %s already exists in hotel %d.", room.Number, room.HotelID)
		}
	}
	updatedRoom := *room
	r.rooms[room.ID] = &updatedRoom
	return nil
}

func (r *MockRoomRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.deleteError != nil {
		return r.deleteError
	}
	if _, exists := r.rooms[id]; !exists {
		return models.ErrNotFound // Use ports.ErrNotFound
	}
	delete(r.rooms, id)
	return nil
}

func (r *MockRoomRepository) FindAvailableRooms(hotelID int, startDate time.Time, endDate time.Time) ([]*models.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.findAvailError != nil {
		return nil, r.findAvailError
	}
	var available []*models.Room
	for _, room := range r.rooms {
		if room.HotelID == hotelID {
			roomCopy := *room
			available = append(available, &roomCopy)
		}
	}
	return available, nil
}

func (r *MockRoomRepository) SearchRooms(startDate time.Time, endDate time.Time, capacity int, priceMin, priceMax float64, hotelChainID int, roomType models.RoomType) ([]*models.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.searchRoomsError != nil {
		return nil, r.searchRoomsError
	}
	var results []*models.Room
	for _, room := range r.rooms {
		include := true
		if roomType != 0 && room.RoomType != roomType {
			include = false
		}
		if priceMin > 0 && room.Price < priceMin {
			include = false
		}
		if priceMax > 0 && room.Price > priceMax {
			include = false
		}
		if capacity > 0 && room.Capacity < capacity {
			include = false
		}
		// TODO: Add HotelChainID filtering
		// TODO: Add Availability check
		if include {
			roomCopy := *room
			results = append(results, &roomCopy)
		}
	}
	return results, nil
}

func (r *MockRoomRepository) CountRoomsForHotel(hotelID int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for _, room := range r.rooms {
		if room.HotelID == hotelID {
			count++
		}
	}
	return count
}

var _ ports.RoomRepository = (*MockRoomRepository)(nil)
