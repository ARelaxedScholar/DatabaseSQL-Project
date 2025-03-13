package models

import (
	"errors"
)

// Amenities Enum
type Amenity int

const (
	WIFI Amenity = iota + 1
	TV
	AC
	MiniFridge
	CoffeeMachine
	AirDryer
	Safe
	Jacuzzi
	Balcony
	RoomService
	KingSizeBed
	QueenSizeBed
	SimpleBed
	Office
)

func (self Amenity) isValid() bool {
	switch self {
	case WIFI,
		TV,
		AC,
		MiniFridge,
		CoffeeMachine,
		AirDryer,
		Safe,
		Jacuzzi,
		Balcony,
		RoomService,
		KingSizeBed,
		QueenSizeBed,
		SimpleBed,
		Office:
		return true
	default:
		return false
	}
}

// View Enum
type ViewType int

const (
	Sea ViewType = iota + 1
	Mountain
	City
	Park
	Courtyard
	Pool
)

func (self ViewType) isValid() bool {
	switch self {
	case Sea, Mountain, City, Park, Courtyard, Pool:
		return true
	default:
		return false
	}
}

// RoomType
type RoomType int

const (
	Simple RoomType = iota + 1
	Double
	Twin
	Queen
	King
	JuniorSuite
	DeluxeSuite
	FamilialSuite
)

func (self RoomType) isValid() bool {
	switch self {
	case Simple,
		Double,
		Twin,
		Queen,
		King,
		JuniorSuite,
		DeluxeSuite,
		FamilialSuite:
		return true
	default:
		return false
	}
}

type HotelChain struct {
	id, numberOfHotel                      int
	name, centralAddress, email, telephone string
}

type Hotel struct {
	id, chainId, rating, numberOfRooms int
	name, address, email, telephone    string
}

type Room struct {
	id, hotelId, capacity, floor int
	price                        float64
	telephone                    string
	viewTypes                    map[ViewType]struct{}
	roomType                     RoomType
	isExtensible                 bool
	amenities                    map[Amenity]struct{} // Hashset of Go (I know stupid)
	problems                     []Problem
}

type ProblemSeverity int
type Problem struct {
	id          int
	severity    ProblemSeverity
	description string
}

func (self Problem) validate() error {
	var err error
	switch {
	case self.id < 0:
		err = errors.New("Problem's id cannot be negative.")
	case !self.severity.isValid():
		err = errors.New("Problem contains invalid severity variant.")
	case self.description == "":
		err = errors.New("Problem's description cannot be empty.")
	}
	return err // note that this will be nil if none of the cases were hit

}

const (
	Minor ProblemSeverity = iota + 1
	Moderate
	Major
	Critical
)

func (self ProblemSeverity) isValid() bool {
	switch self {
	case Minor, Moderate, Major, Critical:
		return true
	default:
		return false
	}
}

// todo: Define the constructors for this
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
		id:             id,
		numberOfHotel:  numberOfHotel,
		name:           name,
		centralAddress: centralAddress,
		email:          email,
		telephone:      telephone,
	}, nil
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
		id:            id,
		chainId:       chainId,
		rating:        rating,
		numberOfRooms: numberOfRooms,
		name:          name,
		address:       address,
		email:         email,
		telephone:     telephone,
	}, nil

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
		for k, _ := range viewTypes {
			if !k.isValid() {
				err = errors.New("The set of view types contains an invalid variant")
				break
			}
		}
	}
	if err == nil {
		// validate amenities
		for k, _ := range amenities {
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
		id:           id,
		hotelId:      hotelId,
		capacity:     capacity,
		floor:        floor,
		price:        price,
		telephone:    telephone,
		viewTypes:    viewTypes,
		roomType:     roomType,
		isExtensible: isExtensible,
		amenities:    amenities,
	}, nil

}
