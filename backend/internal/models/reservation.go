package models

import (
	"errors"
	"time"
)

type Reservation struct {
	id              int
	clientId        int
	roomId          int
	startDate       time.Time
	endDate         time.Time
	totalPrice      float64
	reservationDate time.Time
	status          ReservationStatus
}

func NewReservation(id, clientId int, roomId int, startDate, endDate, reservationDate time.Time, totalPrice float64, status ReservationStatus) (*Reservation, error) {
	var err error
	switch {
	case id < 0:
		err = errors.New("Reservation ID cannot be negative.")
	case clientId < 0:
		err = errors.New("Client's ID cannot be negative.")
	case roomId < 0:
		err = errors.New("Room's ID cannot be negative.")
	case !endDate.After(startDate):
		err = errors.New("End date must be after start date.")
	case totalPrice < 0:
		err = errors.New("Total price cannot be negative.")
	case !status.isValid():
		err = errors.New("Reservation status is invalid.")
	}
	if err != nil {
		return nil, err
	}
	return &Reservation{
		id:              id,
		clientId:        clientId,
		roomId:          roomId,
		startDate:       startDate,
		endDate:         endDate,
		totalPrice:      totalPrice,
		reservationDate: reservationDate,
		status:          status,
	}, nil
}
