package models

import "errors"

type Room struct {
	ID, HotelID, Capacity, Floor int
	Price                        float64
	Telephone                    string
	ViewTypes                    map[ViewType]struct{}
	RoomType                     RoomType
	IsExtensible                 bool
	Amenities                    map[Amenity]struct{} // Hashset of Go (I know stupid)
	Problems                     []Problem
}

func NewRoom(id, hotelId, capacity, floor int,
	price float64,
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
	case price < 0:
		err = errors.New("Room's price cannot be negative.")
	case telephone == "":
		err = errors.New("Room's phone Number cannot be empty")
	case !roomType.isValid():
		err = errors.New("Invalid variant of room type was passed to constructor")
	}
	// If we haven't already found an error
	if err == nil {
		// validate view types
		for k := range viewTypes {
			if !k.isValid() {
				err = errors.New("The set of view types contains an invalid variant")
				break
			}
		}
	}
	if err == nil {
		// validate amenities
		for k := range amenities {
			if !k.isValid() {
				err = errors.New("The set of amenities contains an invalid variant")
				break
			}
		}
	}
	if err == nil {
		// validate amenities
		for _, problem := range problems {
			if e := problem.validate(); e != nil {
				err = errors.Join(errors.New("Problem in the list of problem is invalid:\n"), e)
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
		Floor:        floor,
		Price:        price,
		Telephone:    telephone,
		ViewTypes:    viewTypes,
		RoomType:     roomType,
		IsExtensible: isExtensible,
		Amenities:    amenities,
	}, nil

}
