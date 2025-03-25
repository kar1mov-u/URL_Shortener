-- name: CreateEntry :exec
INSERT INTO urls (
    original_url, hashed_url, created_at)
    VALUES (
        $1, $2, NOW()
    )
;

-- name: GetUrlbyOrig :one
SELECT hashed_url FROM urls WHERE original_url=$1;

-- name: GetUrlbyHash :one
SELECT original_url FROM urls WHERE hashed_url=$1;