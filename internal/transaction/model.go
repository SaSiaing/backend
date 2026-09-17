package transaction

import "time"

type Type string

const (
	TypeExpense Type = "expense"
	TypeIncome  Type = "income"
)

type Transaction struct {
	ID          string    `json:"id"`
	HouseholdID string    `json:"household_id"`
	PayerID     string    `json:"payer_id"`
	CategoryID  string    `json:"category_id"`
	Type        Type      `json:"type"`
	Amount      int64     `json:"amount"`
	OccurredAt  time.Time `json:"occurred_at"`
	Memo        string    `json:"memo"`
	Version     int64     `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
}

type CreateInput struct {
	PayerID    string    `json:"payer_id"`
	CategoryID string    `json:"category_id"`
	Type       Type      `json:"type"`
	Amount     int64     `json:"amount"`
	OccurredAt time.Time `json:"occurred_at"`
	Memo       string    `json:"memo"`
}

type UpdateInput struct {
	PayerID    *string    `json:"payer_id"`
	CategoryID *string    `json:"category_id"`
	Type       *Type      `json:"type"`
	Amount     *int64     `json:"amount"`
	OccurredAt *time.Time `json:"occurred_at"`
	Memo       *string    `json:"memo"`
	Version    int64      `json:"version"`
}

type ListFilter struct {
	From       *time.Time
	To         *time.Time
	PayerID    string
	CategoryID string
	Type       Type
}
