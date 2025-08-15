create table users (
    id bigserial,
    name text,
    email text,
    password text,
    session_token text
);

create table products (
    id bigserial,
    title text,
    description text,
    original_price float,
    discount float,
    thumbnail text,
    photos text
)