package service

import (
	"context"
	"database/sql"

	"github.com/SEU_USUARIO/my-challenge/internal/domain"
)

// OrderService concentra as operações de negócio relacionadas a Orders
type OrderService struct {
	db *sql.DB
}

// NewOrderService retorna uma instância de OrderService
func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{db: db}
}

// ListOrders retorna todas as orders do banco
func (s *OrderService) ListOrders(ctx context.Context) ([]domain.Order, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, price, quantity, total, created_at, updated_at 
		FROM orders
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order

	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID,
			&o.Name,
			&o.Price,
			&o.Quantity,
			&o.Total,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// CreateOrder cria uma nova order no banco
// name: nome do produto/serviço
// price: preço unitário
// quantity: quantidade solicitada
func (s *OrderService) CreateOrder(ctx context.Context, name string, price float64, quantity int) (domain.Order, error) {
	// Exemplo simples: total é calculado aqui
	total := price * float64(quantity)

	var o domain.Order
	query := `
		INSERT INTO orders (name, price, quantity, total)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, price, quantity, total, created_at, updated_at
	`

	// Usamos QueryRowContext para retornar o registro recém-criado
	err := s.db.QueryRowContext(ctx, query, name, price, quantity, total).Scan(
		&o.ID,
		&o.Name,
		&o.Price,
		&o.Quantity,
		&o.Total,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return domain.Order{}, err
	}

	return o, nil
}
