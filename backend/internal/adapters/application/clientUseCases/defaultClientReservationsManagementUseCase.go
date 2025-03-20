package clientUseCases

import (
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultClientReservationsManagementUseCase struct {
	reservationService ports.ReservationService
}

func NewClientReservationsManagementUseCase(reservationService ports.ReservationService) ports.ClientReservationsManagementUseCase {
	return &DefaultClientReservationsManagementUseCase{
		reservationService: reservationService,
	}
}

func (uc *DefaultClientReservationsManagementUseCase) ViewReservations(clientID int) ([]dto.ReservationOutput, error) {
	reservations, err := uc.reservationService.GetReservationsByClient(clientID)
	if err != nil {
		return nil, err
	}

	outputs := make([]dto.ReservationOutput, 0, len(reservations))
	for _, r := range reservations {
		outputs = append(outputs, dto.ReservationOutput{
			ReservationID: r.ID,
			ClientID:      r.ClientID,
			RoomID:        r.RoomID,
			StartDate:     r.StartDate,
			EndDate:       r.EndDate,
			TotalPrice:    r.TotalPrice,
			Status:        int(r.Status), //underlying type is int
		})
	}

	return outputs, nil
}

func (uc *DefaultClientReservationsManagementUseCase) CancelReservation(reservationID int) error {
	return uc.reservationService.CancelReservation(reservationID)
}
