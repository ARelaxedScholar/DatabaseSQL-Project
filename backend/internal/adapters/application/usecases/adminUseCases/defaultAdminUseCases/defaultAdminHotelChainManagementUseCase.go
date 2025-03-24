package defaultAdminUseCases

import (
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultAdminHotelChainManagementUseCase struct {
	hotelChainService ports.HotelChainService
}

func NewAdminHotelChainManagementUseCase(hotelChainService ports.HotelChainService) ports.AdminHotelChainManagementUseCase {
	return &DefaultAdminHotelChainManagementUseCase{
		hotelChainService: hotelChainService,
	}
}

func (uc *DefaultAdminHotelChainManagementUseCase) AddHotelChain(input dto.HotelChainInput) (dto.HotelChainOutput, error) {
	chain, err := uc.hotelChainService.CreateHotelChain(
		input.ID,
		input.NumberOfHotels,
		input.Name,
		input.CentralAddress,
		input.Email,
		input.Telephone,
	)
	if err != nil {
		return dto.HotelChainOutput{}, err
	}
	return dto.HotelChainOutput{ChainID: chain.ID}, nil
}

func (uc *DefaultAdminHotelChainManagementUseCase) UpdateHotelChain(input dto.HotelChainInput) (dto.HotelChainOutput, error) {
	chain, err := uc.hotelChainService.UpdateHotelChain(
		input.ID,
		input.NumberOfHotels,
		input.Name,
		input.CentralAddress,
		input.Email,
		input.Telephone,
	)
	if err != nil {
		return dto.HotelChainOutput{}, err
	}
	return dto.HotelChainOutput{ChainID: chain.ID}, nil
}

func (uc *DefaultAdminHotelChainManagementUseCase) DeleteHotelChain(chainID int) error {
	return uc.hotelChainService.DeleteHotelChain(chainID)
}
