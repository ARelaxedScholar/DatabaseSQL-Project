package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type PostgresClientRepository struct {
	db *sql.DB
}

func NewPostgresClientRepository(db *sql.DB) (ports.ClientRepository, error) {
	if db == nil {
		return nil, errors.New("Db connection pool cannot be nil.")
	}
	return &PostgresClientRepository{db: db}, nil
}

var _ ports.ClientRepository = (*PostgresClientRepository)(nil)

// Helper to scan client data
func scanClient(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Client, error) {
	client := &models.Client{}
	var joinDate time.Time

	err := scanner.Scan(
		&client.ID,
		&client.SIN,
		&client.FirstName,
		&client.LastName,
		&client.Address,
		&client.Phone,
		&client.Email,
		&joinDate,
	)
	if err != nil {
		return nil, err // Let caller handle specific errors like ErrNotFound
	}
	client.JoinDate = joinDate
	return client, nil
}

func (r *PostgresClientRepository) Save(client *models.Client) (*models.Client, error) {
	if client == nil {
		return nil, errors.New("Cannot save a nil client.")
	}
	if client.SIN == "" || len(client.SIN) != 9 || client.FirstName == "" || client.LastName == "" || client.Email == "" || client.JoinDate.IsZero() {
		return nil, errors.New("Invalid client data provided for save.")
	}

	query := `
		INSERT INTO client (sin, first_name, last_name, address, phone, email, join_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.db.QueryRow(query,
		client.SIN,
		client.FirstName,
		client.LastName,
		client.Address,
		client.Phone,
		client.Email,
		client.JoinDate,
	).Scan(&client.ID)

	if err != nil {
		return nil, handlePqError(err) // Checks unique sin/email
	}

	return client, nil
}

func (r *PostgresClientRepository) FindByID(id int) (*models.Client, error) {
	if id <= 0 {
		return nil, errors.New("Invalid client ID provided.")
	}

	query := `
		SELECT id, sin, first_name, last_name, address, phone, email, join_date
		FROM client
		WHERE id = $1`

	row := r.db.QueryRow(query, id)
	c, err := scanClient(row)

	if err != nil {
		return nil, handlePqError(err) // Handles ErrNotFound
	}
	return c, nil
}

func (r *PostgresClientRepository) FindByEmail(email string) (*models.Client, error) {
	if email == "" {
		return nil, errors.New("Email cannot be empty for lookup.")
	}

	query := `
		SELECT id, sin, first_name, last_name, address, phone, email, join_date
		FROM client
		WHERE email = $1`

	row := r.db.QueryRow(query, email)
	c, err := scanClient(row)

	if err != nil {
		return nil, handlePqError(err) // Handles ErrNotFound
	}
	return c, nil
}

func (r *PostgresClientRepository) ListAllClients() ([]*models.Client, error) {
	query := `
		SELECT id, sin, first_name, last_name, address, phone, email, join_date
		FROM client
		ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, handlePqError(err)
	}
	defer rows.Close()

	clients := []*models.Client{}
	for rows.Next() {
		c, err := scanClient(rows)
		if err != nil {
			return nil, handlePqError(err) // Handle row scanning error
		}
		clients = append(clients, c)
	}

	if err = rows.Err(); err != nil {
		return nil, handlePqError(err) // Check for errors during iteration
	}

	return clients, nil
}

func (r *PostgresClientRepository) Update(client *models.Client) (*models.Client, error) {
	if client == nil {
		return nil, errors.New("Cannot update with a nil client.")
	}
	if client.ID <= 0 {
		return nil, errors.New("Invalid ID for client update.")
	}
	if client.SIN == "" || len(client.SIN) != 9 || client.FirstName == "" || client.LastName == "" || client.Email == "" || client.JoinDate.IsZero() {
		return nil, errors.New("Invalid client data provided for update.")
	}

	query := `
		UPDATE client
		SET sin = $1,
		    first_name = $2,
		    last_name = $3,
		    address = $4,
		    phone = $5,
		    email = $6,
		    join_date = $7
		WHERE id = $8`

	result, err := r.db.Exec(query,
		client.SIN,
		client.FirstName,
		client.LastName,
		client.Address,
		client.Phone,
		client.Email,
		client.JoinDate,
		client.ID,
	)
	if err != nil {
		return nil, handlePqError(err) // Checks unique sin/email
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Wrap this specific error for clarity
		return nil, fmt.Errorf("Failed to check rows affected after client update: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound // Use the defined sentinel error
	}

	// Return the original (potentially updated) client pointer as per interface
	return client, nil
}

func (r *PostgresClientRepository) Delete(id int) error {
	if id <= 0 {
		return errors.New("Invalid client ID for deletion.")
	}

	query := `DELETE FROM client WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		// CASCADE should handle FKs in reservation/stay,
		// but handlePqError could catch other issues.
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after client delete: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
