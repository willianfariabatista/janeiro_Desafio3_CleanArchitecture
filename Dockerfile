# Etapa de build
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copia arquivos de módulo e baixa dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Etapa final
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080 50051
CMD ["./server"]
