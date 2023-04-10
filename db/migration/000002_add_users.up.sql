CREATE
EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE "users"
(
    "id"                  uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "first_name"          varchar          NOT NULL,
    "last_name"           varchar          NOT NULL,
    "email"               varchar UNIQUE   NOT NULL,
    "hashed_password"     varchar          NOT NULL,
    "password_changed_at" timestamptz      NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at"          timestamptz               DEFAULT (now())
);
