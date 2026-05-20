-- name: CreateTransaction :one
INSERT INTO transactions (
  user_id,
  amount,
  currency,
  status,
  type,
  reference,
  description,
  provider,
  metadata
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetTransactionByID :one
SELECT *
FROM transactions
WHERE id = $1
LIMIT 1;

-- name: GetTransactionByReference :one
SELECT *
FROM transactions
WHERE reference = $1
LIMIT 1;

-- name: ListTransactionsByUserID :many
SELECT *
FROM transactions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET
  status = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;
