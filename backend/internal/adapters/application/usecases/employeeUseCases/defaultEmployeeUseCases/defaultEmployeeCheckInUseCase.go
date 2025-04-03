package defaultEmployeeUseCases

import (
	"errors"
	"log"

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
	var reservation *models.Reservation
	if input.ReservationID != nil {
		res, err := uc.reservationRepo.FindByID(*input.ReservationID)
		reservation = res
		if err != nil {
			return dto.CheckInOutput{}, err
		}
	}

	if reservation == nil && input.ReservationID != nil {
		return dto.CheckInOutput{}, errors.New("reservation not found")
	} else if reservation == nil {
		log.Printf("At %v: A stay was created without prior reservation", input.CheckInTime)
	}

	if reservation != nil {
		// Do further validations about time
		if input.CheckInTime.Before(reservation.StartDate) {
			return dto.CheckInOutput{}, errors.New("Attempt to check in on a reservation before reservation started.")
		} else if input.CheckInTime.After(reservation.EndDate) {
			return dto.CheckInOutput{}, errors.New("Attempt to check in on a reservation after reservation ended.")
		}
	}

	var err error
	roomID := reservation.RoomID
	if roomID == 0 { // if default room ID is passed we look for a free room
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
		input.ReservationID,
		input.CheckInTime,
		nil,
		checkInEmployeeID,
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
