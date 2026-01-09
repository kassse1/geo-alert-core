package domain

import "time"

// Incident — опасная зона в системе геооповещений
type Incident struct {
	ID        int64
	Title     string
	Lat       float64
	Lon       float64
	RadiusM   int
	Active    bool
	CreatedAt time.Time
}
