package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/sql-project-backend/internal/ports"
)

type PostgresQueryRepository struct {
	db *sql.DB
}

func NewPostgresQueryRepository(db *sql.DB) (ports.QueryRepository, error) {
	if db == nil {
		return nil, errors.New("Db connection pool cannot be nil.")
	}
	return &PostgresQueryRepository{db: db}, nil
}

var _ ports.QueryRepository = (*PostgresQueryRepository)(nil)

// GetHotelRoomCapacity counts the total number of rooms associated with a specific hotel ID.
func (r *PostgresQueryRepository) GetHotelRoomCapacity(hotelId int) (int, error) {
	if hotelId <= 0 {
		return 0, errors.New("Invalid hotel ID provided.")
	}

	// Query to count rooms for the given hotel ID.
	query := `SELECT room_count FROM room_per_hotel rph WHERE rph.hotel_id = $1`

	var count int
	var hotelExists bool

	err := r.db.QueryRow(query, hotelId).Scan(&count)

	hotelExists = count != 0 // an hotel cannot have no rooms

	if err != nil {
		wrappedErr := handlePqError(err) // Use helper for consistency
		return 0, fmt.Errorf("Failed to query hotel room capacity for hotel ID %d: %w", hotelId, wrappedErr)
	}

	// Return
	if !hotelExists {
		return 0, fmt.Errorf("hotel with ID %d does not exist", hotelId)
	}
	return count, nil
}

// GetAvailableRoomsByZone counts the total number of rooms grouped by the city of the hotel.
func (r *PostgresQueryRepository) GetAvailableRoomsByZone() (map[string]int, error) {

	query := `
        SELECT * FROM rooms_by_zones;
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, handlePqError(fmt.Errorf("Failed to query total rooms by zone: %w", err))
	}
	defer rows.Close()

	results := make(map[string]int)
	for rows.Next() {
		var zoneName string
		var count int
		if err := rows.Scan(&zoneName, &count); err != nil {
			return nil, handlePqError(fmt.Errorf("Failed to scan zone result: %w", err))
		}
		results[zoneName] = count
	}
	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, handlePqError(fmt.Errorf("Error iterating zone results: %w", err))
	}

	// Return the map (which will be empty if no rooms/hotels with cities exist)
	return results, nil
}
