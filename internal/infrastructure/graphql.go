package infrastructure

import (
	"database/sql"
	"log"

	"github.com/graphql-go/graphql"

	"github.com/SEU_USUARIO/my-challenge/internal/domain"
)

// Cria e retorna o schema GraphQL
func NewGraphQLSchema(db *sql.DB) graphql.Schema {
	orderType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.Int},
			"name":     &graphql.Field{Type: graphql.String},
			"price":    &graphql.Field{Type: graphql.Float},
			"quantity": &graphql.Field{Type: graphql.Int},
			"total":    &graphql.Field{Type: graphql.Float},
		},
	})

	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"listOrders": &graphql.Field{
				Type: graphql.NewList(orderType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Buscamos as orders do banco
					rows, err := db.Query("SELECT id, name, price, quantity, total FROM orders")
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var orders []map[string]interface{}

					for rows.Next() {
						var o domain.Order
						if err := rows.Scan(&o.ID, &o.Name, &o.Price, &o.Quantity, &o.Total); err != nil {
							return nil, err
						}

						// Cada registro será convertido em map[string]interface{}
						orders = append(orders, map[string]interface{}{
							"id":       o.ID,
							"name":     o.Name,
							"price":    o.Price,
							"quantity": o.Quantity,
							"total":    o.Total,
						})
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
