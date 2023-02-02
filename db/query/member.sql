-- name: CreateMember :one
INSERT INTO members (
  id, first_name, last_name, email
) VALUES (
  $1, $2, $3, $4
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
