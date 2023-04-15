CREATE TABLE "sessions"
(
    "id"            uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "user_id"       uuid NOT NULL,
    "session_token" uuid NOT NULL,
    "expired_at"    timestamptz NOT NULL
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
