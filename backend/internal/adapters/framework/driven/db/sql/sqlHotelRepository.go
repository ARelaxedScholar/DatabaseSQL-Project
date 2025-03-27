package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type PostgresHotelRepository struct {
	db *sql.DB
}

func NewPostgresHotelRepository(db *sql.DB) (ports.HotelRepository, error) {
	if db == nil {
		return nil, errors.New("Db connection pool cannot be nil.")
	}
	return &PostgresHotelRepository{db: db}, nil
}

var _ ports.HotelRepository = (*PostgresHotelRepository)(nil)

func (r *PostgresHotelRepository) Save(hotel *models.Hotel) (*models.Hotel, error) {
	if hotel == nil {
		return nil, errors.New("Cannot save a nil hotel.")
	}
	if hotel.ChainID <= 0 || hotel.Rating < 1 || hotel.Rating > 5 || hotel.NumberOfRooms < 1 || hotel.Name == "" || hotel.Address == "" || hotel.Email == "" || hotel.Telephone == "" {
		return nil, errors.New("Invalid hotel data provided for save.")
	}

	query := `
		INSERT INTO Hotel (id_chaine, nom, adresse, email, telephone, rating, nombre_chambre)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_hotel`

	err := r.db.QueryRow(query,
		hotel.ChainID,
		hotel.Name,
		hotel.Address,
		hotel.Email,
		hotel.Telephone,
		hotel.Rating,
		hotel.NumberOfRooms,
	).Scan(&hotel.ID)

	if err != nil {
		return nil, handlePqError(err)
	}

	return hotel, nil
}

func (r *PostgresHotelRepository) FindByID(id int) (*models.Hotel, error) {
	if id <= 0 {
		return nil, errors.New("Invalid hotel ID provided.")
	}

	query := `
		SELECT id_hotel, id_chaine, nom, adresse, email, telephone, rating, nombre_chambre
		FROM Hotel
		WHERE id_hotel = $1`

	hotel := &models.Hotel{}
	var dbRating float64

	err := r.db.QueryRow(query, id).Scan(
		&hotel.ID,
		&hotel.ChainID,
		&hotel.Name,
		&hotel.Address,
		&hotel.Email,
		&hotel.Telephone,
		&dbRating,
		&hotel.NumberOfRooms,
	)

	if err != nil {
		return nil, handlePqError(err)
	}

	hotel.Rating = int(dbRating)

	return hotel, nil
}

func (r *PostgresHotelRepository) Update(hotel *models.Hotel) error {
	if hotel == nil {
		return errors.New("Cannot update with a nil hotel.")
	}
	if hotel.ID <= 0 {
		return errors.New("Invalid ID for hotel update.")
	}
	if hotel.ChainID <= 0 || hotel.Rating < 1 || hotel.Rating > 5 || hotel.NumberOfRooms < 1 || hotel.Name == "" || hotel.Address == "" || hotel.Email == "" || hotel.Telephone == "" {
		return errors.New("Invalid hotel data provided for update.")
	}

	query := `
		UPDATE Hotel
		SET id_chaine = $1,
		    nom = $2,
		    adresse = $3,
		    email = $4,
		    telephone = $5,
		    rating = $6,
		    nombre_chambre = $7
		WHERE id_hotel = $8`

	result, err := r.db.Exec(query,
		hotel.ChainID,
		hotel.Name,
		hotel.Address,
		hotel.Email,
		hotel.Telephone,
		hotel.Rating,
		hotel.NumberOfRooms,
		hotel.ID,
	)
	if err != nil {
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after update: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresHotelRepository) Delete(id int) error {
	if id <= 0 {
		return errors.New("Invalid hotel ID for deletion.")
	}

	query := `DELETE FROM Hotel WHERE id_hotel = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after delete: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
