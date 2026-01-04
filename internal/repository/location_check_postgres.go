package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/kassse1/geo-alert-core/internal/domain"
)

type LocationCheckPostgresRepository struct {
	db *sql.DB
}

func NewLocationCheckPostgresRepository(db *sql.DB) *LocationCheckPostgresRepository {
	return &LocationCheckPostgresRepository{db: db}
}
func (r *LocationCheckPostgresRepository) CountUniqueUsersSince(minutes int) (int, error) {
	query := `
		SELECT COUNT(DISTINCT user_id)
		FROM location_checks
		WHERE checked_at >= NOW() - make_interval(mins => $1)
	`

	var count int
	err := r.db.QueryRow(query, minutes).Scan(&count)
	return count, err
}

func (r *LocationCheckPostgresRepository) Save(c *domain.LocationCheck) error {
	query := `
		INSERT INTO location_checks (user_id, lat, lon)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, c.UserID, c.Lat, c.Lon)
	return err
}
