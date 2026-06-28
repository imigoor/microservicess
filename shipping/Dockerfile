FROM golang:1.25-alpine AS builder
WORKDIR /app

# 1. Copia arquivos de dependência
COPY microservices/shipping/go.mod microservices/shipping/go.sum ./

# 2. Copia os protos (caminho absoluto no container)
COPY microservices-proto /microservices-proto

# 3. Baixa dependências e garante que o go.mod esteja limpo
RUN go mod download
RUN go mod tidy

# 4. Copia o código fonte
COPY microservices/shipping .

# 5. Compila ignorando qualquer pasta 'vendor' que possa ter sido criada
RUN go build -mod=mod -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3002
CMD ["./main"]