package domain

import "time"

type Incident struct {
	ID        int64
	Title     string
	Lat       float64
	Lon       float64
	RadiusM   int
	Active    bool
	CreatedAt time.Time
}
