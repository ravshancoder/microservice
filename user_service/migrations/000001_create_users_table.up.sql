CREATE TABLE IF NOT EXISTS users (
    id            uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name    VARCHAR(50) NOT NULL,
    last_name     VARCHAR(50) NOT NULL,
    email         VARCHAR(100) NOT NULL UNIQUE,
    user_type     VARCHAR(50),
    password      TEXT NOT NULL,
    acces_token   TEXT,
    refresh_token TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIME
);