package service

import (
	"github.com/kassse1/geo-alert-core/internal/domain"
	"github.com/kassse1/geo-alert-core/internal/repository"
)

type LocationService struct {
	incidentRepo repository.IncidentRepository
	checkRepo    repository.LocationCheckRepository
	webhook      *WebhookService
}

func NewLocationService(
	incidentRepo repository.IncidentRepository,
	checkRepo repository.LocationCheckRepository,
	webhook *WebhookService,
) *LocationService {
	return &LocationService{
		incidentRepo: incidentRepo,
		checkRepo:    checkRepo,
		webhook:      webhook,
	}
}

func (s *LocationService) CheckLocation(
	userID string,
	lat, lon float64,
) ([]domain.Incident, error) {

	//  Получаем только АКТИВНЫЕ инциденты
	incidents, err := s.incidentRepo.GetActive()
	if err != nil {
		return nil, err
	}

	//  Фильтруем по расстоянию
	nearby := make([]domain.Incident, 0)
	for _, i := range incidents {
		distance := DistanceMeters(lat, lon, i.Lat, i.Lon)
		if distance <= float64(i.RadiusM) {
			nearby = append(nearby, i)
		}
	}

	//  Сохраняем факт проверки (не блокирует ответ)
	_ = s.checkRepo.Save(&domain.LocationCheck{
		UserID: userID,
		Lat:    lat,
		Lon:    lon,
	})

	//  Асинхронно отправляем webhook, если есть опасности
	if len(nearby) > 0 && s.webhook != nil {
		go s.webhook.Send(userID, nearby)
	}

	return nearby, nil
}
