package defaultServices

import (
	"errors"
	"fmt"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultReservationService struct {
	reservationRepo ports.ReservationRepository
}

func NewReservationService(repo ports.ReservationRepository) ports.ReservationService {
	return &DefaultReservationService{
		reservationRepo: repo,
	}
}

func (s *DefaultReservationService) CreateReservation(id, clientId, hotelID, roomId int, startDate, endDate, reservationDate time.Time, totalPrice float64, status models.ReservationStatus) (*models.Reservation, error) {
	reservation, err := models.NewReservation(id, clientId, hotelID, roomId, startDate, endDate, reservationDate, totalPrice, status)
	if err != nil {
		return nil, err
	}
	dbReservation, err := s.reservationRepo.Save(reservation)
	if err != nil {
		return nil, err
	}
	return dbReservation, nil
}

func (s *DefaultReservationService) UpdateReservation(id, clientId, hotelId, roomId int, startDate, endDate, reservationDate time.Time, totalPrice float64, status models.ReservationStatus) (*models.Reservation, error) {
	existing, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("Reservation not found.")
	}
	reservation, err := models.NewReservation(id, clientId, hotelId, roomId, startDate, endDate, reservationDate, totalPrice, status)
	if err != nil {
		return nil, err
	}
	if err = s.reservationRepo.Update(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func (s *DefaultReservationService) CancelReservation(id int) error {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return err
	}
	if reservation == nil {
		return errors.New("Reservation not found.")
	}
	// Update the reservation's status to cancelled.
	reservation.Status = models.Cancelled
	if err = s.reservationRepo.Update(reservation); err != nil {
		return err
	}
	return nil
}

func (s *DefaultReservationService) CancelReservationForUser(id, clientID int) error {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return err
	}
	if reservation == nil {
		return errors.New("Reservation not found.")
	}
	if reservation.ClientID != clientID {
		return errors.New(fmt.Sprintf("This reservation (id: %d) does not belong to user %d.", id, clientID))
	}
	// Update the reservation's status to cancelled.
	reservation.Status = models.Cancelled
	if err = s.reservationRepo.Update(reservation); err != nil {
		return err
	}
	return nil
}

func (s *DefaultReservationService) GetReservationsByClient(clientID int) ([]*models.Reservation, error) {
	return s.reservationRepo.GetByClient(clientID)
}
