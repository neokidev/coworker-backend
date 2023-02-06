-- name: CreateMember :one
INSERT INTO members (
  id, first_name, last_name, email, status
) VALUES (
  $1, $2, $3, sqlc.narg(email), sqlc.narg(status)
)
RETURNING *;

-- name: GetMember :one
SELECT * FROM members
WHERE id = $1 LIMIT 1;

-- name: ListMembers :many
SELECT * FROM members
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateMember :one
UPDATE members
SET
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  email = COALESCE(sqlc.narg(email), email),
  status = COALESCE(sqlc.narg(status), status)
WHERE id = $1
RETURNING *;

-- name: DeleteMember :exec
DELETE FROM members
WHERE id = $1;
