package defaultEmployeeUseCases

import (
	"github.com/sql-project-backend/internal/models/dto"
	"github.com/sql-project-backend/internal/ports"
)

type DefaultEmployeeCreateNewStayUseCase struct {
	stayService ports.StayService
}

func NewEmployeeCreateNewStayUseCase(stayService ports.StayService) ports.EmployeeCreateNewStayUseCase {
	return &DefaultEmployeeCreateNewStayUseCase{
		stayService: stayService,
	}
}

func (uc *DefaultEmployeeCreateNewStayUseCase) CreateNewStay(input dto.NewStayInput) (dto.NewStayOutput, error) {
	stay, err := uc.stayService.RegisterStay(
		0, // new stay, ID will be generated
		input.ClientID,
		input.RoomID,
		input.ReservationID, // no reservation ID (e.g., for walk-in) would be nil
		input.CheckInTime,
		nil,
		input.CheckInEmployeeID,
		nil, // no check-out employee at creation
		input.Comments,
	)
	if err != nil {
		return dto.NewStayOutput{}, err
	}
	return dto.NewStayOutput{
		StayID: stay.ID,
	}, nil
}
