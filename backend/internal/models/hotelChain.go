package models

import "errors"

type HotelChain struct {
	ID, NumberOfHotel                      int
	Name, CentralAddress, Email, Telephone string
}

func NewHotelChain(id, numberOfHotel int, name, centralAddress, email, telephone string) (*HotelChain, error) {
	var err error
	// Validate Fields
	switch {
	case id < 0:
		err = errors.New("Hotel chain's ID cannot be negative.")
	case numberOfHotel < 0:
		err = errors.New("Number of hotels cannot be negative.")
	case name == "":
		err = errors.New("Hotel chain's name cannot be empty.")

	case centralAddress == "":
		err = errors.New("Central address cannot be empty.")

	case email == "":
		err = errors.New("Hotel chain's email cannot be empty.")

	case telephone == "":
		err = errors.New("Hotel chain's phone Number cannot be empty.")

	}

	// if a case was hit
	if err != nil {
		return nil, err
	}

	// else just return
	return &HotelChain{
		ID:             id,
		NumberOfHotel:  numberOfHotel,
		Name:           name,
		CentralAddress: centralAddress,
		Email:          email,
		Telephone:      telephone,
	}, nil
}
