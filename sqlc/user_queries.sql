-- name: CreateUser :one
INSERT INTO users (
    login, password
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UserByLogin :one
SELECT * FROM users
WHERE login = $1;

-- name: DeleteUser :exec
DELETE FROM users
    WHERE id = $1;
