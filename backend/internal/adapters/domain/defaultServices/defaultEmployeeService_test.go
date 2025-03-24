package defaultServices_test

import (
	"strings"
	"testing"
	"time"

	"github.com/sql-project-backend/internal/adapters/domain/defaultServices"
	"github.com/sql-project-backend/internal/adapters/framework/driven/db/mocks"
)

func TestHireEmployee_Success(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)
	now := time.Now()

	// Hire a new employee with valid input.
	emp, err := service.HireEmployee(0, "123456789", "John", "Doe", "123 Main St", "555-1234", "john@example.com", "Staff", 1, now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if emp.ID == 0 {
		t.Errorf("expected a non-zero employee ID, got %d", emp.ID)
	}
	if emp.Email != "john@example.com" {
		t.Errorf("expected email %q, got %q", "john@example.com", emp.Email)
	}
}

func TestHireEmployee_DuplicateEmail(t *testing.T) {
	// For this test, we assume our mock repository prevents duplicate emails.
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)
	now := time.Now()
	email := "duplicate@example.com"

	// First hire should succeed.
	_, err := service.HireEmployee(0, "123456789", "John", "Doe", "123 Main St", "555-1234", email, "Staff", 1, now)
	if err != nil {
		t.Fatalf("expected first hire to succeed, got: %v", err)
	}

	// Second hire with the same email should fail.
	_, err = service.HireEmployee(0, "987654321", "Jane", "Smith", "456 Elm St", "555-5678", email, "Staff", 2, now)
	if err == nil {
		t.Fatalf("expected error for duplicate email, got nil")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		t.Errorf("expected duplicate email error, got: %v", err)
	}
}

func TestPromoteEmployeeToManager_Success(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)
	now := time.Now()

	// First, hire an employee.
	emp, err := service.HireEmployee(0, "123456789", "John", "Doe", "123 Main St", "555-1234", "john@example.com", "Staff", 1, now)
	if err != nil {
		t.Fatalf("failed to hire employee: %v", err)
	}

	// Promote the employee to manager.
	mgr, err := service.PromoteEmployeeToManager(emp.ID, "Sales", 3)
	if err != nil {
		t.Fatalf("expected promotion to succeed, got error: %v", err)
	}
	if mgr == nil {
		t.Fatalf("expected a valid manager, got nil")
	}
	if mgr.Department != "Sales" {
		t.Errorf("expected department 'Sales', got %s", mgr.Department)
	}
	if mgr.AuthorizationLevel != 3 {
		t.Errorf("expected authorization level 3, got %d", mgr.AuthorizationLevel)
	}
}

func TestPromoteEmployeeToManager_InvalidInput(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)

	// Negative employee ID.
	_, err := service.PromoteEmployeeToManager(-1, "Sales", 3)
	if err == nil {
		t.Fatal("expected error for negative employee ID")
	}

	// Empty department.
	_, err = service.PromoteEmployeeToManager(1, "", 3)
	if err == nil {
		t.Fatal("expected error for empty department")
	}

	// Invalid authorization level (too low).
	_, err = service.PromoteEmployeeToManager(1, "Sales", 0)
	if err == nil {
		t.Fatal("expected error for invalid authorization level (too low)")
	}
	// Invalid authorization level (too high).
	_, err = service.PromoteEmployeeToManager(1, "Sales", 6)
	if err == nil {
		t.Fatal("expected error for invalid authorization level (too high)")
	}
}

func TestPromoteEmployeeToManager_EmployeeNotFound(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)

	// Attempt to promote a non-existent employee.
	_, err := service.PromoteEmployeeToManager(999, "Sales", 3)
	if err == nil {
		t.Fatal("expected error for non-existent employee, got nil")
	}
}

func TestFireEmployee_Success(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)
	now := time.Now()

	emp, err := service.HireEmployee(0, "123456789", "John", "Doe", "123 Main St", "555-1234", "john@example.com", "Staff", 1, now)
	if err != nil {
		t.Fatalf("HireEmployee failed: %v", err)
	}

	firedEmp, err := service.FireEmployee(emp.ID)
	if err != nil {
		t.Fatalf("FireEmployee failed: %v", err)
	}
	if firedEmp.ID != emp.ID {
		t.Errorf("expected fired employee ID %d, got %d", emp.ID, firedEmp.ID)
	}

	// Verify the employee is removed.
	removedEmp, err := mockRepo.FindByID(emp.ID)
	if err == nil && removedEmp != nil {
		t.Errorf("expected employee to be removed, but found one")
	}
}

func TestFireEmployee_InvalidEmployeeID(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)

	_, err := service.FireEmployee(-1)
	if err == nil {
		t.Fatal("expected error for negative employee ID, got nil")
	}
}

func TestFireEmployee_EmployeeNotFound(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)

	_, err := service.FireEmployee(999)
	if err == nil {
		t.Fatal("expected error for non-existent employee, got nil")
	}
}

func TestUpdateEmployee_Success(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)
	now := time.Now()

	emp, err := service.HireEmployee(0, "123456789", "John", "Doe", "123 Main St", "555-1234", "john@example.com", "Staff", 1, now)
	if err != nil {
		t.Fatalf("HireEmployee failed: %v", err)
	}

	updatedEmp, err := service.UpdateEmployee(emp.ID, "Jane", "Smith", "456 New Ave", "555-5678", "jane@example.com", "Manager", 2)
	if err != nil {
		t.Fatalf("UpdateEmployee failed: %v", err)
	}
	if updatedEmp.FirstName != "Jane" || updatedEmp.LastName != "Smith" {
		t.Errorf("UpdateEmployee did not update names correctly: got %s %s", updatedEmp.FirstName, updatedEmp.LastName)
	}
	if updatedEmp.Address != "456 New Ave" || updatedEmp.Phone != "555-5678" || updatedEmp.Email != "jane@example.com" {
		t.Errorf("UpdateEmployee did not update contact info correctly")
	}
	if updatedEmp.Position != "Manager" || updatedEmp.HotelID != 2 {
		t.Errorf("UpdateEmployee did not update position or hotel ID correctly")
	}
}

func TestUpdateEmployee_EmployeeNotFound(t *testing.T) {
	mockRepo := mocks.NewMockEmployeeRepository()
	service := defaultServices.NewEmployeeService(mockRepo)

	_, err := service.UpdateEmployee(999, "Jane", "Smith", "456 New Ave", "555-5678", "jane@example.com", "Manager", 2)
	if err == nil {
		t.Fatal("expected error for non-existent employee, got nil")
	}
}
