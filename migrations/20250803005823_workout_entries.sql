-- +goose Up
-- +goose StatementBegin
create table if not exists workout_entries (
    id serial primary key,
    workout_id bigint not null references workouts(id) ON DELETE CASCADE,
    exercise_name varchar(255) not null,
    sets integer not null,
    reps integer,
    duration_seconds integer,
    weight decimal(5, 2) not null,
    notes text,
    order_index integer not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()

    constraint valid_workout_entry check (
        (reps is not null or duration_seconds is not null) and
        (reps is null or duration_seconds is null)
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists workout_entries;
-- +goose StatementEnd
