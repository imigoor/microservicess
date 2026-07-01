# Ecossistema de Microsserviços com gRPC e Docker

Projeto final construído em Go. O sistema implementa uma arquitetura hexagonal, comunicação via gRPC, persistência segregada e orquestração resiliente com Docker.

## Cumprimento dos Requisitos

Este repositório atende a todos os critérios exigidos:

1. **Arquitetura Hexagonal (Ports and Adapters):** A lógica de negócio (`internal/application/core`) está estritamente isolada da infraestrutura (`internal/adapters`).
2. **Comunicação gRPC:** Os serviços trocam mensagens utilizando contratos rígidos via Protocol Buffers (`.proto`).
3. **Database per Service:** O sistema garante isolamento de estado. Um único container MySQL gerencia bancos distintos (`orders`, `payment`, `shipping`).
4. **Docker Multi-stage Build:** Os `Dockerfile`s utilizam compilação em múltiplos estágios com Alpine Linux, resultando em imagens seguras e extremamente leves.
5. **Orquestração e Resiliência (Auto-healing):** O `docker-compose.yml` sobe a infraestrutura completa. Os microsserviços possuem políticas de `restart: on-failure` e `healthchecks`, aguardando automaticamente o banco de dados estar pronto sem intervenção manual.

## Microsserviços

- **Order** (porta 3000) — recebe pedidos, valida estoque, aciona Payment e Shipping.
- **Payment** (porta 3001) — processa as transações financeiras.
- **Shipping** (porta 3002) — calcula o prazo de entrega com base no volume.

## Pré-requisitos

- Docker e Docker Compose
- Go 1.21+
- grpcurl (para testes)

## Executando o Projeto (Docker)

Esta é a forma recomendada de avaliar o projeto, pois garante o isolamento total.

**1. Suba a infraestrutura completa:**
```bash
docker-compose up --build
```
*(Nota: Na primeira vez, o MySQL fará sua inicialização. Os microsserviços reiniciarão sozinhos até que o banco esteja pronto).*

**2. Cadastre produtos no banco de dados (Estoque):**
Em um novo terminal, rode o comando abaixo para injetar os dados iniciais:
```bash
docker exec -it mysql-grpc mysql -u root -pminhasenha -e "INSERT INTO orders.products (product_code, name, stock, created_at, updated_at) VALUES ('prod1', 'Produto 1', 100, NOW(), NOW()), ('prod2', 'Produto 2', 50, NOW(), NOW());"
```

## Testando com grpcurl (Cenários de Validação)

Com o ambiente rodando, você pode validar as regras de negócio disparando os comandos abaixo:

**✅ Pedido normal (Sucesso total):**
```bash
grpcurl -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"prod1\", \"unit_price\": 10.0, \"quantity\": 2}]}" -plaintext localhost:3000 Order/Create
```

**❌ Regra 1: Produto inexistente no banco:**
```bash
grpcurl -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"inexistente\", \"unit_price\": 10.0, \"quantity\": 2}]}" -plaintext localhost:3000 Order/Create
```

**❌ Regra 2: Valor do pedido maior que 1000:**
```bash
grpcurl -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"prod1\", \"unit_price\": 200.0, \"quantity\": 6}]}" -plaintext localhost:3000 Order/Create
```

**❌ Regra 3: Pedido com mais de 50 itens:**
```bash
grpcurl -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"prod1\", \"unit_price\": 10.0, \"quantity\": 51}]}" -plaintext localhost:3000 Order/Create
```

## Regra de Negócio: Prazo de Entrega (Shipping)

O prazo de entrega é calculado de forma dinâmica pelo microsserviço de frete, isolado do contexto de pedidos.
- O prazo mínimo padrão é de **1 dia**.
- A cada **5 unidades** somadas no pedido, adiciona-se **1 dia extra** ao frete.

**Exemplos:**
- 2 unidades → 1 dia
- 5 unidades → 2 dias
- 10 unidades → 3 dias
