package repository

import "github.com/kassse1/geo-alert-core/internal/domain"

type LocationCheckRepository interface {
	Save(check *domain.LocationCheck) error
	CountUniqueUsersSince(minutes int) (int, error)
}
