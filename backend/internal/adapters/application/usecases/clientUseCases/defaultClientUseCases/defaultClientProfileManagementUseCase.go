package defaultClientUseCases

import (
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultClientProfileManagementUseCase struct {
	clientService ports.ClientService
	clientRepo    ports.ClientRepository
}

func NewClientProfileManagementUseCase(clientService ports.ClientService, clientRepo ports.ClientRepository) ports.ClientProfileManagementUseCase {
	return &DefaultClientProfileManagementUseCase{
		clientService: clientService,
		clientRepo:    clientRepo,
	}
}

func (uc *DefaultClientProfileManagementUseCase) GetProfile(clientID int) (dto.ClientProfileOutput, error) {
	client, err := uc.clientRepo.FindByID(clientID)
	if err != nil {
		return dto.ClientProfileOutput{}, err
	}

	return dto.ClientProfileOutput{
		ClientID:  client.ID,
		SIN:       client.SIN,
		FirstName: client.FirstName,
		LastName:  client.LastName,
		Address:   client.Address,
		Phone:     client.Phone,
		Email:     client.Email,
		JoinDate:  client.JoinDate,
	}, nil
}

func (uc *DefaultClientProfileManagementUseCase) UpdateProfile(input dto.ClientProfileUpdateInput) (dto.ClientProfileOutput, error) {
	client, err := uc.clientService.UpdateClient(
		input.ClientID,
		input.FirstName,
		input.LastName,
		input.Address,
		input.Phone,
		input.Email,
	)
	if err != nil {
		return dto.ClientProfileOutput{}, err
	}

	return dto.ClientProfileOutput{
		ClientID:  client.ID,
		SIN:       client.SIN,
		FirstName: client.FirstName,
		LastName:  client.LastName,
		Address:   client.Address,
		Phone:     client.Phone,
		Email:     client.Email,
		JoinDate:  client.JoinDate,
	}, nil
}
