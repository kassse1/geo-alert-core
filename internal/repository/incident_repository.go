package repository

import "github.com/kassse1/geo-alert-core/internal/domain"

// IncidentRepository описывает работу с инцидентами в хранилище
type IncidentRepository interface {
	Create(incident *domain.Incident) error
	GetByID(id int64) (*domain.Incident, error)
	List(offset, limit int) ([]domain.Incident, error)
	Update(incident *domain.Incident) error
	Deactivate(id int64) error
	GetActive() ([]domain.Incident, error)
}
