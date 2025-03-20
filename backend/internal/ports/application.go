package ports

import (
	"github.com/sql-project-backend/internal/models"
	"github.com/sql-project-backend/internal/models/dto"
)

//import "github.com/sql-project-backend/internal/models/dto"

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
	CancelReservation(reservationID int) error
}

type ClientProfileManagementUseCase interface {
	GetProfile(clientID int) (dto.ClientProfileOutput, error)
	UpdateProfile(input dto.ClientProfileUpdateInput) (dto.ClientProfileOutput, error)
}

// ## Employee USE CASES
type EmployeeRegistrationUseCase interface {
	RegisterEmployee(input dto.EmployeeRegistrationInput) (dto.EmployeeRegistrationOutput, error)
}

type EmployeeLoginUseCase interface {
	Login(input dto.EmployeeLoginInput) (dto.EmployeeLoginOutput, error)
}

type EmployeeCheckInUseCase interface {
	CheckIn(input dto.CheckInInput) (dto.CheckInOutput, error)
}

type EmployeeCreateNewStayUseCase interface {
	CreateNewStay(input dto.NewStayInput) (dto.NewStayOutput, error)
}

// ## Anonymous
type SearchRoomsUseCase interface {
	SearchRooms(input dto.RoomSearchInput) (dto.RoomSearchOutput, error)
}

// ## Admin USE CASES (Right now no requirement for that so kind of an after thought)
type AdminHotelManagementUseCase interface {
	AddHotel(input dto.HotelInput) (dto.HotelOutput, error)
	UpdateHotel(input dto.HotelInput) (dto.HotelOutput, error)
	DeleteHotel(hotelID int) error
}

type AdminHotelChainUseManagementCase interface {
	AddHotelChain(input dto.HotelChainInput) (dto.HotelChainOutput, error)
	UpdateHotelChain(input dto.HotelChainInput) (dto.HotelChainOutput, error)
	DeleteHotelChain(chainID int) error
}

type AdminRoomManagementCase interface {
	AddRoom(input dto.RoomInput) (dto.RoomOutput, error)
	UpdateRoom(input dto.RoomUpdateInput) (dto.RoomOutput, error)
	DeleteRoom(roomID int) error
}

// ## REPOSITORIES
// The part of the code that handles persistence (still db-technology agnostic)
// While defined in the application layer since other application code will depend on these most likely
// Their concrete implementations will likely live in framework_driven since you do need to know
// which technology to use at that point for the methods to make sense.
type EmployeeRepository interface {
	Save(emp *models.Employee) (*models.Employee, error)
	FindByID(id int) (*models.Employee, error)
	UpdateManager(mgr *models.Manager) error
	Delete(employeeID int) error
}

type ClientRepository interface {
	Save(client *models.Client) (*models.Client, error)
	FindByID(id int) (*models.Client, error)
	FindByEmail(email string) (*models.Client, error)
	Update(client *models.Client) error
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
	Delete(id int) error
}

type QueryRepository interface {
	GetAvailableRoomsByZone() (map[string]int, error)
	GetHotelRoomCapacity(hotelId int) (int, error)
}
