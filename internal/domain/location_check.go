package domain

import "time"

type LocationCheck struct {
	ID        int64
	UserID    string
	Lat       float64
	Lon       float64
	CheckedAt time.Time
}
