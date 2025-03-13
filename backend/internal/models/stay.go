package models

import (
	"errors"
	"time"
)

type Stay struct {
	id                 int
	clientId           int
	roomId             int
	reservationId      *int // nil for walk-ins
	arrivalDate        time.Time
	departureDate      time.Time
	finalPrice         float64
	paymentMethod      string
	checkInEmployeeId  *int // nil if not applicable
	checkOutEmployeeId *int // nil if not applicable
	comments           string
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
		id:                 id,
		clientId:           clientId,
		roomId:             roomId,
		reservationId:      reservationId,
		arrivalDate:        arrivalDate,
		departureDate:      departureDate,
		finalPrice:         finalPrice,
		paymentMethod:      paymentMethod,
		checkInEmployeeId:  checkInEmployeeId,
		checkOutEmployeeId: checkOutEmployeeId,
		comments:           comments,
	}, nil
}
