package defaultClientUseCases

import (
	"errors"
	"fmt"
	"time"

	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultClientLoginUseCase struct {
	clientRepo   ports.ClientRepository
	tokenService ports.TokenService
	emailService ports.EmailService
	appLink      string
}

func NewClientLoginUseCase(clientRepo ports.ClientRepository, tokenService ports.TokenService, emailService ports.EmailService, appLink string) ports.ClientLoginUseCase {
	return &DefaultClientLoginUseCase{
		clientRepo:   clientRepo,
		tokenService: tokenService,
		emailService: emailService,
		appLink:      appLink,
	}
}

func (uc *DefaultClientLoginUseCase) Login(input dto.ClientLoginInput) (dto.ClientLoginOutput, error) {
	client, err := uc.clientRepo.FindByEmail(input.Email)
	if err != nil || client == nil {
		return dto.ClientLoginOutput{}, errors.New("Client not found.")
	}

	token, err := uc.tokenService.GenerateTokenWithDuration(client.ID, "client", 10*time.Minute) // sends a link with a token that is valid for 10 minutes
	if err != nil {
		return dto.ClientLoginOutput{}, err
	}

	// Create a login link with the token
	loginLink := fmt.Sprintf("https://%s/magic?token=%s", uc.appLink, token)

	if err := uc.emailService.SendLoginLink(client.Email, loginLink); err != nil {
		return dto.ClientLoginOutput{}, errors.New("failed to send login email")
	}

	return dto.ClientLoginOutput{
		Message: "A login link has been sent to your email address. Please check your email to proceed.",
	}, nil
}

func (uc *DefaultClientLoginUseCase) MagicLogin(tokenString string) (dto.MagicLoginOutput, error) {
	// Validate the short-lived temporary token
	clientID, role, err := uc.tokenService.ValidateToken(tokenString)
	if err != nil {
		return dto.MagicLoginOutput{}, errors.New("invalid or expired token")
	}

	// Generate a new, longer-lived session token
	sessionToken, err := uc.tokenService.GenerateTokenWithDuration(clientID, role, 24*time.Hour) // generates a new longer lived session token
	if err != nil {
		return dto.MagicLoginOutput{}, err
	}

	return dto.MagicLoginOutput{
		Message:      "You have been successfully logged in.",
		SessionToken: sessionToken,
	}, nil
}
