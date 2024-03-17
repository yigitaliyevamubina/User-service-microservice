CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL,
    first_name TEXT NO NULL,
    last_name TEXT NOT NULL,
    age INT,
    gender INT,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    refresh_token NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
