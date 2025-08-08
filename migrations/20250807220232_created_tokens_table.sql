-- +goose Up
-- +goose StatementBegin
create table if not exists tokens (
    Hash bytea primary key,
    user_id integer not null references users(id) on delete cascade,
    scope varchar(255) not null,
    expires_at timestamp with time zone not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tokens;
-- +goose StatementEnd
