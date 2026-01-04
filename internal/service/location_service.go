package service

import (
	"github.com/kassse1/geo-alert-core/internal/domain"
	"github.com/kassse1/geo-alert-core/internal/repository"
)

type LocationService struct {
	incidentRepo repository.IncidentRepository
	checkRepo    repository.LocationCheckRepository
}

func NewLocationService(
	incidentRepo repository.IncidentRepository,
	checkRepo repository.LocationCheckRepository,
) *LocationService {
	return &LocationService{
		incidentRepo: incidentRepo,
		checkRepo:    checkRepo,
	}
}

func (s *LocationService) CheckLocation(
	userID string,
	lat, lon float64,
) ([]domain.Incident, error) {

	// 1️⃣ Получаем активные инциденты
	incidents, err := s.incidentRepo.GetActive()
	if err != nil {
		return nil, err
	}

	// 2️⃣ Фильтруем по радиусу
	nearby := make([]domain.Incident, 0)
	for _, i := range incidents {
		distance := DistanceMeters(lat, lon, i.Lat, i.Lon)
		if distance <= float64(i.RadiusM) {
			nearby = append(nearby, i)
		}
	}

	// 3️⃣ Сохраняем факт проверки
	_ = s.checkRepo.Save(&domain.LocationCheck{
		UserID: userID,
		Lat:    lat,
		Lon:    lon,
	})

	return nearby, nil
}
