CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    email varchar(100) UNIQUE NOT NULL,
    password bytea NOT NULL,
    username varchar(55) UNIQUE NOT NULL,
    points INTEGER NOT NULL DEFAULT 800,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
)