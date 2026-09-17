package transaction

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("transaction not found")
	ErrConflict = errors.New("transaction version conflict")
)

type Repository interface {
	Create(context.Context, string, CreateInput) (Transaction, error)
	List(context.Context, string, ListFilter) ([]Transaction, error)
	Get(context.Context, string, string) (Transaction, error)
	Update(context.Context, string, string, UpdateInput) (Transaction, error)
	Delete(context.Context, string, string) error
}
