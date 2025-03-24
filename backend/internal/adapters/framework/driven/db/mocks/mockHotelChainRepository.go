package mocks

import (
	"errors"
	"sync"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type MockHotelChainRepository struct {
	mu          sync.Mutex
	hotelChains map[int]*models.HotelChain
	nextID      int
}

func NewMockHotelChainRepository() ports.HotelChainRepository {
	return &MockHotelChainRepository{
		hotelChains: make(map[int]*models.HotelChain),
		nextID:      1,
	}
}

func (r *MockHotelChainRepository) Save(chain *models.HotelChain) (*models.HotelChain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Check for duplicate by chain name
	for _, existingChain := range r.hotelChains {
		if existingChain.Name == chain.Name {
			return nil, errors.New("duplicate hotel chain: a chain with this name already exists")
		}
	}
	chain.ID = r.nextID
	r.nextID++
	r.hotelChains[chain.ID] = chain
	return chain, nil
}

func (r *MockHotelChainRepository) FindByID(id int) (*models.HotelChain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	chain, exists := r.hotelChains[id]
	if !exists {
		return nil, errors.New("hotel chain not found")
	}
	return chain, nil
}

func (r *MockHotelChainRepository) Update(chain *models.HotelChain) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.hotelChains[chain.ID]; !exists {
		return errors.New("hotel chain not found")
	}
	r.hotelChains[chain.ID] = chain
	return nil
}

func (r *MockHotelChainRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.hotelChains[id]; !exists {
		return errors.New("hotel chain not found")
	}
	delete(r.hotelChains, id)
	return nil
}
