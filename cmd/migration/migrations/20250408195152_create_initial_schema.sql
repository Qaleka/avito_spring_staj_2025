-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL,
    hash_password TEXT NOT NULL,
    role TEXT CHECK (role IN ('employee','moderator')) NOT NULL
);

CREATE TABLE pvzs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registration_date TIMESTAMP NOT NULL,
    city TEXT NOT NULL
);

CREATE TABLE receptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_time TIMESTAMP NOT NULL,
    pvz_id UUID NOT NULL,
    status TEXT CHECK (status IN ('in_progress', 'close')) NOT NULL,
    FOREIGN KEY (pvz_id) REFERENCES pvzs(id)
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    received_at TIMESTAMP NOT NULL,
    type TEXT CHECK (type IN ('электроника', 'одежда', 'обувь')) NOT NULL,
    date_time UUID NOT NULL,
    FOREIGN KEY (reception_id) REFERENCES receptions(id)
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS receptions;
DROP TABLE IF EXISTS pickup_points;
