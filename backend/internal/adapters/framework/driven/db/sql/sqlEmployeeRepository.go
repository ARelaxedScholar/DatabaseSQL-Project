package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type PostgresEmployeeRepository struct {
	db *sql.DB
}

func NewPostgresEmployeeRepository(db *sql.DB) (ports.EmployeeRepository, error) {
	if db == nil {
		return nil, errors.New("Db connection pool cannot be nil.")
	}
	return &PostgresEmployeeRepository{db: db}, nil
}

var _ ports.EmployeeRepository = (*PostgresEmployeeRepository)(nil)

// Helper to scan employee data
func scanEmployee(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Employee, error) {
	emp := &models.Employee{}
	var hireDate time.Time

	err := scanner.Scan(
		&emp.ID,
		&emp.SIN,
		&emp.FirstName,
		&emp.LastName,
		&emp.Address,
		&emp.Phone,
		&emp.Email,
		&emp.HotelID,
		&emp.Position,
		&hireDate,
	)
	if err != nil {
		return nil, err // Let caller handle specific errors like ErrNotFound
	}
	emp.HireDate = hireDate
	return emp, nil
}

func (r *PostgresEmployeeRepository) Save(emp *models.Employee) (*models.Employee, error) {
	if emp == nil {
		return nil, errors.New("Cannot save a nil employee.")
	}
	if emp.SIN == "" || len(emp.SIN) != 9 || emp.FirstName == "" || emp.LastName == "" || emp.Email == "" || emp.HotelID <= 0 || emp.Position == "" || emp.HireDate.IsZero() {
		return nil, errors.New("Invalid employee data provided for save.")
	}

	query := `
		INSERT INTO employee (sin, first_name, last_name, address, phone, email, hotel_id, position, hire_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := r.db.QueryRow(query,
		emp.SIN,
		emp.FirstName,
		emp.LastName,
		emp.Address,
		emp.Phone,
		emp.Email,
		emp.HotelID,
		emp.Position,
		emp.HireDate,
	).Scan(&emp.ID)

	if err != nil {
		// Checks unique sin/email, FK violation for hotel_id
		return nil, handlePqError(err)
	}

	return emp, nil
}

func (r *PostgresEmployeeRepository) FindByID(id int) (*models.Employee, error) {
	if id <= 0 {
		return nil, errors.New("Invalid employee ID provided.")
	}

	query := `
		SELECT id, sin, first_name, last_name, address, phone, email, hotel_id, position, hire_date
		FROM employee
		WHERE id = $1`

	row := r.db.QueryRow(query, id)
	e, err := scanEmployee(row)

	if err != nil {
		return nil, handlePqError(err) // Handles ErrNotFound
	}
	return e, nil
}

func (r *PostgresEmployeeRepository) FindByEmail(email string) (*models.Employee, error) {
	if email == "" {
		return nil, errors.New("Email cannot be empty for lookup.")
	}

	query := `
		SELECT id, sin, first_name, last_name, address, phone, email, hotel_id, position, hire_date
		FROM employee
		WHERE email = $1`

	row := r.db.QueryRow(query, email)
	e, err := scanEmployee(row)

	if err != nil {
		return nil, handlePqError(err) // Handles ErrNotFound
	}
	return e, nil
}

func (r *PostgresEmployeeRepository) ListAllEmployees() ([]*models.Employee, error) {
	query := `
		SELECT id, sin, first_name, last_name, address, phone, email, hotel_id, position, hire_date
		FROM employee
		ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, handlePqError(err)
	}
	defer rows.Close()

	employees := []*models.Employee{}
	for rows.Next() {
		e, err := scanEmployee(rows)
		if err != nil {
			return nil, handlePqError(err) // Handle row scanning error
		}
		employees = append(employees, e)
	}

	if err = rows.Err(); err != nil {
		return nil, handlePqError(err) // Check for errors during iteration
	}

	return employees, nil
}

func (r *PostgresEmployeeRepository) UpdateEmployee(emp *models.Employee) (*models.Employee, error) {
	if emp == nil {
		return nil, errors.New("Cannot update with a nil employee.")
	}
	if emp.ID <= 0 {
		return nil, errors.New("Invalid ID for employee update.")
	}
	if emp.SIN == "" || len(emp.SIN) != 9 || emp.FirstName == "" || emp.LastName == "" || emp.Email == "" || emp.HotelID <= 0 || emp.Position == "" || emp.HireDate.IsZero() {
		return nil, errors.New("Invalid employee data provided for update.")
	}

	query := `
		UPDATE employee
		SET sin = $1,
		    first_name = $2,
		    last_name = $3,
		    address = $4,
		    phone = $5,
		    email = $6,
		    hotel_id = $7,
		    position = $8,
		    hire_date = $9
		WHERE id = $10`

	result, err := r.db.Exec(query,
		emp.SIN,
		emp.FirstName,
		emp.LastName,
		emp.Address,
		emp.Phone,
		emp.Email,
		emp.HotelID,
		emp.Position,
		emp.HireDate,
		emp.ID,
	)
	if err != nil {
		// Checks unique sin/email, FK violation for hotel_id
		return nil, handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("Failed to check rows affected after employee update: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Return the original (potentially updated) employee pointer as per interface
	return emp, nil
}

// UpdateManager handles updating base employee data and manager-specific data.
func (r *PostgresEmployeeRepository) UpdateManager(mgr *models.Manager) error {
	if mgr == nil {
		return errors.New("Cannot update with a nil manager.")
	}
	if mgr.ID <= 0 { // Check the embedded employee ID
		return errors.New("Invalid ID for manager update.")
	}
	// Add validation for manager specific fields
	if mgr.Department == "" || mgr.AuthorizationLevel < 1 || mgr.AuthorizationLevel > 5 {
		return errors.New("Invalid manager-specific data provided for update.")
	}

	// Step 1: Update base employee record using the other method
	// We ignore the returned pointer here as the interface demands only error.
	_, err := r.UpdateEmployee(&mgr.Employee)
	if err != nil {
		// If employee update fails (e.g., not found, constraint violation), return error
		return err
	}

	// Step 2: Upsert manager-specific details
	managerQuery := `
		INSERT INTO manager (employee_id, department, authorization_level)
		VALUES ($1, $2, $3)
		ON CONFLICT (employee_id) DO UPDATE
		SET department = EXCLUDED.department,
		    authorization_level = EXCLUDED.authorization_level`

	_, err = r.db.Exec(managerQuery, mgr.ID, mgr.Department, mgr.AuthorizationLevel)
	if err != nil {
		// This could fail if somehow the employee_id FK constraint is violated, wouldn't count on it
		// if the UpdateEmployee above succeeded.
		// Or other DB errors.
		return handlePqError(err)
	}

	return nil // Success
}

func (r *PostgresEmployeeRepository) Delete(employeeID int) error {
	if employeeID <= 0 {
		return errors.New("Invalid employee ID for deletion.")
	}

	query := `DELETE FROM employee WHERE id = $1`

	result, err := r.db.Exec(query, employeeID)
	if err != nil {
		// CASCADE should handle manager, stay FKs.
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after employee delete: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
