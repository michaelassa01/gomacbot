package transactions

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/michaelassa01/gomacbot/internal/transactions/domain"
)

type CreateTransactionReq struct {
	UserID      string          `json:"user_id" binding:"required,uuid"`
	Amount      int64           `json:"amount" binding:"required,gt=0"`
	Currency    string          `json:"currency" binding:"required,len=3"`
	Type        string          `json:"type" binding:"required,oneof=credit debit transfer payment"`
	Reference   string          `json:"reference" binding:"required"`
	Description string          `json:"description"`
	Provider    string          `json:"provider"`
	Metadata    json.RawMessage `json:"metadata"`
}

type TransactionRes struct {
	ID          pgtype.UUID     `json:"id"`
	UserID      pgtype.UUID     `json:"user_id"`
	Amount      int64           `json:"amount"`
	Currency    string          `json:"currency"`
	Status      string          `json:"status"`
	Type        string          `json:"type"`
	Reference   string          `json:"reference"`
	Description string          `json:"description"`
	Provider    string          `json:"provider"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type GetTransactionReq struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type ListTransactionsReq struct {
	UserID string `form:"user_id" binding:"required,uuid"`
	Limit  int32  `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int32  `form:"offset" binding:"omitempty,min=0"`
}

type ListTransactionsRes struct {
	Transactions []TransactionRes `json:"transactions"`
}

type UpdateTransactionStatusReq struct {
	Status string `json:"status" binding:"required,oneof=pending completed failed cancelled"`
}

type DeleteTransactionReq struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func toTransactionResponse(t *domain.Transaction) *TransactionRes {
	var metadata json.RawMessage
	if len(t.Metadata) > 0 {
		metadata = json.RawMessage(t.Metadata)
	}

	return &TransactionRes{
		ID:          t.ID,
		UserID:      t.UserID,
		Amount:      t.Amount,
		Currency:    t.Currency,
		Status:      t.Status,
		Type:        t.Type,
		Reference:   t.Reference,
		Description: t.Description,
		Provider:    t.Provider,
		Metadata:    metadata,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
