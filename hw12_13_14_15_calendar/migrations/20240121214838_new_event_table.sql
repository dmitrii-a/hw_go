-- +goose Up
-- +goose StatementBegin
CREATE TABLE event
(
    id                uuid primary key,
    title             text      not null,
    start_time        timestamp not null,
    end_time          timestamp,
    notify_time       timestamp,
    description       text,
    user_id           bigint    not null,
    created_time      timestamp not null default now(),
    updated_time      timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event;
-- +goose StatementEnd
