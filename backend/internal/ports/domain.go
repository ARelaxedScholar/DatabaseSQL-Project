package ports

import (
	"time"

	"github.com/sql-project-backend/internal/models"
)

// This is where we declare the core functionalities
// we want for the app. So you can read this to know what the app does (is supposed to do) ;)
type EmployeeService interface {
	HireEmployee(id int, sin, firstName, lastName, address, phone, email, position string, hotelId int, hireDate time.Time) (*models.Employee, error)
	PromoteEmployeeToManager(employeeId int, department string, authorizationLevel int) (*models.Manager, error)
	FireEmployee(employeeId int) (*models.Employee, error)
	UpdateEmployee(employeeId int, firstName, lastName, address, phone, email, position string, hotelId int) (*models.Employee, error)
}

type ClientService interface {
	RegisterClient(id int, sin, firstName, lastName, address, phone, email string, joinDate time.Time) (*models.Client, error)
	UpdateClient(id int, firstName, lastName, address, phone, email string) (*models.Client, error)
	RemoveClient(id int) error
}

type HotelChainService interface {
	CreateHotelChain(id, numberOfHotel int, name, centralAddress, email, phone string) (*models.HotelChain, error)
	UpdateHotelChain(id, numberOfHotel int, name, centralAddress, email, phone string) (*models.HotelChain, error)
	DeleteHotelChain(id int) error
}

type HotelService interface {
	AddHotel(id, chainId, rating, numberOfRooms int, name, address, city, email, phone string) (*models.Hotel, error)
	UpdateHotel(id, chainId, rating, numberOfRooms int, name, address, city, email, phone string) (*models.Hotel, error)
	DeleteHotel(id int) error
}

type RoomService interface {
	AddRoom(id, hotelId, capacity int, number, floor string, surfaceArea, price float64, telephone string,
		viewTypes map[models.ViewType]struct{}, roomType models.RoomType, isExtensible bool,
		amenities map[models.Amenity]struct{}, problems []models.Problem) (*models.Room, error)
	UpdateRoom(id, hotelId, capacity int, number, floor string, surfaceArea, price float64, telephone string,
		viewTypes map[models.ViewType]struct{}, roomType models.RoomType, isExtensible bool,
		amenities map[models.Amenity]struct{}, problems []models.Problem) (*models.Room, error)
	DeleteRoom(id int) error
	AssignRoomForReservation(reservation *models.Reservation) (int, error)
	FindAvailableRooms(hotelID int, startDate time.Time, endDate time.Time) ([]*models.Room, error)
}

type ReservationService interface {
	CreateReservation(id, clientId, hotelID, roomId int,
		startDate, endDate, reservationDate time.Time, totalPrice float64, status models.ReservationStatus) (*models.Reservation, error)
	UpdateReservation(id, clientId, hotelID, roomId int,
		startDate, endDate, reservationDate time.Time, totalPrice float64, status models.ReservationStatus) (*models.Reservation, error)
	CancelReservation(id int) error
	GetReservationsByClient(clientID int) ([]*models.Reservation, error)
}

type StayService interface {
	RegisterStay(id, clientId, roomId int, reservationId *int,
		arrivalDate time.Time, departureDate *time.Time,
		checkInEmployeeId int, checkOutEmployeeId *int, comments string) (*models.Stay, error)
	UpdateStay(id, clientId int, roomId int, reservationId *int,
		arrivalDate time.Time, departureDate *time.Time,
		checkInEmployeeId int, checkOutEmployeeId *int, comments string) (*models.Stay, error)
	EndStay(id, employeeID int) error
}

type PaymentService interface {
	ProcessPayment(stayId int, amount float64, paymentMethod string) error
}

type QueryService interface {
	GetAvailableRoomsByZone() (map[string]int, error)
	GetHotelRoomCapacity(hotelId int) (int, error)
}
