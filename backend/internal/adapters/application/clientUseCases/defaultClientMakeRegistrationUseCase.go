package clientUseCases

import (
	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultClientMakeReservationUseCase struct {
	reservationService ports.ReservationService
}

func NewClientMakeReservationUseCase(reservationService ports.ReservationService) ports.ClientMakeReservationUseCase {
	return &DefaultClientMakeReservationUseCase{
		reservationService: reservationService,
	}
}

func (uc *DefaultClientMakeReservationUseCase) MakeReservation(input dto.ReservationInput) (dto.ReservationOutput, error) {
	reservation, err := uc.reservationService.CreateReservation(
		0, // pass a default value, let the db deal with it
		input.ClientID,
		input.RoomID,
		input.StartDate,
		input.EndDate,
		input.ReservationDate,
		input.TotalPrice,
		models.ReservationStatus(input.Status),
	)
	if err != nil {
		return dto.ReservationOutput{}, err
	}

	return dto.ReservationOutput{
		ReservationID: reservation.ID,
		ClientID:      reservation.ClientID,
		RoomID:        reservation.RoomID,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		TotalPrice:    reservation.TotalPrice,
		Status:        int(reservation.Status),
	}, nil
}
