CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER,
    post_title VARCHAR(60),   
    user_id INTEGER,
    post_user_name VARCHAR(60),
    text text,
    created_at TIME DEFAULT CURRENT_TIMESTAMP,
    updated_at TIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIME
);