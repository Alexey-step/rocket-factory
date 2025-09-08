-- +goose UP
create table if not exists users
(
    id bigint generated always as identity primary key,
    uuid uuid not null unique default uuid_generate_v4(),
    login text not null unique,
    email text not null unique,
    notification_methods jsonb not null default '[]'::jsonb,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp,
    password_hash text not null
);

-- +goose Down
drop table if exists users;