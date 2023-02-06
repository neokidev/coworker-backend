CREATE TABLE "members"
(
    "id"         uuid PRIMARY KEY NOT NULL,
    "first_name" varchar          NOT NULL,
    "last_name"  varchar          NOT NULL,
    "email"      varchar,
    "created_at" timestamptz      NOT NULL DEFAULT (now())
);
