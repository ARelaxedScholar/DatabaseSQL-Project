package models

import (
	"errors"
)

type Hotel struct {
	ID, ChainID, Rating, NumberOfRooms int
	Name, Address, Email, Telephone    string
}

func NewHotel(id, chainId, rating, numberOfRooms int,
	name, address, email, telephone string) (*Hotel, error) {
	var err error
	// Validate Fields
	switch {
	case id < 0:
		err = errors.New("Hotel's ID cannot be negative.")
	case chainId < 0:
		err = errors.New("Hotel chain's ID cannot be negative.")
	case rating < 1:
		err = errors.New("Hotel's rating must be at least 1.")
	case rating > 5:
		err = errors.New("Hotel's rating can be at most 5.")
	case numberOfRooms < 1:
		err = errors.New("Hotel must have at least one room.")
	case name == "":
		err = errors.New("Hotel's name cannot be empty")
	case address == "":
		err = errors.New("Hotel's address cannot be empty")
	case email == "":
		err = errors.New("Hotel's email cannot be empty")
	case telephone == "":
		err = errors.New("Hotel's phone Number cannot be empty")
	}

	// if a case was hit
	if err != nil {
		return nil, err
	}

	// else just return
	return &Hotel{
		ID:            id,
		ChainID:       chainId,
		Rating:        rating,
		NumberOfRooms: numberOfRooms,
		Name:          name,
		Address:       address,
		Email:         email,
		Telephone:     telephone,
	}, nil

}
