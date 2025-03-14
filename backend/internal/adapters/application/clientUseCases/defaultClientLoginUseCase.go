package clientUseCases

import (
	"errors"

	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultClientLoginUseCase struct {
	clientRepo   ports.ClientRepository
	tokenService ports.TokenService
}

func NewClientLoginUseCase(clientRepo ports.ClientRepository, tokenService ports.TokenService) ports.ClientLoginUseCase {
	return &DefaultClientLoginUseCase{
		clientRepo:   clientRepo,
		tokenService: tokenService,
	}
}

func (uc *DefaultClientLoginUseCase) Login(input dto.ClientLoginInput) (dto.ClientLoginOutput, error) {
	client, err := uc.clientRepo.FindByEmail(input.Email)
	if err != nil || client == nil {
		return dto.ClientLoginOutput{}, errors.New("Client not found.")
	}

	token, err := uc.tokenService.GenerateToken(client.ID)
	if err != nil {
		return dto.ClientLoginOutput{}, err
	}

	return dto.ClientLoginOutput{
		ClientID: client.ID,
		Token:    token,
	}, nil
}
