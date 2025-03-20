package adminUseCases

import (
	"errors"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultAdminAccountManagementUseCase struct {
	clientRepo      ports.ClientRepository
	employeeRepo    ports.EmployeeRepository
	clientService   ports.ClientService
	employeeService ports.EmployeeService
}

func NewAdminAccountManagementUseCase(
	clientRepo ports.ClientRepository,
	employeeRepo ports.EmployeeRepository,
	clientService ports.ClientService,
	employeeService ports.EmployeeService,
) ports.AdminAccountManagementUseCase {
	return &DefaultAdminAccountManagementUseCase{
		clientRepo:      clientRepo,
		employeeRepo:    employeeRepo,
		clientService:   clientService,
		employeeService: employeeService,
	}
}

func (uc *DefaultAdminAccountManagementUseCase) GetAccount(accountID int) (dto.AccountOutput, error) {
	client, err := uc.clientRepo.FindByID(accountID)
	if err == nil && client != nil {
		return mapClientToAccountOutput(client), nil
	}
	employee, err := uc.employeeRepo.FindByID(accountID)
	if err == nil && employee != nil {
		return mapEmployeeToAccountOutput(employee), nil
	}
	return dto.AccountOutput{}, errors.New("account not found")
}

func (uc *DefaultAdminAccountManagementUseCase) ListClientAccounts() ([]dto.AccountOutput, error) {
	clients, err := uc.clientRepo.ListAllClients()
	if err != nil {
		return nil, err
	}
	var outputs []dto.AccountOutput
	for _, client := range clients {
		outputs = append(outputs, mapClientToAccountOutput(client))
	}
	return outputs, nil
}

func (uc *DefaultAdminAccountManagementUseCase) CreateClientAccount(input dto.ClientAccountInput) (dto.AccountOutput, error) {
	client, err := uc.clientService.RegisterClient(
		0,
		input.SIN,
		input.FirstName,
		input.LastName,
		input.Address,
		input.Phone,
		input.Email,
		time.Now(),
	)
	if err != nil {
		return dto.AccountOutput{}, err
	}
	return mapClientToAccountOutput(client), nil
}

func (uc *DefaultAdminAccountManagementUseCase) UpdateClientAccount(accountID int, input dto.ClientAccountUpdateInput) (dto.AccountOutput, error) {
	client, err := uc.clientService.UpdateClient(
		accountID,
		input.FirstName,
		input.LastName,
		input.Address,
		input.Phone,
		input.Email,
	)
	if err != nil {
		return dto.AccountOutput{}, err
	}
	return mapClientToAccountOutput(client), nil
}

func (uc *DefaultAdminAccountManagementUseCase) DeleteClientAccount(accountID int) error {
	return uc.clientRepo.Delete(accountID)
}

func (uc *DefaultAdminAccountManagementUseCase) ListEmployeeAccounts() ([]dto.AccountOutput, error) {
	employees, err := uc.employeeRepo.ListAllEmployees()
	if err != nil {
		return nil, err
	}
	var outputs []dto.AccountOutput
	for _, employee := range employees {
		outputs = append(outputs, mapEmployeeToAccountOutput(employee))
	}
	return outputs, nil
}

func (uc *DefaultAdminAccountManagementUseCase) CreateEmployeeAccount(input dto.EmployeeAccountInput) (dto.AccountOutput, error) {
	employee, err := uc.employeeService.HireEmployee(
		0,
		input.SIN,
		input.FirstName,
		input.LastName,
		input.Address,
		input.Phone,
		input.Email,
		input.Position,
		input.HotelID,
		time.Now(),
	)
	if err != nil {
		return dto.AccountOutput{}, err
	}
	return mapEmployeeToAccountOutput(employee), nil
}

func (uc *DefaultAdminAccountManagementUseCase) UpdateEmployeeAccount(accountID int, input dto.EmployeeAccountUpdateInput) (dto.AccountOutput, error) {
	employee, err := uc.employeeService.UpdateEmployee(
		accountID,
		input.FirstName,
		input.LastName,
		input.Address,
		input.Phone,
		input.Email,
		input.Position,
		input.HotelID,
	)
	if err != nil {
		return dto.AccountOutput{}, err
	}
	return mapEmployeeToAccountOutput(employee), nil
}

func (uc *DefaultAdminAccountManagementUseCase) DeleteEmployeeAccount(accountID int) error {
	return uc.employeeRepo.Delete(accountID)
}

func mapClientToAccountOutput(client *models.Client) dto.AccountOutput {
	return dto.AccountOutput{
		AccountID: client.ID,
		FirstName: client.FirstName,
		LastName:  client.LastName,
		Email:     client.Email,
		Role:      "client",
		CreatedAt: client.JoinDate,
		UpdatedAt: client.JoinDate,
	}
}

func mapEmployeeToAccountOutput(employee *models.Employee) dto.AccountOutput {
	return dto.AccountOutput{
		AccountID: employee.ID,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Email:     employee.Email,
		Role:      "employee",
		CreatedAt: employee.HireDate,
		UpdatedAt: employee.HireDate,
	}
}
