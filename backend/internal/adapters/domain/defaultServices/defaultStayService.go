package defaultServices

import (
	"errors"
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultStayService struct {
	stayRepo ports.StayRepository
}

func NewStayService(repo ports.StayRepository) ports.StayService {
	return &DefaultStayService{
		stayRepo: repo,
	}
}

func (s *DefaultStayService) RegisterStay(id, clientId, roomId int, reservationId *int, checkInTime time.Time, checkOutTime *time.Time, checkInEmployeeId int, checkOutEmployeeId *int, comments string) (*models.Stay, error) {
	stay, err := models.NewStay(id, clientId, roomId, checkInEmployeeId, checkOutEmployeeId, reservationId, checkInTime, checkOutTime, comments)
	if err != nil {
		return nil, err
	}
	dbStay, err := s.stayRepo.Save(stay)
	if err != nil {
		return nil, err
	}
	return dbStay, nil
}

func (s *DefaultStayService) UpdateStay(id, clientId, roomId int, reservationId *int, checkInTime time.Time, checkOutTime *time.Time, checkInEmployeeId int, checkOutEmployeeId *int, comments string) (*models.Stay, error) {
	stay, err := s.stayRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if stay == nil {
		return nil, errors.New("Stay not found.")
	}
	stay, err = models.NewStay(id, clientId, roomId, checkInEmployeeId, checkOutEmployeeId, reservationId, checkInTime, checkOutTime, comments)
	if err != nil {
		return nil, err
	}
	if err = s.stayRepo.Update(stay); err != nil {
		return nil, err
	}
	return stay, nil
}

func (s *DefaultStayService) EndStay(id, employeeID int) error {
	stay, err := s.stayRepo.FindByID(id)
	if err != nil {
		return err
	}
	if stay == nil {
		return errors.New("Stay not found.")
	}
	// EndStay
	return s.stayRepo.EndStay(id, employeeID)
}
