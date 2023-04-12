CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    password varchar,
    created_at TIME DEFAULT CURRENT_TIMESTAMP,
    updated_at TIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIME
);