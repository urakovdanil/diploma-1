-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id       bigserial CONSTRAINT users_pk PRIMARY KEY,
    login    varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    CONSTRAINT unique_login UNIQUE (login),
    CONSTRAINT unique_login_password UNIQUE (login, password)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users
-- +goose StatementEnd