package domain

import (
	"errors"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultHotelChainService struct {
	hotelChainRepo ports.HotelChainRepository
}

func NewHotelChainService(repo ports.HotelChainRepository) ports.HotelChainService {
	return &DefaultHotelChainService{
		hotelChainRepo: repo,
	}
}

func (s *DefaultHotelChainService) CreateHotelChain(id, numberOfHotel int, name, centralAddress, email, telephone string) (*models.HotelChain, error) {
	chain, err := models.NewHotelChain(id, numberOfHotel, name, centralAddress, email, telephone)
	if err != nil {
		return nil, err
	}
	dbHotelChain, err := s.hotelChainRepo.Save(chain)
	if err != nil {
		return nil, err
	}
	return dbHotelChain, nil
}

func (s *DefaultHotelChainService) UpdateHotelChain(id, numberOfHotel int, name, centralAddress, email, telephone string) (*models.HotelChain, error) {
	chain, err := s.hotelChainRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if chain == nil {
		return nil, errors.New("Hotel chain not found.")
	}
	chain = &models.HotelChain{
		ID:             id,
		NumberOfHotel:  numberOfHotel,
		Name:           name,
		CentralAddress: centralAddress,
		Email:          email,
		Telephone:      telephone,
	}
	if err = s.hotelChainRepo.Update(chain); err != nil {
		return nil, err
	}
	return chain, nil
}

func (s *DefaultHotelChainService) DeleteHotelChain(id int) error {
	chain, err := s.hotelChainRepo.FindByID(id)
	if err != nil {
		return err
	}
	if chain == nil {
		return errors.New("Hotel chain not found.")
	}
	return s.hotelChainRepo.Delete(id)
}
