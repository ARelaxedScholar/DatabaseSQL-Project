package defaultServices

import (
	"errors"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultHotelService struct {
	hotelRepo ports.HotelRepository
}

func NewHotelService(repo ports.HotelRepository) ports.HotelService {
	return &DefaultHotelService{
		hotelRepo: repo,
	}
}

func (s *DefaultHotelService) AddHotel(id, chainId, rating, numberOfRooms int, name, address, email, phone string) (*models.Hotel, error) {
	hotel, err := models.NewHotel(id, chainId, rating, numberOfRooms, name, address, email, phone)
	if err != nil {
		return nil, err
	}
	dbHotel, err := s.hotelRepo.Save(hotel)
	if err != nil {
		return nil, err
	}
	return dbHotel, nil
}

func (s *DefaultHotelService) UpdateHotel(id, chainId, rating, numberOfRooms int, name, address, email, phone string) (*models.Hotel, error) {
	hotel, err := s.hotelRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if hotel == nil {
		return nil, errors.New("Hotel not found.")
	}
	hotel, err = models.NewHotel(id, chainId, rating, numberOfRooms, name, address, email, phone)
	if err != nil {
		return nil, err
	}
	if err = s.hotelRepo.Update(hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}

func (s *DefaultHotelService) DeleteHotel(id int) error {
	hotel, err := s.hotelRepo.FindByID(id)
	if err != nil {
		return err
	}
	if hotel == nil {
		return errors.New("Hotel not found.")
	}
	return s.hotelRepo.Delete(id)
}
