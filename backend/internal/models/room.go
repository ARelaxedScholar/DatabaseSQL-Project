package models

import (
	"errors"
	"fmt"
	"strings"
)

type Room struct {
	ID           int
	HotelID      int
	Capacity     int
	Number       string
	Floor        string
	SurfaceArea  float64
	Price        float64
	Telephone    string
	ViewTypes    map[ViewType]struct{}
	RoomType     RoomType
	IsExtensible bool
	Amenities    map[Amenity]struct{}
	Problems     []Problem
}

// Updated constructor signature to include surfaceArea
func NewRoom(id, hotelId, capacity int, number, floor string, surfaceArea, price float64, // Added surfaceArea parameter
	telephone string,
	viewTypes map[ViewType]struct{},
	roomType RoomType,
	isExtensible bool,
	amenities map[Amenity]struct{},
	problems []Problem) (*Room, error) {

	var err error
	// Validate Fields
	switch {
	case id < 0:
		err = errors.New("Room's ID cannot be negative.")
	case hotelId < 0:
		err = errors.New("Hotel's ID cannot be negative.")
	case capacity < 1:
		err = errors.New("Room's capacity must be at least 1.")
	case strings.TrimSpace(number) == "":
		err = errors.New("Room's number cannot be empty.")
	case strings.TrimSpace(floor) == "":
		err = errors.New("Room's floor cannot be empty.")
	case surfaceArea <= 0: // Added validation for surfaceArea (must be positive)
		err = errors.New("Room's surface area must be positive.")
	case price < 0:
		err = errors.New("Room's price cannot be negative.")
	case telephone == "":
		err = errors.New("Room's phone number cannot be empty.")
	case !roomType.isValid():
		err = errors.New("Invalid variant of room type was passed to constructor.")
	}

	if err == nil {
		for k := range viewTypes {
			if !k.isValid() {
				err = errors.New("The set of view types contains an invalid variant.")
				break
			}
		}
	}
	if err == nil {
		for k := range amenities {
			if !k.isValid() {
				err = errors.New("The set of amenities contains an invalid variant.")
				break
			}
		}
	}
	if err == nil {
		for i := range problems {
			if e := problems[i].Validate(); e != nil {
				err = fmt.Errorf("Problem at index %d is invalid: %w", i, e)
				break
			}
		}
	}

	if err != nil {
		return nil, err
	}

	return &Room{
		ID:           id,
		HotelID:      hotelId,
		Capacity:     capacity,
		Number:       strings.TrimSpace(number),
		Floor:        strings.TrimSpace(floor),
		SurfaceArea:  surfaceArea, // Assign surfaceArea
		Price:        price,
		Telephone:    telephone,
		ViewTypes:    viewTypes,
		RoomType:     roomType,
		IsExtensible: isExtensible,
		Amenities:    amenities,
		Problems:     problems,
	}, nil
}

func (p *Problem) Validate() error {
	var err error
	switch {
	case p.ID < 0 && p.ID != 0:
		err = errors.New("Problem's id cannot be negative.")
	case !p.Severity.isValid():
		err = errors.New("Problem contains invalid severity variant.")
	case p.Description == "":
		err = errors.New("Problem's description cannot be empty.")
	case p.SignaledWhen.IsZero():
		err = errors.New("Problem's signaled date cannot be zero.")
	case p.ResolutionDate.IsZero() && p.IsResolved:
		err = errors.New("Problem is resolved, but no resolution time was passed.")
	case !p.ResolutionDate.IsZero() && !p.IsResolved:
		err = errors.New("Problem has a resolution date but is marked as not resolved.")
	case !p.ResolutionDate.IsZero() && p.ResolutionDate.Before(p.SignaledWhen):
		err = errors.New("Problem's resolution date cannot be before its signaled date.")
	}
	return err
}
