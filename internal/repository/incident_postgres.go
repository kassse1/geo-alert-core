package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/kassse1/geo-alert-core/internal/domain"
)

type IncidentPostgresRepository struct {
	db *sql.DB
}

func NewIncidentPostgresRepository(db *sql.DB) *IncidentPostgresRepository {
	return &IncidentPostgresRepository{db: db}
}

func (r *IncidentPostgresRepository) Create(i *domain.Incident) error {
	query := `
		INSERT INTO incidents (title, lat, lon, radius_m)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.db.QueryRowContext(
		ctx,
		query,
		i.Title,
		i.Lat,
		i.Lon,
		i.RadiusM,
	).Scan(&i.ID, &i.CreatedAt)
}

func (r *IncidentPostgresRepository) GetByID(id int64) (*domain.Incident, error) {
	query := `
		SELECT id, title, lat, lon, radius_m, active, created_at
		FROM incidents
		WHERE id = $1
	`

	var i domain.Incident

	err := r.db.QueryRow(query, id).Scan(
		&i.ID,
		&i.Title,
		&i.Lat,
		&i.Lon,
		&i.RadiusM,
		&i.Active,
		&i.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (r *IncidentPostgresRepository) List(offset, limit int) ([]domain.Incident, error) {
	query := `
		SELECT id, title, lat, lon, radius_m, active, created_at
		FROM incidents
		ORDER BY id
		OFFSET $1 LIMIT $2
	`

	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []domain.Incident

	for rows.Next() {
		var i domain.Incident
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Lat,
			&i.Lon,
			&i.RadiusM,
			&i.Active,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, i)
	}

	return incidents, nil
}

func (r *IncidentPostgresRepository) Update(i *domain.Incident) error {
	query := `
		UPDATE incidents
		SET title = $1, lat = $2, lon = $3, radius_m = $4, active = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(
		query,
		i.Title,
		i.Lat,
		i.Lon,
		i.RadiusM,
		i.Active,
		i.ID,
	)

	return err
}

func (r *IncidentPostgresRepository) Deactivate(id int64) error {
	query := `
		UPDATE incidents
		SET active = FALSE
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	return err
}

func (r *IncidentPostgresRepository) GetActive() ([]domain.Incident, error) {
	query := `
		SELECT id, title, lat, lon, radius_m, active, created_at
		FROM incidents
		WHERE active = TRUE
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []domain.Incident

	for rows.Next() {
		var i domain.Incident
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Lat,
			&i.Lon,
			&i.RadiusM,
			&i.Active,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, i)
	}

	return incidents, nil
}
