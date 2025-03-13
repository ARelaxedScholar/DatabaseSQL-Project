package models

import (
	"errors"
	"time"
)

type Employee struct {
	id        int
	sin       string
	firstName string
	lastName  string
	address   string
	phone     string
	email     string
	hotelId   int
	position  string
	hireDate  time.Time
}

func NewEmployee(sin, firstName, lastName, address, phone, email, position string, id, hotelId int, hireDate time.Time) (*Employee, error) {
	var err error
	switch {
	case id < 0:
		err = errors.New("Employee's ID cannot be negative")
	case sin == "" || len(sin) != 9:
		err = errors.New("Employee's SIN must be 9 characters")
	case firstName == "":
		err = errors.New("Employee's first name cannot be empty.")
	case lastName == "":
		err = errors.New("Employee's last name cannot be empty.")
	case address == "":
		err = errors.New("Employee's address cannot be empty.")
	case phone == "":
		err = errors.New("Employee's phone cannot be empty.")
	case email == "":
		err = errors.New("Employee's email cannot be empty.")
	case hotelId < 0:
		err = errors.New("Hotel id's cannot be negative.")
	case position == "":
		err = errors.New("Employee's position cannot be empty.")
	}
	if err != nil {
		return nil, err
	}
	return &Employee{
		id:        id,
		sin:       sin,
		firstName: firstName,
		lastName:  lastName,
		address:   address,
		phone:     phone,
		email:     email,
		hotelId:   hotelId,
		position:  position,
		hireDate:  hireDate,
	}, nil
}
