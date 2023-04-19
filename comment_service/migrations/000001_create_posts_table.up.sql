CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER,
    post_title VARCHAR(60),   
    user_id INTEGER,
    post_user_name VARCHAR(60),
    text text,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIME
);