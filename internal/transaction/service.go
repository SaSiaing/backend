package transaction

import (
	"context"
	"strings"
)

type ValidationError struct {
	Message string
}

func (err ValidationError) Error() string {
	return err.Message
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (service Service) Create(ctx context.Context, householdID string, input CreateInput) (Transaction, error) {
	if err := validateHouseholdID(householdID); err != nil {
		return Transaction{}, err
	}

	input.PayerID = strings.TrimSpace(input.PayerID)
	input.CategoryID = strings.TrimSpace(input.CategoryID)
	input.Memo = strings.TrimSpace(input.Memo)
	input.OccurredAt = input.OccurredAt.UTC()
	if err := validateCreate(input); err != nil {
		return Transaction{}, err
	}
	return service.repository.Create(ctx, householdID, input)
}

func (service Service) List(ctx context.Context, householdID string, filter ListFilter) ([]Transaction, error) {
	if err := validateHouseholdID(householdID); err != nil {
		return nil, err
	}
	filter.PayerID = strings.TrimSpace(filter.PayerID)
	filter.CategoryID = strings.TrimSpace(filter.CategoryID)
	if filter.From != nil {
		from := filter.From.UTC()
		filter.From = &from
	}
	if filter.To != nil {
		to := filter.To.UTC()
		filter.To = &to
	}
	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		return nil, ValidationError{Message: "from must be before or equal to to"}
	}
	if filter.Type != "" && !filter.Type.Valid() {
		return nil, ValidationError{Message: "type must be expense or income"}
	}
	return service.repository.List(ctx, householdID, filter)
}

func (service Service) Get(ctx context.Context, householdID, transactionID string) (Transaction, error) {
	if err := validateIDs(householdID, transactionID); err != nil {
		return Transaction{}, err
	}
	return service.repository.Get(ctx, householdID, transactionID)
}

func (service Service) Update(ctx context.Context, householdID, transactionID string, input UpdateInput) (Transaction, error) {
	if err := validateIDs(householdID, transactionID); err != nil {
		return Transaction{}, err
	}
	if input.Version < 1 {
		return Transaction{}, ValidationError{Message: "version must be a positive integer"}
	}
	if !input.HasChanges() {
		return Transaction{}, ValidationError{Message: "at least one field must be provided"}
	}

	if input.PayerID != nil {
		payerID := strings.TrimSpace(*input.PayerID)
		if payerID == "" {
			return Transaction{}, ValidationError{Message: "payer_id is required"}
		}
		input.PayerID = &payerID
	}
	if input.CategoryID != nil {
		categoryID := strings.TrimSpace(*input.CategoryID)
		if categoryID == "" {
			return Transaction{}, ValidationError{Message: "category_id is required"}
		}
		input.CategoryID = &categoryID
	}
	if input.Type != nil && !input.Type.Valid() {
		return Transaction{}, ValidationError{Message: "type must be expense or income"}
	}
	if input.Amount != nil && *input.Amount < 1 {
		return Transaction{}, ValidationError{Message: "amount must be greater than zero"}
	}
	if input.OccurredAt != nil {
		if input.OccurredAt.IsZero() {
			return Transaction{}, ValidationError{Message: "occurred_at is required"}
		}
		occurredAt := input.OccurredAt.UTC()
		input.OccurredAt = &occurredAt
	}
	if input.Memo != nil {
		memo := strings.TrimSpace(*input.Memo)
		input.Memo = &memo
	}

	return service.repository.Update(ctx, householdID, transactionID, input)
}

func (service Service) Delete(ctx context.Context, householdID, transactionID string) error {
	if err := validateIDs(householdID, transactionID); err != nil {
		return err
	}
	return service.repository.Delete(ctx, householdID, transactionID)
}

func (transactionType Type) Valid() bool {
	return transactionType == TypeExpense || transactionType == TypeIncome
}

func (input UpdateInput) HasChanges() bool {
	return input.PayerID != nil || input.CategoryID != nil || input.Type != nil || input.Amount != nil || input.OccurredAt != nil || input.Memo != nil
}

func validateCreate(input CreateInput) error {
	if input.PayerID == "" {
		return ValidationError{Message: "payer_id is required"}
	}
	if input.CategoryID == "" {
		return ValidationError{Message: "category_id is required"}
	}
	if !input.Type.Valid() {
		return ValidationError{Message: "type must be expense or income"}
	}
	if input.Amount < 1 {
		return ValidationError{Message: "amount must be greater than zero"}
	}
	if input.OccurredAt.IsZero() {
		return ValidationError{Message: "occurred_at is required"}
	}
	return nil
}

func validateHouseholdID(householdID string) error {
	if strings.TrimSpace(householdID) == "" {
		return ValidationError{Message: "household id is required"}
	}
	return nil
}

func validateIDs(householdID, transactionID string) error {
	if err := validateHouseholdID(householdID); err != nil {
		return err
	}
	if strings.TrimSpace(transactionID) == "" {
		return ValidationError{Message: "transaction id is required"}
	}
	return nil
}
