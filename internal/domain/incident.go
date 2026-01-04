package domain

import "time"

// Incident — опасная зона в системе геооповещений
type Incident struct {
	ID          int64
	Title       string
	Description string      // описание инцидента
	Type        string      // тип: fire, flood, gas, accident
	Lat         float64
	Lon         float64
	RadiusM     int
	Active      bool

	Severity    int         // уровень опасности (1–5)
	Source      string      // кто создал: operator, system, external

	CreatedAt   time.Time
	UpdatedAt   time.Time
}
