-- name: CreateMessage :one
INSERT INTO messages (service, title, body) VALUES ($1, $2, $3) RETURNING *;

-- name: CreateDestination :one
INSERT INTO destinations (message_id, receiver) VALUES ($1, $2) RETURNING *;
