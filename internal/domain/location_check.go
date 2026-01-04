package domain

import "time"

type LocationCheck struct {
	ID           int64
	UserID       string
	Lat          float64
	Lon          float64

	IncidentIDs  []int64    // какие инциденты были найдены
	HasDanger    bool       // был ли риск
	DistanceM    int        // минимальная дистанция до зоны

	CheckedAt    time.Time
}

