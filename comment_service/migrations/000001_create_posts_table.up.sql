CREATE TABLE IF NOT EXISTS "comments" (
    "id" SERIAL PRIMARY KEY,
    "post_id" INTEGER,
    "text" TEXT,
    "user_id" INTEGER,
    "created_at" TIME DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIME
)