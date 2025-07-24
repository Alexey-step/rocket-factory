-- +goose UP
create table orders
(
    id serial primary key,
    uuid uuid not null unique default uuid_generate_v4(),
    user_uuid uuid not null,
	part_uuids uuid[] not null,
    total_price double precision not null,
    transaction_uuid uuid,
    payment_method text,
    status text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp
);

-- +goose Down
drop table if exists orders;