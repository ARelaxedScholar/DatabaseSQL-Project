package defaultEmployeeUseCases

import (
	"errors"
	"fmt"
	"time"

	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultEmployeeLoginUseCase struct {
	employeeRepo ports.EmployeeRepository
	tokenService ports.TokenService
	emailService ports.EmailService
	appLink      string
}

func NewEmployeeLoginUseCase(employeeRepo ports.EmployeeRepository, tokenService ports.TokenService, emailService ports.EmailService, appLink string) ports.EmployeeLoginUseCase {
	return &DefaultEmployeeLoginUseCase{
		employeeRepo: employeeRepo,
		tokenService: tokenService,
		emailService: emailService,
		appLink:      appLink,
	}
}

func (uc *DefaultEmployeeLoginUseCase) Login(input dto.EmployeeLoginInput) (dto.EmployeeLoginOutput, error) {
	employee, err := uc.employeeRepo.FindByEmail(input.Email)
	if err != nil || employee == nil {
		return dto.EmployeeLoginOutput{}, errors.New("employee not found")
	}

	// Generate a temporary token valid for 10 minutes.
	token, err := uc.tokenService.GenerateTokenWithDuration(employee.ID, "employee", 10*time.Minute)
	if err != nil {
		return dto.EmployeeLoginOutput{}, err
	}

	// Create a login link with the token.
	loginLink := fmt.Sprintf("%s?token=%s&role=employe", uc.appLink, token)

	// Send the login link via the email service.
	if err := uc.emailService.SendLoginLink(employee.Email, loginLink); err != nil {
		return dto.EmployeeLoginOutput{}, errors.New("failed to send login email")
	}

	return dto.EmployeeLoginOutput{
		Message: "A login link has been sent to your email address. Please check your email to proceed.",
	}, nil
}

func (uc *DefaultEmployeeLoginUseCase) MagicLogin(tokenString string) (dto.MagicLoginOutput, error) {
	// Validate the short-lived temporary token
	employeeID, role, err := uc.tokenService.ValidateToken(tokenString)
	if err != nil {
		return dto.MagicLoginOutput{}, errors.New("invalid or expired token")
	}

	// Generate a new, longer-lived session token
	sessionToken, err := uc.tokenService.GenerateTokenWithDuration(employeeID, role, 24*time.Hour) // generates a new longer lived session token
	if err != nil {
		return dto.MagicLoginOutput{}, err
	}

	return dto.MagicLoginOutput{
		Message:      "You have been successfully logged in.",
		SessionToken: sessionToken,
	}, nil
}
