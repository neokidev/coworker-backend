CREATE TYPE "member_statuses" AS ENUM (
  'active',
  'offline'
);

CREATE TABLE "members" (
  "id" uuid PRIMARY KEY NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar,
  "status" member_statuses NOT NULL DEFAULT ('offline'),
  "created_at" timestamptz DEFAULT (now())
);
