package clientUseCases

import (
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultClientRegistrationUseCase struct {
	clientService ports.ClientService
}

func NewClientRegistrationUseCase(clientService ports.ClientService) ports.ClientRegistrationUseCase {
	return &DefaultClientRegistrationUseCase{
		clientService: clientService,
	}
}

func (uc *DefaultClientRegistrationUseCase) RegisterClient(input dto.ClientRegistrationInput) (dto.ClientRegistrationOutput, error) {
	client, err := uc.clientService.RegisterClient(
		0, // Pass a default value, we leave it to the DB to actually initialize this
		input.SIN,
		input.FirstName,
		input.LastName,
		input.Address,
		input.Phone,
		input.Email,
		input.JoinDate,
	)
	if err != nil {
		return dto.ClientRegistrationOutput{}, err
	}

	return dto.ClientRegistrationOutput{
		ClientID: client.ID,
	}, nil
}
