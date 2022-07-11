-- name: GetUser :one
SELECT * FROM users
WHERE id = pggen.arg('Id') LIMIT 1;

-- name: GetUserByEmail :many
SELECT * FROM users
WHERE email = pggen.arg('Email') LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = pggen.arg('Id');