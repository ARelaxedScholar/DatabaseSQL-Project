package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sql-project-backend/internal/models" // Use models package for ErrNotFound
	"github.com/sql-project-backend/internal/ports"

	"github.com/lib/pq"
)

type PostgresRoomRepository struct {
	db *sql.DB
}

func NewPostgresRoomRepository(db *sql.DB) (ports.RoomRepository, error) {
	if db == nil {
		return nil, errors.New("Db connection pool cannot be nil.")
	}
	return &PostgresRoomRepository{db: db}, nil
}

var _ ports.RoomRepository = (*PostgresRoomRepository)(nil)

// --- Helper functions for M2M relationships ---

func syncRoomViewTypes(tx *sql.Tx, roomID int, viewTypes map[models.ViewType]struct{}) error {
	_, err := tx.Exec(`DELETE FROM room_view_type WHERE room_id = $1`, roomID)
	if err != nil {
		return fmt.Errorf("Failed to clear existing view types for room %d: %w", roomID, err)
	}
	if len(viewTypes) == 0 {
		return nil
	}
	stmt, err := tx.Prepare(pq.CopyIn("room_view_type", "room_id", "view_type_id"))
	if err != nil {
		return fmt.Errorf("Failed to prepare view type copy statement: %w", err)
	}
	defer stmt.Close()
	viewTypeIDs := make(map[string]int)
	for vtEnum := range viewTypes {
		vtName := vtEnum.String()
		if _, ok := viewTypeIDs[vtName]; !ok {
			var vtID int
			err = tx.QueryRow(`SELECT id FROM view_type WHERE name = $1`, vtName).Scan(&vtID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("View type '%s' not found in lookup table.", vtName)
				}
				return fmt.Errorf("Failed to query view type ID for '%s': %w", vtName, err)
			}
			viewTypeIDs[vtName] = vtID
		}
		_, err = stmt.Exec(int64(roomID), int64(viewTypeIDs[vtName]))
		if err != nil {
			return fmt.Errorf("Failed to stage view type copy for room %d, view type %d: %w", roomID, viewTypeIDs[vtName], err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("Failed to execute view type copy for room %d: %w", roomID, err)
	}
	return nil
}

func syncRoomAmenities(tx *sql.Tx, roomID int, amenities map[models.Amenity]struct{}) error {
	_, err := tx.Exec(`DELETE FROM room_amenity WHERE room_id = $1`, roomID)
	if err != nil {
		return fmt.Errorf("Failed to clear existing amenities for room %d: %w", roomID, err)
	}
	if len(amenities) == 0 {
		return nil
	}
	stmt, err := tx.Prepare(pq.CopyIn("room_amenity", "room_id", "amenity_id"))
	if err != nil {
		return fmt.Errorf("Failed to prepare amenity copy statement: %w", err)
	}
	defer stmt.Close()
	amenityIDs := make(map[string]int)
	for amEnum := range amenities {
		amName := amEnum.String()
		if _, ok := amenityIDs[amName]; !ok {
			var amID int
			err = tx.QueryRow(`SELECT id FROM amenity WHERE name = $1`, amName).Scan(&amID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("Amenity '%s' not found in lookup table.", amName)
				}
				return fmt.Errorf("Failed to query amenity ID for '%s': %w", amName, err)
			}
			amenityIDs[amName] = amID
		}
		_, err = stmt.Exec(int64(roomID), int64(amenityIDs[amName]))
		if err != nil {
			return fmt.Errorf("Failed to stage amenity copy for room %d, amenity %d: %w", roomID, amenityIDs[amName], err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("Failed to execute amenity copy for room %d: %w", roomID, err)
	}
	return nil
}

func syncRoomProblems(tx *sql.Tx, roomID int, problems []models.Problem) error {
	_, err := tx.Exec(`DELETE FROM room_problem WHERE room_id = $1`, roomID)
	if err != nil {
		return fmt.Errorf("Failed to clear existing problems for room %d: %w", roomID, err)
	}
	if len(problems) == 0 {
		return nil
	}
	stmt, err := tx.Prepare(pq.CopyIn("room_problem", "room_id", "description", "signaled_when", "severity", "is_resolved", "resolution_date"))
	if err != nil {
		return fmt.Errorf("Failed to prepare problem copy statement: %w", err)
	}
	defer stmt.Close()
	for _, p := range problems {
		severityStr := p.Severity.String()
		if severityStr == "Invalid Severity" {
			return fmt.Errorf("Invalid problem severity provided for room %d.", roomID)
		}
		var resolutionDate pq.NullTime
		if p.IsResolved && !p.ResolutionDate.IsZero() {
			resolutionDate = pq.NullTime{Time: p.ResolutionDate, Valid: true}
		}
		_, err = stmt.Exec(int64(roomID), p.Description, p.SignaledWhen, severityStr, p.IsResolved, resolutionDate)
		if err != nil {
			return fmt.Errorf("Failed to stage problem copy for room %d: %w", roomID, err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("Failed to execute problem copy for room %d: %w", roomID, err)
	}
	return nil
}

// --- Repository Methods ---

func (r *PostgresRoomRepository) Save(room *models.Room) (*models.Room, error) {
	if room == nil {
		return nil, errors.New("Cannot save a nil room.")
	}
	if room.HotelID <= 0 || room.Capacity < 1 || room.Price < 0 || room.SurfaceArea <= 0 || room.RoomType == 0 || room.Number == "" || room.Floor == "" {
		return nil, errors.New("Invalid room data provided for save.")
	}

	var roomTypeID int
	err := r.db.QueryRow(`SELECT id FROM room_type WHERE name = $1`, room.RoomType.String()).Scan(&roomTypeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Room type '%s' not found in lookup table.", room.RoomType.String())
		}
		return nil, fmt.Errorf("Failed to query room type ID: %w.", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Failed to begin transaction: %w.", err)
	}
	defer tx.Rollback()

	roomQuery := `
		INSERT INTO room (hotel_id, room_type_id, number, floor, capacity, surface_area, price, telephone, is_extensible)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err = tx.QueryRow(roomQuery,
		room.HotelID, roomTypeID, room.Number, room.Floor, room.Capacity,
		room.SurfaceArea, room.Price, room.Telephone, room.IsExtensible, // Added surface_area
	).Scan(&room.ID)
	if err != nil {
		return nil, handlePqError(err)
	}

	if err = syncRoomViewTypes(tx, room.ID, room.ViewTypes); err != nil {
		return nil, err
	}
	if err = syncRoomAmenities(tx, room.ID, room.Amenities); err != nil {
		return nil, err
	}
	if err = syncRoomProblems(tx, room.ID, room.Problems); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("Failed to commit transaction: %w.", err)
	}
	return room, nil
}

// fetchRoomsWithDetails is a helper used by FindByID, FindAvailableRooms, SearchRooms
func (r *PostgresRoomRepository) fetchRoomsWithDetails(roomIDs []int) ([]*models.Room, error) {
	if len(roomIDs) == 0 {
		return []*models.Room{}, nil
	}
	idArray := pq.Array(roomIDs)

	queryMain := `
        SELECT
            r.id, r.hotel_id, r.number, r.floor, r.capacity, r.surface_area, r.price, r.telephone, r.is_extensible,
            rt.name as room_type_name
        FROM room r
        JOIN room_type rt ON r.room_type_id = rt.id
        WHERE r.id = ANY($1) ORDER BY r.id`
	rowsMain, err := r.db.Query(queryMain, idArray)
	if err != nil {
		return nil, handlePqError(fmt.Errorf("Failed to query main room details: %w", err))
	}
	defer rowsMain.Close()
	roomsMap := make(map[int]*models.Room, len(roomIDs))
	processedOrder := []int{}
	for rowsMain.Next() {
		room := &models.Room{ViewTypes: make(map[models.ViewType]struct{}), Amenities: make(map[models.Amenity]struct{}), Problems: []models.Problem{}}
		var roomTypeName string
		err = rowsMain.Scan(&room.ID, &room.HotelID, &room.Number, &room.Floor, &room.Capacity, &room.SurfaceArea, // Added surface_area
			&room.Price, &room.Telephone, &room.IsExtensible, &roomTypeName)
		if err != nil {
			return nil, handlePqError(fmt.Errorf("Failed to scan main room details: %w", err))
		}
		rtEnum, parseErr := models.ParseRoomType(roomTypeName)
		if parseErr != nil {
			return nil, fmt.Errorf("Failed to parse room type '%s': %w", roomTypeName, parseErr)
		}
		room.RoomType = rtEnum
		roomsMap[room.ID] = room
		processedOrder = append(processedOrder, room.ID)
	}
	if err = rowsMain.Err(); err != nil {
		return nil, handlePqError(fmt.Errorf("Error iterating main room details: %w", err))
	}
	if len(roomsMap) == 0 {
		return []*models.Room{}, nil
	}

	// Fetch View Types
	rowsVt, errVt := r.db.Query(`SELECT rvt.room_id, vt.name FROM view_type vt JOIN room_view_type rvt ON vt.id = rvt.view_type_id WHERE rvt.room_id = ANY($1)`, idArray)
	if errVt != nil {
		return nil, handlePqError(fmt.Errorf("Failed to query view types: %w", errVt))
	}
	defer rowsVt.Close()
	for rowsVt.Next() {
		var roomID int
		var vtName string
		if err := rowsVt.Scan(&roomID, &vtName); err != nil {
			return nil, handlePqError(fmt.Errorf("Failed to scan view type: %w", err))
		}
		if room, ok := roomsMap[roomID]; ok {
			vtEnum, _ := models.ParseViewType(vtName)
			if vtEnum != 0 {
				room.ViewTypes[vtEnum] = struct{}{}
			}
		}
	}
	if err = rowsVt.Err(); err != nil {
		return nil, handlePqError(fmt.Errorf("Error iterating view types: %w", err))
	}

	// Fetch Amenities
	rowsAm, errAm := r.db.Query(`SELECT ra.room_id, a.name FROM amenity a JOIN room_amenity ra ON a.id = ra.amenity_id WHERE ra.room_id = ANY($1)`, idArray)
	if errAm != nil {
		return nil, handlePqError(fmt.Errorf("Failed to query amenities: %w", errAm))
	}
	defer rowsAm.Close()
	for rowsAm.Next() {
		var roomID int
		var amName string
		if err := rowsAm.Scan(&roomID, &amName); err != nil {
			return nil, handlePqError(fmt.Errorf("Failed to scan amenity: %w", err))
		}
		if room, ok := roomsMap[roomID]; ok {
			amEnum, _ := models.ParseAmenity(amName)
			if amEnum != 0 {
				room.Amenities[amEnum] = struct{}{}
			}
		}
	}
	if err = rowsAm.Err(); err != nil {
		return nil, handlePqError(fmt.Errorf("Error iterating amenities: %w", err))
	}

	// Fetch Problems
	rowsPr, errPr := r.db.Query(`SELECT room_id, id, description, signaled_when, severity, is_resolved, resolution_date FROM room_problem WHERE room_id = ANY($1) ORDER BY room_id, signaled_when DESC`, idArray)
	if errPr != nil {
		return nil, handlePqError(fmt.Errorf("Failed to query problems: %w", errPr))
	}
	defer rowsPr.Close()
	for rowsPr.Next() {
		var roomID int
		prob := models.Problem{}
		var severityStr string
		var signaled, resolution sql.NullTime
		err := rowsPr.Scan(&roomID, &prob.ID, &prob.Description, &signaled, &severityStr, &prob.IsResolved, &resolution)
		if err != nil {
			return nil, handlePqError(fmt.Errorf("Failed to scan problem: %w", err))
		}
		if room, ok := roomsMap[roomID]; ok {
			if signaled.Valid {
				prob.SignaledWhen = signaled.Time
			}
			if resolution.Valid {
				prob.ResolutionDate = resolution.Time
			} else {
				prob.ResolutionDate = time.Time{}
			}
			sevEnum, _ := models.ParseProblemSeverity(severityStr)
			if sevEnum != 0 {
				prob.Severity = sevEnum
			}
			room.Problems = append(room.Problems, prob)
		}
	}
	if err = rowsPr.Err(); err != nil {
		return nil, handlePqError(fmt.Errorf("Error iterating problems: %w", err))
	}

	finalRooms := make([]*models.Room, len(processedOrder))
	for i, id := range processedOrder {
		finalRooms[i] = roomsMap[id]
	}
	return finalRooms, nil
}

func (r *PostgresRoomRepository) FindByID(id int) (*models.Room, error) {
	if id <= 0 {
		return nil, errors.New("Invalid room ID provided.")
	}
	results, err := r.fetchRoomsWithDetails([]int{id})
	if err != nil {
		return nil, err
	} // Error already handled in fetch
	if len(results) == 0 {
		return nil, models.ErrNotFound
	} // Use models.ErrNotFound
	return results[0], nil
}

func (r *PostgresRoomRepository) Update(room *models.Room) error {
	if room == nil {
		return errors.New("Cannot update with a nil room.")
	}
	if room.ID <= 0 {
		return errors.New("Invalid ID for room update.")
	}
	if room.HotelID <= 0 || room.Capacity < 1 || room.Price < 0 || room.SurfaceArea <= 0 || room.RoomType == 0 || room.Number == "" || room.Floor == "" {
		return errors.New("Invalid room data provided for update.")
	}
	var roomTypeID int
	err := r.db.QueryRow(`SELECT id FROM room_type WHERE name = $1`, room.RoomType.String()).Scan(&roomTypeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Room type '%s' not found in lookup table.", room.RoomType.String())
		}
		return fmt.Errorf("Failed to query room type ID: %w.", err)
	}
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w.", err)
	}
	defer tx.Rollback()

	roomQuery := `
		UPDATE room SET hotel_id = $1, room_type_id = $2, number = $3, floor = $4,
		    capacity = $5, surface_area = $6, price = $7, telephone = $8, is_extensible = $9
		WHERE id = $10` // Added surface_area

	result, err := tx.Exec(roomQuery,
		room.HotelID, roomTypeID, room.Number, room.Floor, room.Capacity,
		room.SurfaceArea, room.Price, room.Telephone, room.IsExtensible, room.ID, // Added surfaceArea
	)
	if err != nil {
		return handlePqError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after room update: %w.", err)
	}
	if rowsAffected == 0 {
		return models.ErrNotFound
	} // Use models.ErrNotFound

	if err = syncRoomViewTypes(tx, room.ID, room.ViewTypes); err != nil {
		return err
	}
	if err = syncRoomAmenities(tx, room.ID, room.Amenities); err != nil {
		return err
	}
	if err = syncRoomProblems(tx, room.ID, room.Problems); err != nil {
		return err
	} // Assuming overwrite for problems

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %w.", err)
	}
	return nil
}

func (r *PostgresRoomRepository) Delete(id int) error {
	if id <= 0 {
		return errors.New("Invalid room ID for deletion.")
	}
	query := `DELETE FROM room WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return handlePqError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check rows affected after room delete: %w.", err)
	}
	if rowsAffected == 0 {
		return models.ErrNotFound
	} // Use models.ErrNotFound
	return nil
}

func (r *PostgresRoomRepository) FindAvailableRooms(hotelID int, startDate time.Time, endDate time.Time) ([]*models.Room, error) {
	if hotelID <= 0 {
		return nil, errors.New("Invalid hotel ID provided.")
	}
	if startDate.IsZero() || endDate.IsZero() || !endDate.After(startDate) {
		return nil, errors.New("Invalid start or end date provided.")
	}
	queryIDs := ` SELECT r.id FROM room r WHERE r.hotel_id = $1
          AND NOT EXISTS ( SELECT 1 FROM reservation res WHERE res.room_id = r.id AND res.status != 3 AND res.start_date < $2 AND res.end_date > $3 )
          AND NOT EXISTS ( SELECT 1 FROM stay s WHERE s.room_id = r.id AND s.arrival_date < $2 AND s.departure_date > $3 )
        ORDER BY r.id `
	rowsIDs, err := r.db.Query(queryIDs, hotelID, endDate, startDate)
	if err != nil {
		return nil, handlePqError(fmt.Errorf("Failed to query available room IDs: %w", err))
	}
	defer rowsIDs.Close()
	var availableRoomIDs []int
	for rowsIDs.Next() {
		var id int
		if err := rowsIDs.Scan(&id); err != nil {
			return nil, handlePqError(fmt.Errorf("Failed to scan available room ID: %w", err))
		}
		availableRoomIDs = append(availableRoomIDs, id)
	}
	if err = rowsIDs.Err(); err != nil {
		return nil, handlePqError(fmt.Errorf("Error iterating available room IDs: %w", err))
	}
	if len(availableRoomIDs) == 0 {
		return []*models.Room{}, nil
	}
	return r.fetchRoomsWithDetails(availableRoomIDs)
}

func (r *PostgresRoomRepository) SearchRooms(startDate time.Time, endDate time.Time, capacity int, priceMin, priceMax float64, hotelChainID int, roomType models.RoomType) ([]*models.Room, error) {
	var queryFilter strings.Builder
	args := []interface{}{}
	argID := 1

	queryFilter.WriteString(" WHERE 1=1 ")

	if hotelChainID > 0 {
		queryFilter.WriteString(fmt.Sprintf("AND h.hotel_chain_id = $%d ", argID))
		args = append(args, hotelChainID)
		argID++
	}
	if capacity > 0 {
		queryFilter.WriteString(fmt.Sprintf("AND r.capacity >= $%d ", argID))
		args = append(args, capacity)
		argID++
	}
	if priceMin > 0 {
		queryFilter.WriteString(fmt.Sprintf("AND r.price >= $%d ", argID))
		args = append(args, priceMin)
		argID++
	}
	if priceMax > priceMin {
		queryFilter.WriteString(fmt.Sprintf("AND r.price <= $%d ", argID))
		args = append(args, priceMax)
		argID++
	}
	if roomType != 0 {
		rtName := roomType.String()
		if rtName == "Invalid Room Type" {
			return nil, errors.New("Invalid room type.")
		}
		queryFilter.WriteString(fmt.Sprintf("AND rt.name = $%d ", argID))
		args = append(args, rtName)
		argID++
	}
	// Only add date availability filtering if both dates are provided.
	if !startDate.IsZero() && !endDate.IsZero() {
		// Append endDate and startDate as parameters.
		endDateArgIdx := argID
		args = append(args, endDate)
		argID++
		startDateArgIdx := argID
		args = append(args, startDate)
		argID++

		// Add conditions to exclude rooms with overlapping reservations
		queryFilter.WriteString(fmt.Sprintf(
			"AND NOT EXISTS ( SELECT 1 FROM reservation res WHERE res.room_id = r.id AND res.status != 3 AND res.start_date < $%d AND res.end_date > $%d ) ",
			endDateArgIdx, startDateArgIdx,
		))
		// Add conditions to exclude rooms with overlapping stays
		queryFilter.WriteString(fmt.Sprintf(
			"AND NOT EXISTS ( SELECT 1 FROM stay s WHERE s.room_id = r.id AND s.arrival_date < $%d AND s.departure_date > $%d ) ",
			endDateArgIdx, startDateArgIdx,
		))
	}

	queryIDs := fmt.Sprintf(`
        SELECT r.id 
        FROM room r 
        JOIN room_type rt ON r.room_type_id = rt.id 
        JOIN hotel h ON r.hotel_id = h.id 
        %s
        ORDER BY r.price, r.id
    `, queryFilter.String())

	rowsIDs, err := r.db.Query(queryIDs, args...)
	if err != nil {
		return nil, handlePqError(fmt.Errorf("Failed to query searched room IDs: %w", err))
	}
	defer rowsIDs.Close()

	var availableRoomIDs []int
	for rowsIDs.Next() {
		var id int
		if err := rowsIDs.Scan(&id); err != nil {
			return nil, handlePqError(fmt.Errorf("Failed to scan searched room ID: %w", err))
		}
		availableRoomIDs = append(availableRoomIDs, id)
	}
	if err = rowsIDs.Err(); err != nil {
		return nil, handlePqError(fmt.Errorf("Error iterating searched room IDs: %w", err))
	}
	return r.fetchRoomsWithDetails(availableRoomIDs)
}
