-- +goose Up
-- +goose StatementBegin
create table if not exists workouts (
    id serial primary key,
    title varchar(255) not null,
    description text not null,
    duration_minutes integer not null,
    calories_burned integer not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists workouts;
-- +goose StatementEnd
