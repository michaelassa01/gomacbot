package transactions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/michaelassa01/gomacbot/internal/transactions/domain"
	db "github.com/michaelassa01/gomacbot/internal/transactions/models"
	u "github.com/michaelassa01/gomacbot/utils"
)

type PgRepo struct {
	q *db.Queries
}

func NewPgRepo(conn *pgxpool.Pool) *PgRepo {
	return &PgRepo{q: db.New(conn)}
}

func toDomain(t db.Transaction) *domain.Transaction {
	tx := &domain.Transaction{
		ID:        t.ID,
		UserID:    t.UserID,
		Amount:    t.Amount,
		Currency:  t.Currency,
		Status:    t.Status,
		Type:      t.Type,
		Reference: t.Reference,
		Metadata:  t.Metadata,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.Description != nil {
		tx.Description = *t.Description
	}
	if t.Provider != nil {
		tx.Provider = *t.Provider
	}
	return tx
}

func fromDomain(t *domain.Transaction) db.CreateTransactionParams {
	params := db.CreateTransactionParams{
		UserID:    t.UserID,
		Amount:    t.Amount,
		Currency:  t.Currency,
		Status:    t.Status,
		Type:      t.Type,
		Reference: t.Reference,
		Metadata:  t.Metadata,
	}
	if t.Description != "" {
		params.Description = &t.Description
	}
	if t.Provider != "" {
		params.Provider = &t.Provider
	}
	return params
}

func (r *PgRepo) Create(t *domain.Transaction) (*domain.Transaction, error) {
	ctx := context.Background()
	created, err := r.q.CreateTransaction(ctx, fromDomain(t))
	if err != nil {
		return nil, err
	}
	return toDomain(created), nil
}

func (r *PgRepo) GetByID(id string) (*domain.Transaction, error) {
	ctx := context.Background()
	txID, err := u.ConvertToPgUUIDFromString(id)
	if err != nil {
		return nil, err
	}

	tx, err := r.q.GetTransactionByID(ctx, txID)
	if err != nil {
		return nil, err
	}
	return toDomain(tx), nil
}

func (r *PgRepo) GetByReference(reference string) (*domain.Transaction, error) {
	ctx := context.Background()
	tx, err := r.q.GetTransactionByReference(ctx, reference)
	if err != nil {
		return nil, err
	}
	return toDomain(tx), nil
}

func (r *PgRepo) ListByUserID(userID string, limit, offset int32) ([]*domain.Transaction, error) {
	ctx := context.Background()
	uid, err := u.ConvertToPgUUIDFromString(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListTransactionsByUserID(ctx, db.ListTransactionsByUserIDParams{
		UserID: uid,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	items := make([]*domain.Transaction, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomain(row))
	}
	return items, nil
}

func (r *PgRepo) UpdateStatus(id, status string) (*domain.Transaction, error) {
	ctx := context.Background()
	txID, err := u.ConvertToPgUUIDFromString(id)
	if err != nil {
		return nil, err
	}

	updated, err := r.q.UpdateTransactionStatus(ctx, db.UpdateTransactionStatusParams{
		ID:     txID,
		Status: status,
	})
	if err != nil {
		return nil, err
	}
	return toDomain(updated), nil
}

func (r *PgRepo) Delete(id string) error {
	ctx := context.Background()
	txID, err := u.ConvertToPgUUIDFromString(id)
	if err != nil {
		return err
	}
	return r.q.DeleteTransaction(ctx, txID)
}
