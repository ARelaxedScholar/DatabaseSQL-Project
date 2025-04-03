package ports

import (
	"time"

	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/models/dto"
)

//import "github.com/sql-project-backend/internal/models/dto"

// TokenService (For the Login)
type TokenService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (int, error)
}

// Outline all possible usecases
// at this point I foresee some of em might be dropped/just directly implemented without interface
// this is for architecting (For now the Admin use cases can be more or less ignored until the rest is done)

// ## Client USE CASES
type ClientRegistrationUseCase interface {
	RegisterClient(input dto.ClientRegistrationInput) (dto.ClientRegistrationOutput, error)
}

type ClientLoginUseCase interface {
	Login(input dto.ClientLoginInput) (dto.ClientLoginOutput, error)
}

type ClientMakeReservationUseCase interface {
	MakeReservation(input dto.ReservationInput) (dto.ReservationOutput, error)
}

type ClientReservationsManagementUseCase interface {
	ViewReservations(clientID int) ([]dto.ReservationOutput, error)
	CancelReservation(reservationID int, userID int) error
}

type ClientProfileManagementUseCase interface {
	GetProfile(clientID int) (dto.ClientProfileOutput, error)
	UpdateProfile(input dto.ClientProfileUpdateInput) (dto.ClientProfileOutput, error)
}

// ## Employee USE CASES
type EmployeeLoginUseCase interface {
	Login(input dto.EmployeeLoginInput) (dto.EmployeeLoginOutput, error)
}

// This when we already have a reservation
type EmployeeCheckInUseCase interface {
	CheckIn(input dto.CheckInInput) (dto.CheckInOutput, error)
}

type EmployeeCheckoutUseCase interface {
	Checkout(input dto.CheckoutInput) (dto.CheckoutOutput, error)
}

// This is for when stays are created outside of check-in context
type EmployeeCreateNewStayUseCase interface {
	CreateNewStay(input dto.NewStayInput) (dto.NewStayOutput, error)
}

// ## Anonymous
type SearchRoomsUseCase interface {
	SearchRooms(input dto.RoomSearchInput) (dto.RoomSearchOutput, error)
	GetNumberOfRoomsForHotel(hotelID int) (int, error)
	GetNumberOfRoomsPerZone()  (map[string]int, error)  // Zone == City
}

// ## Admin USE CASES (Right now no requirement for that so kind of an after thought)
type AdminHotelManagementUseCase interface {
	AddHotel(input dto.HotelInput) (dto.HotelOutput, error)
	UpdateHotel(input dto.HotelInput) (dto.HotelOutput, error)
	DeleteHotel(hotelID int) error
}

type AdminHotelChainManagementUseCase interface {
	AddHotelChain(input dto.HotelChainInput) (dto.HotelChainOutput, error)
	UpdateHotelChain(input dto.HotelChainInput) (dto.HotelChainOutput, error)
	DeleteHotelChain(chainID int) error
}

type AdminRoomManagementUseCase interface {
	AddRoom(input dto.RoomInput) (dto.RoomOutput, error)
	UpdateRoom(input dto.RoomUpdateInput) (dto.RoomOutput, error)
	DeleteRoom(roomID int) error
}

type AdminAccountManagementUseCase interface {
	GetAccount(accountID int) (dto.AccountOutput, error)
	ListClientAccounts() ([]dto.AccountOutput, error)
	CreateClientAccount(input dto.ClientAccountInput) (dto.AccountOutput, error)
	UpdateClientAccount(accountID int, input dto.ClientAccountUpdateInput) (dto.AccountOutput, error)
	DeleteClientAccount(accountID int) error
	ListEmployeeAccounts() ([]dto.AccountOutput, error)
	CreateEmployeeAccount(input dto.EmployeeAccountInput) (dto.AccountOutput, error)
	UpdateEmployeeAccount(accountID int, input dto.EmployeeAccountUpdateInput) (dto.AccountOutput, error)
	DeleteEmployeeAccount(accountID int) error
}

// ## REPOSITORIES
// The part of the code that handles persistence (still db-technology agnostic)
// While defined in the application layer since other application code will depend on these most likely
// Their concrete implementations will likely live in framework_driven since you do need to know
// which technology to use at that point for the methods to make sense.
type EmployeeRepository interface {
	Save(emp *models.Employee) (*models.Employee, error)
	FindByID(id int) (*models.Employee, error)
	FindByEmail(email string) (*models.Employee, error)
	ListAllEmployees() ([]*models.Employee, error)
	UpdateEmployee(emp *models.Employee) (*models.Employee, error)
	UpdateManager(mgr *models.Manager) error
	Delete(employeeID int) error
}

type ClientRepository interface {
	Save(client *models.Client) (*models.Client, error)
	FindByID(id int) (*models.Client, error)
	FindByEmail(email string) (*models.Client, error)
	ListAllClients() ([]*models.Client, error)
	Update(client *models.Client) (*models.Client, error)
	Delete(id int) error
}
type HotelChainRepository interface {
	Save(chain *models.HotelChain) (*models.HotelChain, error)
	FindByID(id int) (*models.HotelChain, error)
	Update(chain *models.HotelChain) error
	Delete(id int) error
}

type HotelRepository interface {
	Save(hotel *models.Hotel) (*models.Hotel, error)
	FindByID(id int) (*models.Hotel, error)
	Update(hotel *models.Hotel) error
	Delete(id int) error
}

type RoomRepository interface {
	Save(room *models.Room) (*models.Room, error)
	FindByID(id int) (*models.Room, error)
	Update(room *models.Room) error
	Delete(id int) error
	FindAvailableRooms(hotelID int, startDate time.Time, endDate time.Time) ([]*models.Room, error)
	SearchRooms(startDate time.Time, endDate time.Time, capacity int, priceMin, priceMax float64, hotelChainID int, roomType models.RoomType) ([]*models.Room, error)
}

type ReservationRepository interface {
	Save(reservation *models.Reservation) (*models.Reservation, error)
	FindByID(id int) (*models.Reservation, error)
	GetByClient(clientID int) ([]*models.Reservation, error)
	Update(reservation *models.Reservation) error
	Delete(id int) error
}

type StayRepository interface {
	Save(stay *models.Stay) (*models.Stay, error)
	FindByID(id int) (*models.Stay, error)
	Update(stay *models.Stay) error
	EndStay(id, employeeID int) error
	Delete(id int) error
}

type QueryRepository interface {
	GetAvailableRoomsByZone() (map[string]int, error)
	GetHotelRoomCapacity(hotelId int) (int, error)
}
