-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role INT NOT NULL,
    created_at TIMESTAMP NOT NULL default now(),
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE users;