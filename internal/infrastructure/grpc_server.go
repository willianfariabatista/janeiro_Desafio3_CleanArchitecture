package infrastructure

import (
	"context"
	"database/sql"
	"log"
	"net"

	"github.com/willianfariabatista/my-challenge/internal/domain"
	"github.com/willianfariabatista/my-challenge/internal/service/orderspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server implementa o serviço gRPC para Orders.
type Server struct {
	orderspb.UnimplementedOrderServiceServer
	DB *sql.DB
}

// ListOrders implementa o método ListOrders via gRPC.
func (s *Server) ListOrders(ctx context.Context, req *orderspb.ListOrdersRequest) (*orderspb.ListOrdersResponse, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id, name, price, quantity, total, created_at, updated_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*orderspb.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.Name, &o.Price, &o.Quantity, &o.Total, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &orderspb.Order{
			Id:       int32(o.ID),
			Name:     o.Name,
			Price:    o.Price,
			Quantity: int32(o.Quantity),
			Total:    o.Total,
		})
	}
	return &orderspb.ListOrdersResponse{
		Orders: orders,
	}, nil
}

// StartGRPCServer inicia o servidor gRPC na porta informada.
func StartGRPCServer(db *sql.DB, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Erro ao ouvir na porta %s: %v", port, err)
	}
	s := grpc.NewServer()
	orderspb.RegisterOrderServiceServer(s, &Server{DB: db})
	// Habilita a reflexão para ferramentas como grpcurl
	reflection.Register(s)
	log.Printf("gRPC rodando na porta %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Erro ao servir gRPC: %v", err)
	}
}
