package domain

import "time"

// Order representa nosso modelo de domínio
type Order struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Price     float64   `db:"price"`
	Quantity  int       `db:"quantity"`
	Total     float64   `db:"total"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
