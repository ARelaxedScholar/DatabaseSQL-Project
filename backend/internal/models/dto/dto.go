package dto

import "time"

// New unified Client/Employee output
type AccountOutput struct {
	AccountID int       `json:"accountId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // e.g. "client" or "employee"
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

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
	HotelID         int       `josn:"hotelID"`
	RoomID          int       `json:"roomId"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	ReservationDate time.Time `json:"reservationDate"`
	TotalPrice      float64   `json:"totalPrice"`
	Status          int       `json:"status"` // we have in-house representation of this
}

type ReservationOutput struct {
	ReservationID int       `json:"reservationId"`
	ClientID      int       `json:"clientId"`
	RoomID        int       `json:"roomId"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	TotalPrice    float64   `json:"totalPrice"`
	Status        int       `json:"status"`
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

// Used by Client
type ClientProfileUpdateInput struct {
	ClientID  int    `json:"clientId"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Address   string `json:"address,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
}

// Used by Admin (Split object in case of future divergence)
type ClientAccountInput struct {
	SIN       string `json:"sin"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type ClientAccountUpdateInput struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Address   string `json:"address,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
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

// Used by Admin
type EmployeeAccountInput struct {
	SIN       string `json:"sin"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Position  string `json:"position"`
	HotelID   int    `json:"hotelId"`
	Password  string `json:"password"`
}

type EmployeeAccountUpdateInput struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Address   string `json:"address,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Position  string `json:"position,omitempty"`
	HotelID   int    `json:"hotelId,omitempty"`
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
	StartDate    *time.Time `json:"startDate,omitempty"`
	EndDate      *time.Time `json:"endDate,omitempty"`
	Capacity     *int       `json:"capacity,omitempty"`
	PriceMin     *float64   `json:"priceMin,omitempty"`
	PriceMax     *float64   `json:"priceMax,omitempty"`
	HotelChainID *int       `json:"hotelChainId,omitempty"`
	RoomType     *string    `json:"roomType,omitempty"`
}

type RoomSearchOutput struct {
	Rooms []RoomOutput `json:"rooms"`
}

type RoomOutput struct {
	RoomID       int      `json:"roomId"`
	HotelID      int      `json:"hotelId"`
	Capacity     int      `json:"capacity"`
	Number       string   `json:"number"`
	Floor        string   `json:"floor"`
	Price        float64  `json:"price"`
	Telephone    string   `json:"telephone"`
	ViewTypes    []string `json:"viewTypes"`
	RoomType     string   `json:"roomType"`
	IsExtensible bool     `json:"isExtensible"`
	Amenities    []string `json:"amenities"`
	Problems     []string `json:"problems"`
}

// Admin DTOs
type HotelInput struct {
	ID            int    `json:"id"`
	ChainID       int    `json:"chainId"`
	Rating        int    `json:"rating"`
	NumberOfRooms int    `json:"numberOfRooms"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	City          string `json:"city"`
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

// RoomInput is used for creating a new room.
type RoomInput struct {
	ID           int      `json:"id"` // Usually ignored/0 for create, set by DB
	HotelID      int      `json:"hotelId"`
	Capacity     int      `json:"capacity"`
	Number       string   `json:"number"`
	Floor        string   `json:"floor"`
	SurfaceArea  float64  `json:"surfaceArea"`
	Price        float64  `json:"price"`
	Telephone    string   `json:"telephone"`
	ViewTypes    []string `json:"viewTypes,omitempty"`
	RoomType     string   `json:"roomType"`
	IsExtensible bool     `json:"isExtensible"` // Defaults to false if omitted
	Amenities    []string `json:"amenities,omitempty"`
	Problems     []string `json:"problems,omitempty"`
}

// RoomUpdateInput is used for updating an existing room.
type RoomUpdateInput struct {
	ID           int       `json:"id"`
	HotelID      *int      `json:"hotelId,omitempty"`
	Capacity     *int      `json:"capacity,omitempty"`
	Number       *string   `json:"number,omitempty"`
	Floor        *string   `json:"floor,omitempty"`
	SurfaceArea  *float64  `json:"surfaceArea,omitempty"`
	Price        *float64  `json:"price,omitempty"`
	Telephone    *string   `json:"telephone,omitempty"`
	ViewTypes    *[]string `json:"viewTypes,omitempty"`
	RoomType     *string   `json:"roomType,omitempty"`
	IsExtensible *bool     `json:"isExtensible,omitempty"`
	Amenities    *[]string `json:"amenities,omitempty"`
	Problems     *[]string `json:"problems,omitempty"`
}

type RoomOutputAdmin struct {
	RoomID int `json:"roomId"`
}

// CheckoutInput represents the data required to perform a checkout.
type CheckoutInput struct {
	StayID        int
	FinalPrice    float64
	PaymentMethod string
}

// CheckoutOutput represents the result of the checkout operation.
type CheckoutOutput struct {
	StayID  int
	Message string
}
