# janeiro_Desafio3_CleanArchitecture

Desafio de Janeiro 2025 Faculdade Full Cycle - Go Expert - CleanArchitecture

# Desafio 3 - Clean Architecture

Este projeto implementa um sistema de cadastro e listagem de pedidos (orders), utilizando Go, PostgreSQL, REST, GraphQL e gRPC, seguindo conceitos de Clean Architecture.

## Funcionalidades

1. **CRUD de Pedidos:**

2.1 - **POST /order – Cria um pedido.**

2.2 - **GET /order – Lista todos os pedidos.**

4. **GraphQL Endpoint:**
Query listOrders para listar pedidos.

5. **gRPC Endpoint:**
Método ListOrders para recuperar pedidos via Protobuf.

## Requisitos

Go 1.21+ (ou versão compatível com as bibliotecas utilizadas) se for executar localmente sem Docker.
Docker e Docker Compose (opcional, se for usar containers para rodar tudo).

## Estrutura de Pastas

```
janeiro_Desafio3_CleanArchitecture/
├── cmd/
│   └── server/
│       └── main.go            # Arquivo principal que inicia HTTP e gRPC
├── db/
│   └── migrations/
│       └── 001_create_orders.sql
├── internal/
│   ├── domain/
│   │   └── order.go
│   ├── infrastructure/
│   │   ├── database.go
│   │   ├── graphql.go
│   │   └── grpc_server.go
│   └── service/
│       ├── order_service.go
│       └── orderspb/          # Gerações do protoc (orders.pb.go, orders_grpc.pb.go)
├── proto/
│   └── orders.proto
├── api.http                   # Exemplos de requisições HTTP
├── docker-compose.yaml
├── Dockerfile
├── teste.env                  # Exemplo de variáveis de ambiente
└── README.md                  # Este arquivo
```

## Configuração de Ambiente

### 1. **Variáveis de Ambiente**

O arquivo teste.env (ou .env) deve conter:
```
HTTP_PORT=8080
GRPC_PORT=50051
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=orders_db
```
**DESCRIÇÃO CAMPOS**
```
HTTP_PORT: Porta para o servidor HTTP e GraphQL.
GRPC_PORT: Porta para o servidor gRPC.
DB_HOST: Host do banco de dados (no Docker Compose, utilizamos db).
DB_PORT: Porta do banco de dados.
DB_USER: Usuário do PostgreSQL.
DB_PASSWORD: Senha do PostgreSQL.
DB_NAME: Nome do banco de dados.
```
### 2. **Subindo o Projeto**

2.1 **Editar o arquivo docker-compose.yaml para incluir os serviços app (nossa aplicação) e db (Postgres). Por exemplo:**
```
version: '3.9'
services:
  db:
    image: postgres:13
    container_name: db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: orders_db
    ports:
      - "5432:5432"

  app:
    image: my-challenge:latest
    container_name: my-app
    env_file:
      - teste.env
    ports:
      - "8080:8080"
      - "50051:50051"
    depends_on:
      - db
```
2.2 **Buildar a imagem da aplicação (caso ainda não tenha feito):**
```
docker build -t my-challenge:latest .
```
2.3 **Subir os serviços:**
```
docker compose up -d
```
2.4 **Verifique se ambos estão rodando:**
```
docker ps
```
Se tudo estiver correto, você terá my-app (expondo 8080 e 50051) e db (expondo 5432).

### 3. **Opção B: Docker Manualmente**

3.1 **Criar uma rede:**
```
docker network create meu-network
```
3.2 **Rodar o Postgres:**
```
docker run -d \
  --name db \
  --network meu-network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=orders_db \
  -p 5432:5432 \
  postgres:13
```
3.3 **Buildar a imagem da aplicação:**
docker build -t my-challenge:latest .

3.4 **Rodar a aplicação na mesma rede:**
```
docker run -d \
  --name my-app \
  --network meu-network \
  -p 8080:8080 \
  -p 50051:50051 \
  --env-file teste.env \
  my-challenge:latest
```
### 4. **Opção C: Executando Localmente Sem Docker**

4.1 **Instale e configure o PostgreSQL localmente (ou rode em Docker somente o banco).**

4.2 **Exporte as variáveis de ambiente ou crie um .env e chame godotenv.Load() no main.go.**

4.3 **Execute:**
```
go run cmd/server/main.go
```
4.4 **O servidor iniciará, escutando em :8080 para HTTP/GraphQL e :50051 para gRPC, conforme seu código.**

## Migrações (Banco de Dados)

Se as migrações não forem executadas automaticamente, rode manualmente o script no container do Postgres ou localmente:
```
docker exec -it db bash
psql -U postgres -d orders_db -f /caminho/para/db/migrations/001_create_orders.sql
```
Verifique também se há bibliotecas ou funções no main.go para executar migrações em runtime.

## Testes dos Endpoints

### 1. REST

**Criar Pedido:**

Endpoint: POST http://localhost:8080/order
```
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--data '{
  "name": "Produto Exemplo",
  "price": 100.50,
  "quantity": 2
}'
```
### Listar Pedidos:
```
Endpoint: GET http://localhost:8080/order
```
Você pode usar o arquivo api.http (com a extensão REST Client no VSCode) ou ferramentas como Postman.

### 2. GraphQL

### Abra seu navegador em:
```
http://localhost:8080/graphql
```
### Para usar a interface GraphiQL e rode uma query como:
```
{
  listOrders {
    id
    name
    price
    quantity
    total
    created_at
    updated_at
  }
}
```
### 3. gRPC

### Utilize o grpcurl ou outro cliente gRPC:
```
grpcurl -plaintext -d '{}' localhost:50051 orders.OrderService/ListOrders
```
Isso deve retornar a lista de pedidos.

### Possíveis Erros e Soluções

1. Porta em uso: Se a porta 8080 ou 50051 estiver em uso, encerre o processo ocupando a porta ou mapeie outra porta local (-p 8081:8080, por exemplo).
2. "dial tcp: lookup db: no such host": Significa que o container não encontra o host db. Use Docker Compose ou rodar ambos containers na mesma rede.
3. "Database is uninitialized and superuser password is not specified": Defina POSTGRES_PASSWORD para o usuário postgres no seu container do PostgreSQL.

### Contribuição

Faça um fork do repositório.
Crie um branch para sua feature/fix: git checkout -b my-feature.
Faça o commit: git commit -m 'Adiciona minha feature'.
Faça o push para o branch: git push origin my-feature.
Crie um novo Pull Request.




