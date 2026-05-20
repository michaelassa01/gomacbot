package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	transactions "github.com/michaelassa01/gomacbot/internal/transactions"
	users "github.com/michaelassa01/gomacbot/internal/users"
)

type Services struct {
	User         *users.Service
	Transactions *transactions.Service
}

func NewServices(dbConn *pgxpool.Pool) *Services {
	return &Services{
		User:         users.NewService(users.NewPgRepo(dbConn)),
		Transactions: transactions.NewService(transactions.NewPgRepo(dbConn)),
	}
}
