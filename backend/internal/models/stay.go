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
	ArrivalDate        time.Time
	DepartureDate      time.Time
	FinalPrice         float64
	PaymentMethod      string
	CheckInEmployeeId  *int // nil if not applicable
	CheckOutEmployeeId *int // nil if not applicable
	Comments           string
}

func NewStay(id int, clientId int, roomId int, checkInEmployeeId, checkOutEmployeeId, reservationId *int, arrivalDate, departureDate time.Time, finalPrice float64, paymentMethod string, comments string) (*Stay, error) {
	var err error
	switch {
	case id < 0:
		err = errors.New("Stay's ID cannot be negative.")
	case clientId < 0:
		err = errors.New("Client's ID must be 9 characters.")
	case roomId < 0:
		err = errors.New("Room ID cannot be negative.")
	case departureDate.Before(arrivalDate):
		err = errors.New("Departure date must be after or equal to arrival date.")
	case finalPrice < 0:
		err = errors.New("Final price cannot be negative.")
	case paymentMethod == "":
		err = errors.New("Payment method cannot be empty.")
	}
	if err != nil {
		return nil, err
	}
	return &Stay{
		ID:                 id,
		ClientID:           clientId,
		RoomID:             roomId,
		ReservationID:      reservationId,
		ArrivalDate:        arrivalDate,
		DepartureDate:      departureDate,
		FinalPrice:         finalPrice,
		PaymentMethod:      paymentMethod,
		CheckInEmployeeId:  checkInEmployeeId,
		CheckOutEmployeeId: checkOutEmployeeId,
		Comments:           comments,
	}, nil
}
