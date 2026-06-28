# Microsserviços com gRPC

Projeto de microsserviços implementado em Go com arquitetura hexagonal e comunicação via gRPC.

## Microsserviços

- **Order** (porta 3000) — recebe pedidos, valida estoque, aciona Payment e Shipping
- **Payment** (porta 3001) — processa cobranças
- **Shipping** (porta 3002) — calcula prazo de entrega

## Pré-requisitos

- Docker e Docker Compose
- Go 1.21+
- grpcurl (para testes)

## Executando com Docker Compose

```bash
docker-compose up --build
```

## Executando manualmente (Windows)

**1. MySQL:**
```bat
docker run -p 3306:3306 -e MYSQL_ROOT_PASSWORD=minhasenha --name mysql-grpc -d mysql
docker exec -it mysql-grpc mysql -u root -pminhasenha -e "CREATE DATABASE IF NOT EXISTS orders; CREATE DATABASE IF NOT EXISTS payment; CREATE DATABASE IF NOT EXISTS shipping;"
```

**2. Payment (CMD 1):**
```bat
cd payment
set DATA_SOURCE_URL=root:minhasenha@tcp(127.0.0.1:3306)/payment
set APPLICATION_PORT=3001
set ENV=development
go run cmd/main.go
```

**3. Shipping (CMD 2):**
```bat
cd shipping
set DATA_SOURCE_URL=root:minhasenha@tcp(127.0.0.1:3306)/shipping
set APPLICATION_PORT=3002
set ENV=development
go run cmd/main.go
```

**4. Order (CMD 3):**
```bat
cd order
set DATA_SOURCE_URL=root:minhasenha@tcp(127.0.0.1:3306)/orders
set APPLICATION_PORT=3000
set ENV=development
set PAYMENT_SERVICE_URL=localhost:3001
set SHIPPING_SERVICE_URL=localhost:3002
go run cmd/main.go
```

**5. Cadastrar produtos no estoque:**
```bat
docker exec -it mysql-grpc mysql -u root -pminhasenha -e "INSERT INTO orders.products (product_code, name, stock, created_at, updated_at) VALUES ('prod1', 'Produto 1', 100, NOW(), NOW()), ('prod2', 'Produto 2', 50, NOW(), NOW());"
```

## Testando com grpcurl

**Pedido normal:**
```bash
grpcurl -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"prod1\", \"unit_price\": 10.0, \"quantity\": 2}]}" -plaintext localhost:3000 Order/Create
```

**Produto inexistente:**
```bash
grpcurl -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"inexistente\", \"unit_price\": 10.0, \"quantity\": 2}]}" -plaintext localhost:3000 Order/Create
```

**Valor maior que 1000:**
```bash
grpcurl -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"prod1\", \"unit_price\": 200.0, \"quantity\": 6}]}" -plaintext localhost:3000 Order/Create
```

**Mais de 50 itens:**
```bash
grpcurl -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"prod1\", \"unit_price\": 10.0, \"quantity\": 51}]}" -plaintext localhost:3000 Order/Create
```

## Prazo de entrega

O prazo mínimo é 1 dia. A cada 5 unidades (somando todos os itens) é adicionado 1 dia.

Exemplos:
- 2 unidades → 1 dia
- 5 unidades → 2 dias
- 10 unidades → 3 dias
