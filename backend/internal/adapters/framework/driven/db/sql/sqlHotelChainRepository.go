package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"

	"github.com/lib/pq"
)

var (
	ErrDuplicateEntry      = errors.New("Database constraint violation: duplicate entry.")
	ErrForeignKeyViolation = errors.New("Database constraint violation: foreign key.")
	ErrNotFound            = sql.ErrNoRows
)

type PostgresHotelChainRepository struct {
	db *sql.DB
}

func NewPostgresHotelChainRepository(db *sql.DB) (ports.HotelChainRepository, error) {
	if db == nil {
		return nil, errors.New("Db connection pool cannot be nil.")
	}
	return &PostgresHotelChainRepository{db: db}, nil
}

var _ ports.HotelChainRepository = (*PostgresHotelChainRepository)(nil)

func handlePqError(err error) error {
	if err == nil {
		return nil
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s.", ErrDuplicateEntry, pqErr.Constraint)
		case "23503":
			return fmt.Errorf("%w: %s.", ErrForeignKeyViolation, pqErr.Constraint)
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	// Return error with more info
	return fmt.Errorf("Database operation failed: %w", err)
}

func (r *PostgresHotelChainRepository) Save(chain *models.HotelChain) (*models.HotelChain, error) {
	if chain == nil {
		return nil, errors.New("Cannot save a nil hotel chain.")
	}
	if chain.Name == "" || chain.CentralAddress == "" || chain.Email == "" || chain.Telephone == "" || chain.NumberOfHotel < 0 {
		return nil, errors.New("Invalid hotel chain data provided for save.")
	}

	query := `
		INSERT INTO Chaine_Hotel (nom, adresse_centrale, nombre_hotel, email, telephone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_chaine`

	err := r.db.QueryRow(query,
		chain.Name,
		chain.CentralAddress,
		chain.NumberOfHotel,
		chain.Email,
		chain.Telephone,
	).Scan(&chain.ID)

	if err != nil {
		return nil, handlePqError(err)
	}

	return chain, nil
}

func (r *PostgresHotelChainRepository) FindByID(id int) (*models.HotelChain, error) {
	if id <= 0 {
		return nil, errors.New("Invalid hotel chain ID provided.")
	}

	query := `
		SELECT id_chaine, nom, adresse_centrale, nombre_hotel, email, telephone
		FROM Chaine_Hotel
		WHERE id_chaine = $1`

	chain := &models.HotelChain{}

	err := r.db.QueryRow(query, id).Scan(
		&chain.ID,
		&chain.Name,
		&chain.CentralAddress,
		&chain.NumberOfHotel,
		&chain.Email,
		&chain.Telephone,
	)

	if err != nil {
		return nil, handlePqError(err)
	}

	return chain, nil
}

func (r *PostgresHotelChainRepository) Update(chain *models.HotelChain) error {
	if chain == nil {
		return errors.New("Cannot update with a nil hotel chain.")
	}
	if chain.ID <= 0 {
		return errors.New("Invalid ID for hotel chain update.")
	}
	if chain.Name == "" || chain.CentralAddress == "" || chain.Email == "" || chain.Telephone == "" || chain.NumberOfHotel < 0 {
		return errors.New("Invalid hotel chain data provided for update.")
	}

	query := `
		UPDATE Chaine_Hotel
		SET nom = $1,
		    adresse_centrale = $2,
		    nombre_hotel = $3,
		    email = $4,
		    telephone = $5
		WHERE id_chaine = $6`

	result, err := r.db.Exec(query,
		chain.Name,
		chain.CentralAddress,
		chain.NumberOfHotel,
		chain.Email,
		chain.Telephone,
		chain.ID,
	)
	if err != nil {
		return handlePqError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Format error message correctly
		return fmt.Errorf("Failed to check rows affected after update: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresHotelChainRepository) Delete(id int) error {
	if id <= 0 {
		return errors.New("Invalid hotel chain ID for deletion.")
	}

	query := `DELETE FROM Chaine_Hotel WHERE id_chaine = $1`

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
