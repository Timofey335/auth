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

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    endpoint VARCHAR(255) NOT NULL UNIQUE,
    role INT NOT NULL
);

INSERT INTO roles (endpoint, role) VALUES ('/chat_server_v1.Chat_server_v1/CreateChat', 1);
INSERT INTO roles (endpoint, role) VALUES ('/chat_server_v1.Chat_server_v1/DeleteChat', 1);
INSERT INTO roles (endpoint, role) VALUES ('/chat_server_v1.Chat_server_v1/SendMessage', 0);

-- +goose Down
DROP TABLE users, roles;