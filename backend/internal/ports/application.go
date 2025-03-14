package ports

import "github.com/sql-project-backend/internal/models"

// ## REPOSITORIES
// The part of the code that handles persistence (still db-technology agnostic)
// While defined in the application layer since other application code will depend on these most likely
// Their concrete implementations will likely live in framework_driven since you do need to know
// which technology to use at that point for the methods to make sense.
type EmployeeRepository interface {
	Save(emp *models.Employee) error
	FindByID(id int) (*models.Employee, error)
	UpdateManager(mgr *models.Manager) error
	Delete(employeeID int) error
}

type ClientRepository interface {
	Save(client *models.Client) error
	FindByID(id int) (*models.Client, error)
	Update(client *models.Client) error
	Delete(id int) error
}

type HotelChainRepository interface {
	Save(chain *models.HotelChain) error
	FindByID(id int) (*models.HotelChain, error)
	Update(chain *models.HotelChain) error
	Delete(id int) error
}

type HotelRepository interface {
	Save(hotel *models.Hotel) error
	FindByID(id int) (*models.Hotel, error)
	Update(hotel *models.Hotel) error
	Delete(id int) error
}

type RoomRepository interface {
	Save(room *models.Room) error
	FindByID(id int) (*models.Room, error)
	Update(room *models.Room) error
	Delete(id int) error
}

type ReservationRepository interface {
	Save(reservation *models.Reservation) error
	FindByID(id int) (*models.Reservation, error)
	Update(reservation *models.Reservation) error
	Delete(id int) error
}

type StayRepository interface {
	Save(stay *models.Stay) error
	FindByID(id int) (*models.Stay, error)
	Update(stay *models.Stay) error
	Delete(id int) error
}

type QueryRepository interface {
	GetAvailableRoomsByZone() (map[string]int, error)
	GetHotelRoomCapacity(hotelId int) (int, error)
}
