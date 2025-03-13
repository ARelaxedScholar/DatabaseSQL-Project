package models

// Amenities Enum
type Amenities int

const (
	WIFI Amenities = iota + 1
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
	viewType                     ViewType
	roomType                     RoomType
	isExtensible                 bool
}

// todo: Define the constructors for this
