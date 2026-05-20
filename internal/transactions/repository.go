package transactions

import "github.com/michaelassa01/gomacbot/internal/transactions/domain"

type Repository interface {
	Create(t *domain.Transaction) (*domain.Transaction, error)
	GetByID(id string) (*domain.Transaction, error)
	GetByReference(reference string) (*domain.Transaction, error)
	ListByUserID(userID string, limit, offset int32) ([]*domain.Transaction, error)
	UpdateStatus(id, status string) (*domain.Transaction, error)
	Delete(id string) error
}
