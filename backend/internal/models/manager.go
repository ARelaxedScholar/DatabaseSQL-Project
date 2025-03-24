package models

import (
	"errors"
	"time"
)

type Manager struct {
	Employee
	Department         string
	AuthorizationLevel int
}

func NewManager(sin, firstName, lastName, address, phone, email, position string, id, hotelId int, hireDate time.Time, department string, authorizationLevel int) (*Manager, error) {
	var err error
	switch {
	case department == "":
		err = errors.New("Department cannot be empty")
	case authorizationLevel < 1 || authorizationLevel > 5:
		err = errors.New("Authorization level must be between 1 and 5")
	}
	if err != nil {
		return nil, err
	}
	emp, err := NewEmployee(sin, firstName, lastName, address, phone, email, position, id, hotelId, hireDate)
	if err != nil {
		return nil, err
	}
	return &Manager{
		Employee:           *emp,
		Department:         department,
		AuthorizationLevel: authorizationLevel,
	}, nil
}
