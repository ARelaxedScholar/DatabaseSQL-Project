package domain

import (
	"errors"
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

func (s *DefaultRoomService) AddRoom(id, hotelId, capacity, floor int, price float64, telephone string,
	viewTypes map[models.ViewType]struct{}, roomType models.RoomType, isExtensible bool,
	amenities map[models.Amenity]struct{}, problems []models.Problem) (*models.Room, error) {

	room, err := models.NewRoom(id, hotelId, capacity, floor, price, telephone, viewTypes, roomType, isExtensible, amenities, problems)
	if err != nil {
		return nil, err
	}
	dbRoom, err := s.roomRepo.Save(room)
	if err != nil {
		return nil, err
	}
	return dbRoom, nil
}

func (s *DefaultRoomService) UpdateRoom(id, hotelId, capacity, floor int, price float64, telephone string,
	viewTypes map[models.ViewType]struct{}, roomType models.RoomType, isExtensible bool,
	amenities map[models.Amenity]struct{}, problems []models.Problem) (*models.Room, error) {

	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("Room not found.")
	}
	// Recreate the room with new values; alternatively, update mutable fields if your model permits.
	room, err = models.NewRoom(id, hotelId, capacity, floor, price, telephone, viewTypes, roomType, isExtensible, amenities, problems)
	if err != nil {
		return nil, err
	}
	if err = s.roomRepo.Update(room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *DefaultRoomService) DeleteRoom(id int) error {
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("Room not found.")
	}
	return s.roomRepo.Delete(id)
}

func (s *DefaultRoomService) FindAvailableRooms(hotelID int, startDate time.Time, endDate time.Time) ([]*models.Room, error) {
	rooms, err := s.roomRepo.FindAvailableRooms(hotelID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *DefaultRoomService) AssignRoomForReservation(reservation *models.Reservation) (int, error) {
	rooms, err := s.FindAvailableRooms(reservation.HotelID, reservation.StartDate, reservation.EndDate)
	if err != nil {
		return 0, err
	}
	if len(rooms) == 0 {
		return 0, errors.New("no available rooms")
	}
	// For simplicity, we return the first available room's ID.
	return rooms[0].ID, nil
}
