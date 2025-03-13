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
