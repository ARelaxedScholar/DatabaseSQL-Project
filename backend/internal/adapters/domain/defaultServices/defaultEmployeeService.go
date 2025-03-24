package defaultServices

import (
	"errors"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultEmployeeService struct {
	employeeRepo ports.EmployeeRepository
}

func NewEmployeeService(repo ports.EmployeeRepository) ports.EmployeeService {
	return &DefaultEmployeeService{
		employeeRepo: repo,
	}
}

func (s *DefaultEmployeeService) HireEmployee(id int, sin, firstName, lastName, address, phone, email, position string, hotelId int, hireDate time.Time) (*models.Employee, error) {
	emp, err := models.NewEmployee(sin, firstName, lastName, address, phone, email, position, id, hotelId, hireDate)
	if err != nil {
		return nil, err
	}
	dbEmployee, err := s.employeeRepo.Save(emp)
	if err != nil {
		return nil, err
	}
	return dbEmployee, nil
}

func (s *DefaultEmployeeService) PromoteEmployeeToManager(employeeId int, department string, authorizationLevel int) (*models.Manager, error) {
	var err error
	switch {
	case employeeId <= 0:
		err = errors.New("Employee's ID cannot be negative.")
	case department == "":
		err = errors.New("Employee's department cannot be empty.")
	case authorizationLevel < 1 || authorizationLevel > 5:
		err = errors.New("Employee's authorization level must be between 1 and 5.")
	}
	if err != nil {
		return nil, err
	}

	emp, err := s.employeeRepo.FindByID(employeeId)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, errors.New("Employee not found.")
	}

	mgr, err := models.NewManager(emp.SIN, emp.FirstName, emp.LastName, emp.Address, emp.Phone, emp.Email, emp.Position, emp.ID, emp.HotelID, emp.HireDate, department, authorizationLevel)
	if err != nil {
		return nil, err
	}

	if err = s.employeeRepo.UpdateManager(mgr); err != nil {
		return nil, err
	}
	return mgr, nil
}

func (s *DefaultEmployeeService) FireEmployee(employeeId int) (*models.Employee, error) {
	if employeeId <= 0 {
		return nil, errors.New("Employee's ID cannot be negative.")
	}

	emp, err := s.employeeRepo.FindByID(employeeId)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, errors.New("Employee not found.")
	}

	if err = s.employeeRepo.Delete(employeeId); err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *DefaultEmployeeService) UpdateEmployee(employeeId int, firstName, lastName, address, phone, email, position string, hotelId int) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(employeeId)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, errors.New("employee not found")
	}

	employee.FirstName = firstName
	employee.LastName = lastName
	employee.Address = address
	employee.Phone = phone
	employee.Email = email
	employee.Position = position
	employee.HotelID = hotelId

	updatedEmployee, err := s.employeeRepo.UpdateEmployee(employee)
	if err != nil {
		return nil, err
	}

	return updatedEmployee, nil
}
