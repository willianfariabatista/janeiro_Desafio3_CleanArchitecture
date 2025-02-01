package domain

import "time"

// Order representa o modelo de domínio para um pedido.
type Order struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Price     float64   `db:"price" json:"price"`
	Quantity  int       `db:"quantity" json:"quantity"`
	Total     float64   `db:"total" json:"total"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
