package infrastructure

import (
	"context"
	"database/sql"
	"log"
	"net"

	"github.com/SEU_USUARIO/my-challenge/internal/domain"
	"github.com/SEU_USUARIO/my-challenge/internal/service/orderspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server implementa nosso serviço gRPC
type Server struct {
	orderspb.UnimplementedOrderServiceServer
	DB *sql.DB
}

// ListOrders implementa o método gRPC
func (s *Server) ListOrders(ctx context.Context, req *orderspb.ListOrdersRequest) (*orderspb.ListOrdersResponse, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id, name, price, quantity, total FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*orderspb.Order

	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.Name, &o.Price, &o.Quantity, &o.Total); err != nil {
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

// StartGRPCServer inicia o servidor gRPC
func StartGRPCServer(db *sql.DB, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Erro ao ouvir na porta %s: %v", port, err)
	}

	s := grpc.NewServer()
	orderspb.RegisterOrderServiceServer(s, &Server{DB: db})

	// Reflection (para usar ferramentas como grpcurl)
	reflection.Register(s)

	log.Printf("gRPC rodando na porta %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Erro ao servir gRPC: %v", err)
	}
}
