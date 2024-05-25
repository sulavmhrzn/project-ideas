CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    username text not null unique,
    email text not null unique,
    hash_password text not null,
    created_at timestamptz not null default now()
);

