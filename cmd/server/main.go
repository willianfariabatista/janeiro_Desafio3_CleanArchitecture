package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/graphql-go/handler"
	// Caso queira carregar variáveis do .env, descomente a linha abaixo e adicione o pacote "github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	// Atualize os imports abaixo de acordo com o seu module path
	"github.com/willianfariabatista/my-challenge/internal/infrastructure"
	"github.com/willianfariabatista/my-challenge/internal/service"
)

func main() {
	// Cria a conexão com o banco de dados
	db, err := infrastructure.NewDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	// (Opcional) Aqui você pode executar suas migrações automaticamente se desejar

	// Inicia o servidor gRPC em uma goroutine separada
	go infrastructure.StartGRPCServer(db, "50051")

	// Cria a camada de serviço para Orders
	orderService := service.NewOrderService(db)

	// Configura as rotas HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListOrders(w, r, orderService)
		case http.MethodPost:
			handleCreateOrder(w, r, orderService)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	// Configura o endpoint GraphQL
	schema := infrastructure.NewGraphQLSchema(db)
	graphqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true, // habilita a interface GraphiQL no navegador
	})
	mux.Handle("/graphql", graphqlHandler)

	// Inicia o servidor HTTP
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Servidor HTTP iniciado na porta %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Falha ao iniciar servidor HTTP: %v", err)
	}
}

type createOrderRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func handleListOrders(w http.ResponseWriter, r *http.Request, s *service.OrderService) {
	ctx := r.Context()
	orders, err := s.ListOrders(ctx)
	if err != nil {
		http.Error(w, "Erro ao listar orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func handleCreateOrder(w http.ResponseWriter, r *http.Request, s *service.OrderService) {
	ctx := r.Context()
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	order, err := s.CreateOrder(ctx, req.Name, req.Price, req.Quantity)
	if err != nil {
		http.Error(w, "Erro ao criar order: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
