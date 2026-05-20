package transactions

import (
	"errors"
	"strings"

	"github.com/michaelassa01/gomacbot/internal/transactions/domain"
	u "github.com/michaelassa01/gomacbot/utils"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(t *domain.Transaction) (*domain.Transaction, error) {
	if t.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if !u.IsSupportedCurrency(t.Currency) {
		return nil, errors.New("unsupported currency")
	}
	if t.Reference == "" {
		return nil, errors.New("reference is required")
	}
	if !isValidType(t.Type) {
		return nil, errors.New("invalid transaction type")
	}

	t.Currency = strings.ToUpper(t.Currency)
	t.Status = domain.StatusPending

	return s.repo.Create(t)
}

func (s *Service) GetByID(id string) (*domain.Transaction, error) {
	if id == "" {
		return nil, errors.New("transaction ID is required")
	}
	return s.repo.GetByID(id)
}

func (s *Service) GetByReference(reference string) (*domain.Transaction, error) {
	if reference == "" {
		return nil, errors.New("reference is required")
	}
	return s.repo.GetByReference(reference)
}

func (s *Service) ListByUserID(userID string, limit, offset int32) ([]*domain.Transaction, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByUserID(userID, limit, offset)
}

func (s *Service) UpdateStatus(id, status string) (*domain.Transaction, error) {
	if id == "" {
		return nil, errors.New("transaction ID is required")
	}
	if !isValidStatus(status) {
		return nil, errors.New("invalid transaction status")
	}
	return s.repo.UpdateStatus(id, status)
}

func (s *Service) Delete(id string) error {
	if id == "" {
		return errors.New("transaction ID is required")
	}
	return s.repo.Delete(id)
}

func isValidType(t string) bool {
	switch t {
	case domain.TypeCredit, domain.TypeDebit, domain.TypeTransfer, domain.TypePayment:
		return true
	default:
		return false
	}
}

func isValidStatus(status string) bool {
	switch status {
	case domain.StatusPending, domain.StatusCompleted, domain.StatusFailed, domain.StatusCancelled:
		return true
	default:
		return false
	}
}
