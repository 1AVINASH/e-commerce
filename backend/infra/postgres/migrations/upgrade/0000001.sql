create table users (
    id bigserial,
    name text,
    email text,
    password text,
    session_token text,
    role text
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

create table product_photos (
    id bigserial,
    product_id bigint,
    path text
)