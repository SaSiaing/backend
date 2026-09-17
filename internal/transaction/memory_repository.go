package transaction

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	nextID uint64
	items  map[string]Transaction
	now    func() time.Time
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: make(map[string]Transaction),
		now:   time.Now,
	}
}

func (repository *MemoryRepository) Create(ctx context.Context, householdID string, input CreateInput) (Transaction, error) {
	if err := ctx.Err(); err != nil {
		return Transaction{}, err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()

	repository.nextID++
	now := repository.now().UTC()
	transaction := Transaction{
		ID:          fmt.Sprintf("tx-%d", repository.nextID),
		HouseholdID: householdID,
		PayerID:     input.PayerID,
		CategoryID:  input.CategoryID,
		Type:        input.Type,
		Amount:      input.Amount,
		OccurredAt:  input.OccurredAt,
		Memo:        input.Memo,
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	repository.items[transaction.ID] = transaction
	return transaction, nil
}

func (repository *MemoryRepository) List(ctx context.Context, householdID string, filter ListFilter) ([]Transaction, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	repository.mu.RLock()
	defer repository.mu.RUnlock()

	transactions := make([]Transaction, 0)
	for _, transaction := range repository.items {
		if transaction.HouseholdID != householdID || !transaction.DeletedAt.IsZero() {
			continue
		}
		if filter.From != nil && transaction.OccurredAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && transaction.OccurredAt.After(*filter.To) {
			continue
		}
		if filter.PayerID != "" && transaction.PayerID != filter.PayerID {
			continue
		}
		if filter.CategoryID != "" && transaction.CategoryID != filter.CategoryID {
			continue
		}
		if filter.Type != "" && transaction.Type != filter.Type {
			continue
		}
		transactions = append(transactions, transaction)
	}

	sort.Slice(transactions, func(i, j int) bool {
		if transactions[i].OccurredAt.Equal(transactions[j].OccurredAt) {
			return transactions[i].CreatedAt.After(transactions[j].CreatedAt)
		}
		return transactions[i].OccurredAt.After(transactions[j].OccurredAt)
	})
	return transactions, nil
}

func (repository *MemoryRepository) Get(ctx context.Context, householdID, transactionID string) (Transaction, error) {
	if err := ctx.Err(); err != nil {
		return Transaction{}, err
	}

	repository.mu.RLock()
	defer repository.mu.RUnlock()

	transaction, ok := repository.items[transactionID]
	if !ok || transaction.HouseholdID != householdID || !transaction.DeletedAt.IsZero() {
		return Transaction{}, ErrNotFound
	}
	return transaction, nil
}

func (repository *MemoryRepository) Update(ctx context.Context, householdID, transactionID string, input UpdateInput) (Transaction, error) {
	if err := ctx.Err(); err != nil {
		return Transaction{}, err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()

	transaction, ok := repository.items[transactionID]
	if !ok || transaction.HouseholdID != householdID || !transaction.DeletedAt.IsZero() {
		return Transaction{}, ErrNotFound
	}
	if transaction.Version != input.Version {
		return Transaction{}, ErrConflict
	}

	if input.PayerID != nil {
		transaction.PayerID = *input.PayerID
	}
	if input.CategoryID != nil {
		transaction.CategoryID = *input.CategoryID
	}
	if input.Type != nil {
		transaction.Type = *input.Type
	}
	if input.Amount != nil {
		transaction.Amount = *input.Amount
	}
	if input.OccurredAt != nil {
		transaction.OccurredAt = *input.OccurredAt
	}
	if input.Memo != nil {
		transaction.Memo = *input.Memo
	}
	transaction.Version++
	transaction.UpdatedAt = repository.now().UTC()
	repository.items[transactionID] = transaction
	return transaction, nil
}

func (repository *MemoryRepository) Delete(ctx context.Context, householdID, transactionID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()

	transaction, ok := repository.items[transactionID]
	if !ok || transaction.HouseholdID != householdID || !transaction.DeletedAt.IsZero() {
		return ErrNotFound
	}

	now := repository.now().UTC()
	transaction.DeletedAt = now
	transaction.UpdatedAt = now
	transaction.Version++
	repository.items[transactionID] = transaction
	return nil
}
