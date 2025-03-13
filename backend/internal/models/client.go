package models

import (
	"errors"
	"time"
)

type Client struct {
	id        int
	sin       string
	firstName string
	lastName  string
	address   string
	phone     string
	email     string
	joinDate  time.Time
}

func NewClient(id int, sin, firstName, lastName, address, phone, email string, joinDate time.Time) (*Client, error) {
	var err error
	switch {
	case id < 0:
		err = errors.New("Client ID cannot be negative.")
	case sin == "" || len(sin) != 9:
		err = errors.New("Client's SIN must be 9 characters.")
	case firstName == "":
		err = errors.New("Client's first name cannot be empty.")
	case lastName == "":
		err = errors.New("Client's last name cannot be empty.")
	case joinDate.IsZero():
		err = errors.New("Client's join date must be provided.")
	case email == "":
		err = errors.New("Client's email cannot be empty.")
	}
	if err != nil {
		return nil, err
	}
	return &Client{
		id:        id,
		sin:       sin,
		firstName: firstName,
		lastName:  lastName,
		address:   address,
		phone:     phone,
		email:     email,
		joinDate:  joinDate,
	}, nil
}
