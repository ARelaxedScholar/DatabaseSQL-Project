package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type PostgresReservationRepository struct {
	db *sql.DB
}

func NewPostgresReservationRepository(db *sql.DB) (ports.ReservationRepository, error) {
	if db == nil {
		return nil, errors.New("Db connection pool cannot be nil.")
	}
	return &PostgresReservationRepository{db: db}, nil
}

var _ ports.ReservationRepository = (*PostgresReservationRepository)(nil)

// Helper to scan reservation data
func scanReservation(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Reservation, error) {
	res := &models.Reservation{}
	var statusStr string
	var startDate, endDate, reservationDate time.Time
	var totalPrice float64 // Use float64 for DECIMAL

	// Ensure Scan order matches SELECT columns
	err := scanner.Scan(
		&res.ID,
		&res.ClientID,
		&res.RoomID,
		&res.HotelID, 
		&startDate,
		&endDate,
		&totalPrice,
		&reservationDate,
		&statusStr,
	)
	if err != nil {
		return nil, err // Let caller handle errors like ErrNotFound
	}

	// Assign time values
	res.StartDate = startDate
	res.EndDate = endDate
	res.ReservationDate = reservationDate
	res.TotalPrice = totalPrice

	// Convert status string from DB back to enum int value
	status, statusErr := models.ParseReservationStatus(statusStr)
	if statusErr != nil {
		// Handle error: data in DB doesn't match expected enum strings
		return nil, fmt.Errorf("Invalid reservation status '%s' found in database: %w", statusStr, statusErr)
	}
	res.Status = status

	return res, nil
}

func (r *PostgresReservationRepository) Save(res *models.Reservation) (*models.Reservation, error) {
	if res == nil {
		return nil, errors.New("Cannot save a nil reservation.")
	}
	// Basic checks
	if res.ClientID <= 0 || res.RoomID <= 0 || res.HotelID <= 0 || res.StartDate.IsZero() || res.EndDate.IsZero() || res.EndDate.Before(res.StartDate) || res.TotalPrice < 0 {
		return nil, errors.New("Invalid reservation data provided for save.")
	}
	// Convert Go enum status to string for DB
	statusStr := res.Status.String()
	if statusStr == "Invalid Status" { 
		return nil, errors.New("Invalid reservation status provided for save.")
	}

	query := `
		INSERT INTO reservation (client_id, room_id, hotel_id, start_date, end_date, total_price, reservation_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	// Use current time if ReservationDate is zero in the model
	resDate := res.ReservationDate
	if resDate.IsZero() {
		resDate = time.Now()
	}

	err := r.db.QueryRow(query,
		res.ClientID,
		res.RoomID,
		res.HotelID,
		res.StartDate,
		res.EndDate,
		res.TotalPrice,
		resDate,   
		statusStr, 
	).Scan(&res.ID)

	if err != nil {
		// Checks FK violations (client, room, hotel), date constraints, potentially unique reservation overlap
		return nil, handlePqError(err)
	}

	// Update the model's ReservationDate if it was defaulted
	res.ReservationDate = resDate

	return res, nil
}

func (r *PostgresReservationRepository) FindByID(id int) (*models.Reservation, error) {
	if id <= 0 {
		return nil, errors.New("Invalid reservation ID provided.")
	}

	query := `
		SELECT id, client_id, room_id, hotel_id, start_date, end_date, total_price, reservation_date, status
		FROM reservation
		WHERE id = $1`

	row := r.db.QueryRow(query, id)
	reservation, err := scanReservation(row)

	if err != nil {
		return nil, handlePqError(err) // Handles ErrNotFound
	}
	return reservation, nil
}

func (r *PostgresReservationRepository) GetByClient(clientID int) ([]*models.Reservation, error) {
	if clientID <= 0 {
		return nil, errors.New("Invalid client ID provided.")
	}

	query := `
		SELECT id, client_id, room_id, hotel_id, start_date, end_date, total_price, reservation_date, status
		FROM reservation
		WHERE client_id = $1
		ORDER BY start_date DESC` // Order with most recent on top

	rows, err := r.db.Query(query, clientID)
	if err != nil {
		return nil, handlePqError(err)
	}
	defer rows.Close()

	reservations := []*models.Reservation{}
	for rows.Next() {
		res, err := scanReservation(rows)
		if err != nil {
			return nil, handlePqError(err) 
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, handlePqError(err) 
	}

	// Returns an empty slice if no reservations are found, not ErrNotFound
	return reservations, nil
}

func (r *PostgresReservationRepository) Update(res *models.Reservation) error {
	if res == nil {
		return errors.New("Cannot update with a nil reservation.")
	}
	if res.ID <= 0 {
		return errors.New("Invalid ID for reservation update.")
	}
	// Basic checks...
	if res.ClientID <= 0 || res.RoomID <= 0 || res.HotelID <= 0 || res.StartDate.IsZero() || res.EndDate.IsZero() || res.EndDate.Before(res.StartDate) || res.TotalPrice < 0 {
		return errors.New("Invalid reservation data provided for update.")
	}
	// Convert Go enum status to string for DB
	statusStr := res.Status.String()
	if statusStr == "Invalid Status" {
		return errors.New("Invalid reservation status provided for update.")
	}

	query := `
		UPDATE reservation
		SET client_id = $1,
		    room_id = $2,
		    hotel_id = $3,
		    start_date = $4,
		    end_date = $5,
		    total_price = $6,
		    -- reservation_date is usually not updated, but status is
		    status = $7
		WHERE id = $8`

	result, err := r.db.Exec(query,
		res.ClientID,
		res.RoomID,
		res.HotelID,
		res.StartDate,
		res.EndDate,
		res.TotalPrice,
		statusStr, // Update status string
		res.ID,
	)
	if err != nil {
		// Checks FK violations, date constraints, etc.
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after reservation update: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil // no errors
}

func (r *PostgresReservationRepository) Delete(id int) error {
	if id <= 0 {
		return errors.New("Invalid reservation ID for deletion.")
	}

	// Note: stay.reservation_id is ON DELETE SET NULL
	query := `DELETE FROM reservation WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after reservation delete: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
