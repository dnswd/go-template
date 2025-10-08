-- users
CREATE TABLE users (
    id TEXT PRIMARY KEY DEFAULT uuid_generate_v7(),
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
