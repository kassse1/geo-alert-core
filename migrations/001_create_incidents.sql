CREATE TABLE incidents (
                           id BIGSERIAL PRIMARY KEY,
                           title TEXT NOT NULL,
                           lat DOUBLE PRECISION NOT NULL,
                           lon DOUBLE PRECISION NOT NULL,
                           radius_m INTEGER NOT NULL,
                           active BOOLEAN NOT NULL DEFAULT TRUE,
                           created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_incidents_active ON incidents(active);
