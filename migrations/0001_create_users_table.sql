-- +goose Up
CREATE table if not exists users (
    id serial primary key,
    name varchar,
    age integer
);