package models

import (
	"errors"
	"time"
)

type Stay struct {
	ID                 int
	ClientID           int
	RoomID             int
	ReservationID      *int // nil for walk-ins
	CheckInTime        time.Time
	CheckOutTime       *time.Time
	FinalPrice         *float64
	PaymentMethod      *string
	CheckInEmployeeId  int  // nil if not applicable
	CheckOutEmployeeId *int // nil if not applicable
	Comments           string
}

func NewStay(id int, clientId int, roomId, checkInEmployeeId int, checkOutEmployeeId, reservationId *int, checkInTime time.Time, checkOutTime *time.Time, comments string) (*Stay, error) {
	var err error
	switch {
	case id < 0:
		err = errors.New("Stay's ID cannot be negative.")
	case checkInEmployeeId < 0:
		err = errors.New("Employee's ID (Check in) cannot be negative.")
	case clientId < 0:
		err = errors.New("Client's ID must be 9 characters.")
	case roomId < 0:
		err = errors.New("Room ID cannot be negative.")
	}
	if err != nil {
		return nil, err
	}
	if checkOutTime != nil {
		if !(*checkOutTime).Before(checkInTime) {
			err = errors.New("CheckOut time must be after or equal to checkin time")
		}

	}
	return &Stay{
		ID:                 id,
		ClientID:           clientId,
		RoomID:             roomId,
		ReservationID:      reservationId,
		CheckInTime:        checkInTime,
		CheckOutTime:       checkOutTime,
		FinalPrice:         nil,
		PaymentMethod:      nil,
		CheckInEmployeeId:  checkInEmployeeId,
		CheckOutEmployeeId: checkOutEmployeeId,
		Comments:           comments,
	}, nil
}
