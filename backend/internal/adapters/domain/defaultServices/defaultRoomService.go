// In internal/adapters/domain/defaultServices/default_room_service.go
package defaultServices

import (
	"errors"
	"fmt" // Import fmt for better error context if needed
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultRoomService struct {
	roomRepo ports.RoomRepository
}

func NewRoomService(repo ports.RoomRepository) ports.RoomService {
	return &DefaultRoomService{
		roomRepo: repo,
	}
}


func (s *DefaultRoomService) AddRoom(id, hotelId, capacity int, number, floor string, surfaceArea, price float64, telephone string,
	viewTypes map[models.ViewType]struct{}, roomType models.RoomType, isExtensible bool,
	amenities map[models.Amenity]struct{}, problems []models.Problem) (*models.Room, error) {
	room, err := models.NewRoom(id, hotelId, capacity, number, floor, surfaceArea, price, telephone, viewTypes, roomType, isExtensible, amenities, problems)
	if err != nil {
		// Return validation errors from the constructor
		return nil, fmt.Errorf("Validation failed for new room: %w", err)
	}

	// Delegate saving to the repository
	dbRoom, err := s.roomRepo.Save(room)
	if err != nil {
		// Return repository errors
		return nil, fmt.Errorf("Failed to save room: %w", err)
	}
	return dbRoom, nil
}

// UpdateRoom signature updated to include surfaceArea
func (s *DefaultRoomService) UpdateRoom(id, hotelId, capacity int, number, floor string, surfaceArea, price float64, telephone string,
	viewTypes map[models.ViewType]struct{}, roomType models.RoomType, isExtensible bool,
	amenities map[models.Amenity]struct{}, problems []models.Problem) (*models.Room, error) {

	// Fetch existing room first to ensure it exists (Update repo method also checks, but good practice here)
	_, err := s.roomRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, fmt.Errorf("Room with ID %d not found for update: %w", id, err)
		}
		return nil, fmt.Errorf("Failed to find room %d before update: %w", id, err)
	}

	// Validate parameters by attempting to create a new Room object with them
	updatedRoom, err := models.NewRoom(id, hotelId, capacity, number, floor, surfaceArea, price, telephone, viewTypes, roomType, isExtensible, amenities, problems)
	if err != nil {
		return nil, fmt.Errorf("Validation failed for updated room data: %w", err)
	}

	// Call repository update with the validated, complete room object
	err = s.roomRepo.Update(updatedRoom)
	if err != nil {
		// Return repository errors (like not found, constraints)
		return nil, fmt.Errorf("Failed to update room %d: %w", id, err)
	}

	// Return the representation of the room as it *should* be after the update
	return updatedRoom, nil
}

func (s *DefaultRoomService) DeleteRoom(id int) error {
	err := s.roomRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("Failed to delete room %d: %w", id, err)
	}
	return nil
}

func (s *DefaultRoomService) FindAvailableRooms(hotelID int, startDate time.Time, endDate time.Time) ([]*models.Room, error) {
	rooms, err := s.roomRepo.FindAvailableRooms(hotelID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("Failed to find available rooms for hotel %d: %w", hotelID, err)
	}
	return rooms, nil
}

func (s *DefaultRoomService) AssignRoomForReservation(reservation *models.Reservation) (int, error) {
	if reservation == nil {
		return 0, errors.New("Cannot assign room for nil reservation.")
	}

	rooms, err := s.FindAvailableRooms(reservation.HotelID, reservation.StartDate, reservation.EndDate)
	if err != nil {
		return 0, err
	}
	if len(rooms) == 0 {
		return 0, errors.New("No available rooms found matching the criteria.")
	}

	// Return the first available room's ID.
	return rooms[0].ID, nil
}

// Compile-time check
var _ ports.RoomService = (*DefaultRoomService)(nil)
