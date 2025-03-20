package models

import (
	"errors"
	"time"
)

type Reservation struct {
	ID              int
	HotelID         int
	ClientID        int
	RoomID          int
	StartDate       time.Time
	EndDate         time.Time
	TotalPrice      float64
	ReservationDate time.Time
	Status          ReservationStatus
}

func NewReservation(id, clientId, hotelID, roomId int, startDate, endDate, reservationDate time.Time, totalPrice float64, status ReservationStatus) (*Reservation, error) {
	var err error
	switch {
	case id < 0:
		err = errors.New("Reservation ID cannot be negative.")
	case clientId < 0:
		err = errors.New("Client's ID cannot be negative.")
	case hotelID < 0: // <- I might turn this into an ENUM later
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
		ID:              id,
		ClientID:        clientId,
		HotelID:         hotelID,
		RoomID:          roomId,
		StartDate:       startDate,
		EndDate:         endDate,
		TotalPrice:      totalPrice,
		ReservationDate: reservationDate,
		Status:          status,
	}, nil
}
