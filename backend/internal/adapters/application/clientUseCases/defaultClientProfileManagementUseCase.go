package clientUseCases

import (
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultClientProfileManagementUseCase struct {
	clientService ports.ClientService
}

func NewClientProfileManagementUseCase(clientService ports.ClientService) ports.ClientProfileManagementUseCase {
	return &DefaultClientProfileManagementUseCase{
		clientService: clientService,
	}
}

func (uc *DefaultClientProfileManagementUseCase) GetProfile(clientID int) (dto.ClientProfileOutput, error) {
	client, err := uc.clientService.GetClientByID(clientID)
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
