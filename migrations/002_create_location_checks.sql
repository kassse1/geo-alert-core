CREATE TABLE location_checks (
                                 id BIGSERIAL PRIMARY KEY,
                                 user_id TEXT NOT NULL,
                                 lat DOUBLE PRECISION NOT NULL,
                                 lon DOUBLE PRECISION NOT NULL,
                                 checked_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_location_checks_checked_at ON location_checks(checked_at);
CREATE INDEX idx_location_checks_user_id ON location_checks(user_id);
