package ports

import "github.com/sql-project-backend/internal/models"

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
