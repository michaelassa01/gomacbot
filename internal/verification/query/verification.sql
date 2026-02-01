-- name: CreateVerification :one
INSERT INTO verification (
  code,
  type,
  identifier,
  expired_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetVerificationByID :one
SELECT *
FROM verification
WHERE id = $1
LIMIT 1;


-- name: GetValidVerification :one
SELECT *
FROM verification
WHERE
  code = $1
  AND type = $2
  AND identifier = $3
  AND expired_at > now()
LIMIT 1;


-- name: ListVerificationsByIdentifier :many
SELECT *
FROM verification
WHERE identifier = $1
ORDER BY created_at DESC;


-- name: DeleteVerification :exec
DELETE FROM verification
WHERE id = $1;


-- name: DeleteExpiredVerifications :exec
DELETE FROM verification
WHERE expired_at <= now();


-- name: DeleteVerificationByIdentifierAndType :exec
DELETE FROM verification
WHERE identifier = $1
  AND type = $2;


