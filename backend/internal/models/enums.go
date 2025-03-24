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

func (self Amenity) String() string {
	switch self {

	case WIFI:
		return "WIFI"
	case TV:
		return "TV"
	case AC:
		return "AC"
	case MiniFridge:
		return "MiniFridge"
	case CoffeeMachine:
		return "CoffeeMachine"
	case AirDryer:
		return "AirDryer"
	case Safe:
		return "Safe"
	case Jacuzzi:
		return "Jacuzzi"
	case Balcony:
		return "Balcony"
	case RoomService:
		return "RoomService"
	case KingSizeBed:
		return "KingSizeBed"
	case QueenSizeBed:
		return "QueenSizeBed"
	case SimpleBed:
		return "SimpleBed"
	case Office:
		return "Office"
	default:
		return "Invalid Amenity"
	}

}

func (_ Amenity) ParseAmenity(s string) (Amenity, error) {
	// trim and then lowercase the string
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "wifi":
		return WIFI, nil
	case "tv":
		return TV, nil
	case "ac":
		return AC, nil
	case "minifridge":
		return MiniFridge, nil
	case "coffeemachine":
		return CoffeeMachine, nil
	case "airdryer":
		return AirDryer, nil
	case "safe":
		return Safe, nil
	case "jacuzzi":
		return Jacuzzi, nil
	case "balcony":
		return Balcony, nil
	case "roomservice":
		return RoomService, nil
	case "kingsizebed":
		return KingSizeBed, nil
	case "queensizebed":
		return QueenSizeBed, nil
	case "simplebed":
		return SimpleBed, nil
	case "office":
		return Office, nil
	default:
		return 999, errors.New("Invalid amenity was passed.")
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

func (self RoomType) String() string {
	switch self {
	case Simple:
		return "Simple"
	case Double:
		return "Double"
	case Twin:
		return "Twin"
	case Queen:
		return "Queen"
	case King:
		return "King"
	case JuniorSuite:
		return "Junior Suite"
	case DeluxeSuite:
		return "Deluxe Suite"
	case FamilialSuite:
		return "Familial Suite"
	default:
		return "Invalid Room Type"
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

func (self ReservationStatus) String() string {
	switch self {
	case Confirmed:
		return "Confirmed"
	case Waiting:
		return "Waiting"
	case Cancelled:
		return "Cancelled"
	case Finished:
		return "Finished"
	default:
		return "Invalid Status"
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

func (self ViewType) String() string {
	switch self {
	case Sea:
		return "Sea"
	case Mountain:
		return "Mountain"
	case City:
		return "City"
	case Park:
		return "Park"
	case Courtyard:
		return "Courtyard"
	case Pool:
		return "Pool"
	default:
		return "Invalid View Type"
	}

}
