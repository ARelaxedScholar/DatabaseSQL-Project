package mocks

import (
	"errors"
	"sync"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type MockEmployeeRepository struct {
	mu        sync.Mutex
	employees map[int]*models.Employee
	nextID    int
}

func NewMockEmployeeRepository() ports.EmployeeRepository {
	return &MockEmployeeRepository{
		employees: make(map[int]*models.Employee),
		nextID:    1,
	}
}

func (r *MockEmployeeRepository) Save(emp *models.Employee) (*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate email
	for _, existingEmp := range r.employees {
		if existingEmp.Email == emp.Email {
			return nil, errors.New("duplicate email")
		}
	}

	emp.ID = r.nextID
	r.nextID++
	r.employees[emp.ID] = emp
	return emp, nil
}

func (r *MockEmployeeRepository) FindByID(id int) (*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	emp, exists := r.employees[id]
	if !exists {
		return nil, errors.New("employee not found")
	}
	return emp, nil
}

func (r *MockEmployeeRepository) FindByEmail(email string) (*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, emp := range r.employees {
		if emp.Email == email {
			return emp, nil
		}
	}
	return nil, errors.New("employee not found")
}

func (r *MockEmployeeRepository) ListAllEmployees() ([]*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.Employee
	for _, emp := range r.employees {
		list = append(list, emp)
	}
	return list, nil
}

func (r *MockEmployeeRepository) UpdateEmployee(emp *models.Employee) (*models.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.employees[emp.ID]; !exists {
		return nil, errors.New("employee not found")
	}
	r.employees[emp.ID] = emp
	return emp, nil
}

func (r *MockEmployeeRepository) UpdateManager(mgr *models.Manager) error {
	// For simplicity, we assume UpdateManager is similar to updating an employee.
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.employees[mgr.ID]; !exists {
		return errors.New("manager not found")
	}
	// In a real implementation, update manager-specific fields.
	r.employees[mgr.ID] = &models.Employee{
		ID:        mgr.ID,
		FirstName: mgr.FirstName,
		LastName:  mgr.LastName,
		Email:     mgr.Email,
		// … other fields as needed
	}
	return nil
}

func (r *MockEmployeeRepository) Delete(employeeID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.employees[employeeID]; !exists {
		return errors.New("employee not found")
	}
	delete(r.employees, employeeID)
	return nil
}
