-- users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    balance NUMERIC NOT NULL
);

-- CREATE TABLE accounts (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
--     email TEXT NOT NULL UNIQUE,
--     name TEXT NOT NULL,
--     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
--     belongs_to UUID NOT NULL
-- );
