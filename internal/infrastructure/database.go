package infrastructure

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq" // driver Postgres
)

// NewDB cria e retorna a conexão com o Postgres
func NewDB() (*sql.DB, error) {
	// Você pode usar variáveis de ambiente para ler as configs:
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Ajuste de pool, timeouts etc.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
