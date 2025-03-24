package defaultAdminUseCases

import (
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultAdminHotelManagementUseCase struct {
	hotelService ports.HotelService
}

func NewAdminHotelManagementUseCase(hotelService ports.HotelService) ports.AdminHotelManagementUseCase {
	return &DefaultAdminHotelManagementUseCase{
		hotelService: hotelService,
	}
}

func (uc *DefaultAdminHotelManagementUseCase) AddHotel(input dto.HotelInput) (dto.HotelOutput, error) {
	hotel, err := uc.hotelService.AddHotel(
		input.ID,
		input.ChainID,
		input.Rating,
		input.NumberOfRooms,
		input.Name,
		input.Address,
		input.Email,
		input.Phone,
	)
	if err != nil {
		return dto.HotelOutput{}, err
	}
	return dto.HotelOutput{HotelID: hotel.ID}, nil
}

func (uc *DefaultAdminHotelManagementUseCase) UpdateHotel(input dto.HotelInput) (dto.HotelOutput, error) {
	hotel, err := uc.hotelService.UpdateHotel(
		input.ID,
		input.ChainID,
		input.Rating,
		input.NumberOfRooms,
		input.Name,
		input.Address,
		input.Email,
		input.Phone,
	)
	if err != nil {
		return dto.HotelOutput{}, err
	}
	return dto.HotelOutput{HotelID: hotel.ID}, nil
}

func (uc *DefaultAdminHotelManagementUseCase) DeleteHotel(hotelID int) error {
	return uc.hotelService.DeleteHotel(hotelID)
}
