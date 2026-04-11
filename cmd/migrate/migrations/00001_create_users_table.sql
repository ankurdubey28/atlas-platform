-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name VARCHAR(50) NOT NULL,
        age INT NOT NULL,
        email CITEXT NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
