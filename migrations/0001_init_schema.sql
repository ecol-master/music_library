-- migrations/0001_init_schema.sql
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    "group" TEXT NOT NULL,
    song TEXT NOT NULL
);