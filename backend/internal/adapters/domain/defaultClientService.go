package domain

import (
	"errors"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultClientService struct {
	clientRepo ports.ClientRepository
}

func NewClientService(repo ports.ClientRepository) ports.ClientService {
	return &DefaultClientService{
		clientRepo: repo,
	}
}

func (s *DefaultClientService) RegisterClient(id int, sin, firstName, lastName, address, phone, email string, joinDate time.Time) (*models.Client, error) {
	client, err := models.NewClient(id, sin, firstName, lastName, address, phone, email, joinDate)
	if err != nil {
		return nil, err
	}
	dbClient, err := s.clientRepo.Save(client)
	if err != nil {
		return nil, err
	}
	return dbClient, nil
}

func (s *DefaultClientService) UpdateClient(id int, firstName, lastName, address, phone, email string) (*models.Client, error) {
	client, err := s.clientRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("Client not found.")
	}

	client.FirstName = firstName
	client.LastName = lastName
	client.Address = address
	client.Phone = phone
	client.Email = email

	if client, err = s.clientRepo.Update(client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *DefaultClientService) RemoveClient(id int) error {
	client, err := s.clientRepo.FindByID(id)
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("Client not found.")
	}
	return s.clientRepo.Delete(id)
}
