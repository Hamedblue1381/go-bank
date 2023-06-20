-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    curreny
) VALUES (
    $1, $2, $3
) RETURNING *;