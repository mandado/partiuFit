-- +goose Up
-- +goose StatementBegin
alter table workouts add column user_id integer not null references users(id) on delete cascade;
alter table workout_entries add column user_id integer not null references users(id) on delete cascade;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table workouts drop column if exists user_id;
alter table workout_entries drop column if exists user_id;
-- +goose StatementEnd
