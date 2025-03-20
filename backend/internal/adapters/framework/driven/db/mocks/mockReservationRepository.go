package mocks

import (
	"errors"
	"sync"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type MockReservationRepository struct {
	mu           sync.Mutex
	reservations map[int]*models.Reservation
	nextID       int
}

func NewMockReservationRepository() ports.ReservationRepository {
	return &MockReservationRepository{
		reservations: make(map[int]*models.Reservation),
		nextID:       1,
	}
}

func (r *MockReservationRepository) Save(reservation *models.Reservation) (*models.Reservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	reservation.ID = r.nextID
	r.nextID++
	r.reservations[reservation.ID] = reservation
	return reservation, nil
}

func (r *MockReservationRepository) FindByID(id int) (*models.Reservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	reservation, exists := r.reservations[id]
	if !exists {
		return nil, errors.New("reservation not found")
	}
	return reservation, nil
}

func (r *MockReservationRepository) GetByClient(clientID int) ([]*models.Reservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.Reservation
	for _, reservation := range r.reservations {
		if reservation.ClientID == clientID {
			list = append(list, reservation)
		}
	}
	return list, nil
}

func (r *MockReservationRepository) Update(reservation *models.Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.reservations[reservation.ID]; !exists {
		return errors.New("reservation not found")
	}
	r.reservations[reservation.ID] = reservation
	return nil
}

func (r *MockReservationRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.reservations[id]; !exists {
		return errors.New("reservation not found")
	}
	delete(r.reservations, id)
	return nil
}
