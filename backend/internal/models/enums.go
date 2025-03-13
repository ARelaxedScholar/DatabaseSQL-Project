package models

// ### AMENITIES SECTION
//
//	Amenities Enum
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

// Amenities Methods
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

// ### PROBLEM SECTION
// Problem Severity
type ProblemSeverity int

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

// ### ROOM TYPE SECTION
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

// ### RESERVATION STATUS SECTION
type ReservationStatus int

const (
	Confirmed ReservationStatus = iota + 1
	Waiting
	Cancelled
	Finished
)

func (self ReservationStatus) isValid() bool {
	switch self {
	case Confirmed, Waiting, Cancelled, Finished:
		return true
	default:
		return false
	}
}

// ### VIEW TYPE SECTION
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

// View Methods
func (self ViewType) isValid() bool {
	switch self {
	case Sea, Mountain, City, Park, Courtyard, Pool:
		return true
	default:
		return false
	}
}
