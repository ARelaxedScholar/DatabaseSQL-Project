package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type PostgresStayRepository struct {
	db *sql.DB
}

func NewPostgresStayRepository(db *sql.DB) (ports.StayRepository, error) {
	if db == nil {
		return nil, errors.New("Db connection pool cannot be nil.")
	}
	return &PostgresStayRepository{db: db}, nil
}

var _ ports.StayRepository = (*PostgresStayRepository)(nil)

// Helper to scan stay data, handling nullable integers
func scanStay(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Stay, error) {
	stay := &models.Stay{}
	// Use sql.Null types for nullable columns
	var reservationID sql.NullInt64
	var checkinEmpID sql.NullInt64
	var checkoutEmpID sql.NullInt64
	var arrivalDate, departureDate time.Time
	var finalPrice float64

	// Ensure Scan order matches SELECT columns
	err := scanner.Scan(
		&stay.ID,
		&stay.ClientID,
		&stay.RoomID,
		&reservationID, // Scan into sql.NullInt64
		&arrivalDate,
		&departureDate,
		&finalPrice,
		&stay.PaymentMethod,
		&checkinEmpID,  // Scan into sql.NullInt64
		&checkoutEmpID, // Scan into sql.NullInt64
		&stay.Comments,
	)
	if err != nil {
		return nil, err
	}

	// Assign time values
	stay.ArrivalDate = arrivalDate
	stay.DepartureDate = departureDate
	stay.FinalPrice = finalPrice

	// Convert sql.NullInt64 back to *int pointers
	if reservationID.Valid {
		id := int(reservationID.Int64) 
		stay.ReservationID = &id
	} else {
		stay.ReservationID = nil
	}

	if checkinEmpID.Valid {
		id := int(checkinEmpID.Int64)
		stay.CheckInEmployeeId = &id
	} else {
		stay.CheckInEmployeeId = nil
	}

	if checkoutEmpID.Valid {
		id := int(checkoutEmpID.Int64)
		stay.CheckOutEmployeeId = &id
	} else {
		stay.CheckOutEmployeeId = nil
	}

	return stay, nil
}

func (r *PostgresStayRepository) Save(stay *models.Stay) (*models.Stay, error) {
	if stay == nil {
		return nil, errors.New("Cannot save a nil stay.")
	}
	// Basic validation
	if stay.ClientID <= 0 || stay.RoomID <= 0 || stay.ArrivalDate.IsZero() || stay.DepartureDate.IsZero() || stay.DepartureDate.Before(stay.ArrivalDate) || stay.FinalPrice < 0 {
		return nil, errors.New("Invalid stay data provided for save.")
	}

	query := `
		INSERT INTO stay (client_id, room_id, reservation_id, arrival_date, departure_date, final_price, payment_method, checkin_employee_id, checkout_employee_id, comments)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	// Convert *int pointers to sql.NullInt64 for insertion
	var resID sql.NullInt64
	if stay.ReservationID != nil {
		resID = sql.NullInt64{Int64: int64(*stay.ReservationID), Valid: true}
	}
	var checkinID sql.NullInt64
	if stay.CheckInEmployeeId != nil {
		checkinID = sql.NullInt64{Int64: int64(*stay.CheckInEmployeeId), Valid: true}
	}
	var checkoutID sql.NullInt64
	if stay.CheckOutEmployeeId != nil {
		checkoutID = sql.NullInt64{Int64: int64(*stay.CheckOutEmployeeId), Valid: true}
	}

	err := r.db.QueryRow(query,
		stay.ClientID,
		stay.RoomID,
		resID, // Use sql.NullInt64
		stay.ArrivalDate,
		stay.DepartureDate,
		stay.FinalPrice,
		stay.PaymentMethod,
		checkinID,  // Use sql.NullInt64
		checkoutID, // Use sql.NullInt64
		stay.Comments,
	).Scan(&stay.ID)

	if err != nil {
		// Checks FK violations (client, room, reservation, employee), date constraints
		return nil, handlePqError(err)
	}

	return stay, nil
}

func (r *PostgresStayRepository) FindByID(id int) (*models.Stay, error) {
	if id <= 0 {
		return nil, errors.New("Invalid stay ID provided.")
	}

	query := `
		SELECT id, client_id, room_id, reservation_id, arrival_date, departure_date, final_price, payment_method, checkin_employee_id, checkout_employee_id, comments
		FROM stay
		WHERE id = $1`

	row := r.db.QueryRow(query, id)
	s, err := scanStay(row)

	if err != nil {
		return nil, handlePqError(err) // Handles ErrNotFound
	}
	return s, nil
}

func (r *PostgresStayRepository) Update(stay *models.Stay) error {
	if stay == nil {
		return errors.New("Cannot update with a nil stay.")
	}
	if stay.ID <= 0 {
		return errors.New("Invalid ID for stay update.")
	}
	// Basic validation
	if stay.ClientID <= 0 || stay.RoomID <= 0 || stay.ArrivalDate.IsZero() || stay.DepartureDate.IsZero() || stay.DepartureDate.Before(stay.ArrivalDate) || stay.FinalPrice < 0 {
		return errors.New("Invalid stay data provided for update.")
	}

	// Convert *int pointers to sql.NullInt64 for update
	var resID sql.NullInt64
	if stay.ReservationID != nil {
		resID = sql.NullInt64{Int64: int64(*stay.ReservationID), Valid: true}
	}
	var checkinID sql.NullInt64
	if stay.CheckInEmployeeId != nil {
		checkinID = sql.NullInt64{Int64: int64(*stay.CheckInEmployeeId), Valid: true}
	}
	var checkoutID sql.NullInt64
	if stay.CheckOutEmployeeId != nil {
		checkoutID = sql.NullInt64{Int64: int64(*stay.CheckOutEmployeeId), Valid: true}
	}

	query := `
		UPDATE stay
		SET client_id = $1,
		    room_id = $2,
		    reservation_id = $3,
		    arrival_date = $4,
		    departure_date = $5,
		    final_price = $6,
		    payment_method = $7,
		    checkin_employee_id = $8,
		    checkout_employee_id = $9,
		    comments = $10
		WHERE id = $11`

	result, err := r.db.Exec(query,
		stay.ClientID,
		stay.RoomID,
		resID, // Use sql.NullInt64
		stay.ArrivalDate,
		stay.DepartureDate,
		stay.FinalPrice,
		stay.PaymentMethod,
		checkinID,  // Use sql.NullInt64
		checkoutID, // Use sql.NullInt64
		stay.Comments,
		stay.ID,
	)
	if err != nil {
		// Checks FK violations, date constraints, etc.
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after stay update: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil // Interface expects only error for Update
}

func (r *PostgresStayRepository) Delete(id int) error {
	if id <= 0 {
		return errors.New("Invalid stay ID for deletion.")
	}

	query := `DELETE FROM stay WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after stay delete: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
