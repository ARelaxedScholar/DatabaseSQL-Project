package ports

import "github.com/sql-project-backend/internal/models"

type EmployeeRepository interface {
	Save(emp *models.Employee) error
	FindByID(id int) (*models.Employee, error)
	UpdateManager(mgr *models.Manager) error
	Delete(employeeID int) error
}
