package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store store providers all functions to execute db queries and transactions
type Store interface {
	Querier
}

// SQLStore providers all functions to execute SQL DB queries and transactions
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new store instance
func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
