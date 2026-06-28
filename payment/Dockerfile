FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copia os arquivos de dependência
COPY microservices/payment/go.mod microservices/payment/go.sum ./

# Copia os protos (caminho absoluto no container)
COPY microservices-proto /microservices-proto

# Copia o código fonte
COPY microservices/payment .

# Compila
RUN go build -mod=mod -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3001
CMD ["./main"]