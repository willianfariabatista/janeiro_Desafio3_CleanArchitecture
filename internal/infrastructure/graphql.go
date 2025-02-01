package infrastructure

import (
	"database/sql"
	"log"

	"github.com/graphql-go/graphql"
	"github.com/willianfariabatista/my-challenge/internal/domain"
)

func NewGraphQLSchema(db *sql.DB) graphql.Schema {
	// Define o objeto GraphQL para Order
	orderType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
			"quantity": &graphql.Field{
				Type: graphql.Int,
			},
			"total": &graphql.Field{
				Type: graphql.Float,
			},
			"created_at": &graphql.Field{
				// Para simplificar, formataremos a data como string.
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// Define a query root
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"listOrders": &graphql.Field{
				Type: graphql.NewList(orderType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					rows, err := db.Query("SELECT id, name, price, quantity, total, created_at, updated_at FROM orders")
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var orders []domain.Order
					for rows.Next() {
						var o domain.Order
						if err := rows.Scan(&o.ID, &o.Name, &o.Price, &o.Quantity, &o.Total, &o.CreatedAt, &o.UpdatedAt); err != nil {
							return nil, err
						}
						orders = append(orders, o)
					}
					return orders, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err != nil {
		log.Fatalf("Falha ao criar o schema GraphQL: %v", err)
	}
	return schema
}
