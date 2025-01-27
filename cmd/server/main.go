package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/graphql-go/handler"

	"github.com/SEU_USUARIO/my-challenge/internal/infrastructure"
	"github.com/SEU_USUARIO/my-challenge/internal/service"
)

func main() {
	// (Opcional) Carregar variáveis de ambiente, se necessário
	// Exemplo se usar github.com/joho/godotenv:
	// _ = godotenv.Load()

	// Cria a conexão com o banco de dados
	db, err := infrastructure.NewDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	// (Opcional) Executar migrações, por ex:
	// err = runMigrations("db/migrations", db)
	// if err != nil {
	// 	log.Fatalf("Falha ao rodar migrações: %v", err)
	// }

	// Inicia o servidor gRPC em goroutine separada
	go func() {
		infrastructure.StartGRPCServer(db, "50051")
	}()

	// Cria o service de Order
	orderService := service.NewOrderService(db)

	// Configura as rotas do servidor HTTP
	mux := http.NewServeMux()

	// REST: /order (GET -> listar, POST -> criar)
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

	// GraphQL
	schema := infrastructure.NewGraphQLSchema(db)
	graphqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true, // se quiser a interface GraphiQL no navegador
	})
	mux.Handle("/graphql", graphqlHandler)

	// Configura e inicia o servidor HTTP
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080" // fallback padrão
	}
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("HTTP Server iniciado na porta %s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Falha ao iniciar servidor HTTP: %v", err)
	}
}

// handleListOrders lida com GET /order
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

// handleCreateOrder lida com POST /order
func handleCreateOrder(w http.ResponseWriter, r *http.Request, s *service.OrderService) {
	ctx := r.Context()

	type createOrderRequest struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

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
