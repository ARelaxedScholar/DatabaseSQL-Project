package defaultEmployeeUseCases

import (
	"errors"

	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultEmployeeLoginUseCase struct {
	employeeRepo ports.EmployeeRepository
	tokenService ports.TokenService
}

func NewEmployeeLoginUseCase(employeeRepo ports.EmployeeRepository, tokenService ports.TokenService) ports.EmployeeLoginUseCase {
	return &DefaultEmployeeLoginUseCase{
		employeeRepo: employeeRepo,
		tokenService: tokenService,
	}
}

func (uc *DefaultEmployeeLoginUseCase) Login(input dto.EmployeeLoginInput) (dto.EmployeeLoginOutput, error) {
	employee, err := uc.employeeRepo.FindByEmail(input.Email)
	if err != nil || employee == nil {
		return dto.EmployeeLoginOutput{}, errors.New("employee not found")
	}

	token, err := uc.tokenService.GenerateToken(employee.ID)
	if err != nil {
		return dto.EmployeeLoginOutput{}, err
	}

	return dto.EmployeeLoginOutput{
		EmployeeID: employee.ID,
		Token:      token,
	}, nil
}
