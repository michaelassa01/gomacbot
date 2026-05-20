package domain

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"

	TypeCredit   = "credit"
	TypeDebit    = "debit"
	TypeTransfer = "transfer"
	TypePayment  = "payment"
)

type Transaction struct {
	ID          pgtype.UUID `json:"id"`
	UserID      pgtype.UUID `json:"user_id"`
	Amount      int64       `json:"amount"`
	Currency    string      `json:"currency"`
	Status      string      `json:"status"`
	Type        string      `json:"type"`
	Reference   string      `json:"reference"`
	Description string      `json:"description"`
	Provider    string      `json:"provider"`
	Metadata    []byte      `json:"metadata"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
