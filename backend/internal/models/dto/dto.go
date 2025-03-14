package dto

import "time"

// Client DTOs
type ClientRegistrationInput struct {
	SIN       string    `json:"sin"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"` // this is the username 
	JoinDate  time.Time `json:"joinDate"`
}

type ClientRegistrationOutput struct {
	ClientID int `json:"clientId"`
}

type ClientLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ClientLoginOutput struct {
	ClientID int    `json:"clientId"`
	Token    string `json:"token"`
}

type ReservationInput struct {
	ClientID        int       `json:"clientId"`
	RoomID          int       `json:"roomId"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	ReservationDate time.Time `json:"reservationDate"`
	TotalPrice      float64   `json:"totalPrice"`
	Status          string    `json:"status"`
}

type ReservationOutput struct {
	ReservationID int       `json:"reservationId"`
	ClientID      int       `json:"clientId"`
	RoomID        int       `json:"roomId"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	TotalPrice    float64   `json:"totalPrice"`
	Status        string    `json:"status"`
}

type ClientProfileOutput struct {
	ClientID  int       `json:"clientId"`
	SIN       string    `json:"sin"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	JoinDate  time.Time `json:"joinDate"`
}

type ClientProfileUpdateInput struct {
	ClientID  int    `json:"clientId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

// Employee DTOs
type EmployeeRegistrationInput struct {
	EmployeeID int       `json:"employeeId"`
	SIN        string    `json:"sin"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Position   string    `json:"position"`
	HotelID    int       `json:"hotelId"`
	HireDate   time.Time `json:"hireDate"`
}

type EmployeeRegistrationOutput struct {
	EmployeeID int `json:"employeeId"`
}

type EmployeeLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmployeeLoginOutput struct {
	EmployeeID int    `json:"employeeId"`
	Token      string `json:"token"`
}

type CheckInInput struct {
	ReservationID int       `json:"reservationId"`
	EmployeeID    int       `json:"employeeId"`
	CheckInTime   time.Time `json:"checkInTime"`
}

type CheckInOutput struct {
	StayID int `json:"stayId"`
}

type NewStayInput struct {
	ClientID          int       `json:"clientId"`
	RoomID            int       `json:"roomId"`
	CheckInEmployeeID int       `json:"checkInEmployeeId"`
	ArrivalDate       time.Time `json:"arrivalDate"`
	DepartureDate     time.Time `json:"departureDate"`
	FinalPrice        float64   `json:"finalPrice"`
	PaymentMethod     string    `json:"paymentMethod"`
	Comments          string    `json:"comments"`
}

type NewStayOutput struct {
	StayID int `json:"stayId"`
}

// Anonymous Use Case DTOs
type RoomSearchInput struct {
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	Capacity     int       `json:"capacity"`
	PriceMin     float64   `json:"priceMin"`
	PriceMax     float64   `json:"priceMax"`
	HotelChainID int       `json:"hotelChainId"`
	Category     string    `json:"category"`
}

type RoomSearchOutput struct {
	Rooms []RoomOutput `json:"rooms"`
}

type RoomOutput struct {
	RoomID       int      `json:"roomId"`
	HotelID      int      `json:"hotelId"`
	Capacity     int      `json:"capacity"`
	Floor        int      `json:"floor"`
	Price        float64  `json:"price"`
	Telephone    string   `json:"telephone"`
	ViewTypes    []string `json:"viewTypes"`
	RoomType     string   `json:"roomType"`
	IsExtensible bool     `json:"isExtensible"`
	Amenities    []string `json:"amenities"`
}

// Admin DTOs
type HotelInput struct {
	ID            int    `json:"id"`
	ChainID       int    `json:"chainId"`
	Rating        int    `json:"rating"`
	NumberOfRooms int    `json:"numberOfRooms"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
}

type HotelOutput struct {
	HotelID int `json:"hotelId"`
}

type HotelChainInput struct {
	ID             int    `json:"id"`
	NumberOfHotels int    `json:"numberOfHotels"`
	Name           string `json:"name"`
	CentralAddress string `json:"centralAddress"`
	Email          string `json:"email"`
	Telephone      string `json:"telephone"`
}

type HotelChainOutput struct {
	ChainID int `json:"chainId"`
}

type RoomInput struct {
	ID           int      `json:"id"`
	HotelID      int      `json:"hotelId"`
	Capacity     int      `json:"capacity"`
	Floor        int      `json:"floor"`
	Price        float64  `json:"price"`
	Telephone    string   `json:"telephone"`
	ViewTypes    []string `json:"viewTypes"`
	RoomType     string   `json:"roomType"`
	IsExtensible bool     `json:"isExtensible"`
	Amenities    []string `json:"amenities"`
}

type RoomUpdateInput struct {
	ID           int      `json:"id"`
	HotelID      int      `json:"hotelId"`
	Capacity     int      `json:"capacity"`
	Floor        int      `json:"floor"`
	Price        float64  `json:"price"`
	Telephone    string   `json:"telephone"`
	ViewTypes    []string `json:"viewTypes"`
	RoomType     string   `json:"roomType"`
	IsExtensible bool     `json:"isExtensible"`
	Amenities    []string `json:"amenities"`
}

type RoomOutputAdmin struct {
	RoomID int `json:"roomId"`
}
