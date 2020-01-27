package account

import "context"

type Account struct {
	ID       string  `json:"id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type AccountRepository interface {
	GetUser(ctx context.Context, userID string) (*Account, error)
	GetAllUsers(ctx context.Context) ([]*Account, error)
}
