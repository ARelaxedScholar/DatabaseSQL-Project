package mocks

import (
	"errors"
	"sync"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type MockStayRepository struct {
	mu    sync.Mutex
	stays map[int]*models.Stay
	nextID int
}

func NewMockStayRepository() ports.StayRepository {
	return &MockStayRepository{
		stays: make(map[int]*models.Stay),
		nextID: 1,
	}
}

func (r *MockStayRepository) Save(stay *models.Stay) (*models.Stay, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stay.ID = r.nextID
	r.nextID++
	r.stays[stay.ID] = stay
	return stay, nil
}

func (r *MockStayRepository) FindByID(id int) (*models.Stay, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stay, exists := r.stays[id]
	if !exists {
		return nil, errors.New("stay not found")
	}
	return stay, nil
}

func (r *MockStayRepository) Update(stay *models.Stay) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.stays[stay.ID]; !exists {
		return errors.New("stay not found")
	}
	r.stays[stay.ID] = stay
	return nil
}

func (r *MockStayRepository) EndStay(stay, employeeID int) error {
	return nil
}

func (r *MockStayRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.stays[id]; !exists {
		return errors.New("stay not found")
	}
	delete(r.stays, id)
	return nil
}
