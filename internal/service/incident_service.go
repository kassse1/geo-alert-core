package service

import (
	"errors"

	"github.com/kassse1/geo-alert-core/internal/domain"
	"github.com/kassse1/geo-alert-core/internal/repository"
)

type IncidentService struct {
	repo      repository.IncidentRepository
	checkRepo repository.LocationCheckRepository
}

func NewIncidentService(
	repo repository.IncidentRepository,
	checkRepo repository.LocationCheckRepository,
) *IncidentService {
	return &IncidentService{
		repo:      repo,
		checkRepo: checkRepo,
	}
}

// =====================
// CRUD
// =====================

func (s *IncidentService) Create(incident *domain.Incident) error {
	if incident == nil {
		return errors.New("incident is nil")
	}
	return s.repo.Create(incident)
}

func (s *IncidentService) List(page, limit int) ([]domain.Incident, error) {
	offset := (page - 1) * limit
	return s.repo.List(limit, offset)
}

func (s *IncidentService) GetByID(id int64) (*domain.Incident, error) {
	return s.repo.GetByID(id)
}

func (s *IncidentService) Update(incident *domain.Incident) error {
	if incident == nil {
		return errors.New("incident is nil")
	}
	return s.repo.Update(incident)
}

func (s *IncidentService) Deactivate(id int64) error {
	return s.repo.Deactivate(id)
}

// =====================
// Stats
// =====================

func (s *IncidentService) Stats(minutes int) (int, error) {
	if minutes <= 0 {
		return 0, errors.New("minutes must be positive")
	}
	return s.checkRepo.CountUniqueUsersSince(minutes)
}
