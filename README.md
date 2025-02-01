# Desafio 3 - Clean Architecture

Projeto em Go que expõe:
- API REST para criar e listar orders (`GET/POST /order`)
- API GraphQL para listar orders
- Serviço gRPC para listar orders

## Como rodar

1. Certifique-se de ter [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados.
2. Clone este repositório:
   git clone https://github.com/SEU_USUARIO/my-challenge.git
   cd my-challenge

3. Construa e suba os containers:
docker compose up --build

4. A aplicação subirá na porta8080(HTTP) e `550051(gRP5432.)

Criar ordem :

{
  "name": "Produto X",
  "price": 20.5,
  "quantity": 3
}

Listar ordens:

GET http://localhost:8080/order


Testando Graphql

Listar ordens:
POST http://localhost:8080/graphqlcom

corpo: 
{
  "query": "{ listOrders { id name price quantity total } }"
}


Testando gRPC:
grpcurl -plaintext -proto proto/orders.proto localhost:50051 orders.OrderService/ListOrders

{
  "orders": [
    {
      "id": 1,
      "name": "Produto X",
      "price": 20.5,
      "quantity": 3,
      "total": 61.5
    }
  ]
}


---

## Observações Finais

- Sinta-se à vontade para alterar nomes de variáveis, portas etc. conforme o seu gosto.
- Se quiser rodar migrações de forma diferente (por exemplo, usando um **container** separado só para `migrate`), basta remover do `main.go` e criar esse outro container no `docker-compose.yaml`.
- Se quiser usar [Goose](https://github.com/pressly/goose) ou qualquer outro gerenciador de migrações, basta adaptar.
- Para compilar o proto e gerar o código Go do gRPC, você provavelmente já fez algo como:
  ```bash
  protoc --go_out=. --go-grpc_out=. proto/orders.proto







