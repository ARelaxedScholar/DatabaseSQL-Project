package mocks

import (
	"errors"
	"sync"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type MockClientRepository struct {
	mu      sync.Mutex
	clients map[int]*models.Client
	nextID  int
}

func NewMockClientRepository() ports.ClientRepository {
	return &MockClientRepository{
		clients: make(map[int]*models.Client),
		nextID:  1,
	}
}

func (r *MockClientRepository) Save(client *models.Client) (*models.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	client.ID = r.nextID
	r.nextID++
	r.clients[client.ID] = client
	return client, nil
}

func (r *MockClientRepository) FindByID(id int) (*models.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	client, exists := r.clients[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (r *MockClientRepository) FindByEmail(email string) (*models.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, client := range r.clients {
		if client.Email == email {
			return client, nil
		}
	}
	return nil, errors.New("client not found")
}

func (r *MockClientRepository) ListAllClients() ([]*models.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.Client
	for _, client := range r.clients {
		list = append(list, client)
	}
	return list, nil
}

func (r *MockClientRepository) Update(client *models.Client) (*models.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.clients[client.ID]; !exists {
		return nil, errors.New("client not found")
	}
	r.clients[client.ID] = client
	return client, nil
}

func (r *MockClientRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.clients[id]; !exists {
		return errors.New("client not found")
	}
	delete(r.clients, id)
	return nil
}
