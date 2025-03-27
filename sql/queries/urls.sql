-- name: CreateEntry :exec
INSERT INTO urls (
    original_url, hashed_url, created_at, ttl)
    VALUES (
        $1, $2, NOW(), $3
    )
;

-- name: GetUrlbyOrig :one
SELECT hashed_url FROM urls WHERE original_url=$1;

-- name: GetUrlbyHash :one
SELECT original_url FROM urls WHERE hashed_url=$1;

-- name: DeleteTtl :exec
DELETE FROM urls WHERE ttl<NOW()::time;