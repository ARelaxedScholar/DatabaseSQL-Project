package employeeUseCases

import (
	"errors"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultEmployeeCheckInUseCase struct {
	stayService     ports.StayService
	reservationRepo ports.ReservationRepository
	roomRepo        ports.RoomRepository
}

func NewEmployeeCheckInUseCase(
	stayService ports.StayService,
	reservationRepo ports.ReservationRepository,
	roomRepo ports.RoomRepository,
) ports.EmployeeCheckInUseCase {
	return &DefaultEmployeeCheckInUseCase{
		stayService:     stayService,
		reservationRepo: reservationRepo,
		roomRepo:        roomRepo,
	}
}

func (uc *DefaultEmployeeCheckInUseCase) CheckIn(input dto.CheckInInput) (dto.CheckInOutput, error) {
	reservation, err := uc.reservationRepo.FindByID(input.ReservationID)
	if err != nil {
		return dto.CheckInOutput{}, err
	}
	if reservation == nil {
		return dto.CheckInOutput{}, errors.New("reservation not found")
	}

	roomID := reservation.RoomID
	if roomID == 0 {
		roomID, err = uc.AssignRoomForReservation(reservation)
		if err != nil {
			return dto.CheckInOutput{}, err
		}
	}

	checkInEmployeeID := input.EmployeeID

	stay, err := uc.stayService.RegisterStay(
		0,
		reservation.ClientID,
		roomID,
		&reservation.ID,
		input.CheckInTime,
		time.Time{},
		0,
		"",
		&checkInEmployeeID,
		nil,
		"",
	)
	if err != nil {
		return dto.CheckInOutput{}, err
	}

	return dto.CheckInOutput{
		StayID: stay.ID,
	}, nil
}

func (uc *DefaultEmployeeCheckInUseCase) AssignRoomForReservation(reservation *models.Reservation) (int, error) {
	availableRooms, err := uc.roomRepo.FindAvailableRooms(reservation.HotelID, reservation.StartDate, reservation.EndDate)
	if err != nil {
		return 0, err
	}

	if len(availableRooms) == 0 {
		return 0, errors.New("no available rooms")
	}

	selectedRoom := availableRooms[0]
	return selectedRoom.ID, nil
}
