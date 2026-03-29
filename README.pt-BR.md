🇧🇷 Português | [🇺🇸 English](./README.md)

# ⚡ pix-simulator

> Sistema de pagamentos Pix de alta performance construído em Go — capaz de processar **5.000+ transações por segundo** com zero falhas.

---

## 📋 Índice

- [Visão Geral](#-visão-geral)
- [Arquitetura](#-arquitetura)
- [Fluxo de uma Transação](#-fluxo-de-uma-transação)
- [Tecnologias](#-tecnologias)
- [Estrutura do Projeto](#-estrutura-do-projeto)
- [Pré-requisitos](#-pré-requisitos)
- [Como Rodar](#-como-rodar)
- [Teste de Carga](#-teste-de-carga)
- [Performance](#-performance)
- [API](#-api)

---

## 🎯 Visão Geral

O **pix-simulator** é um sistema de pagamentos assíncrono que simula o processamento de transações Pix em alta escala. A arquitetura é projetada para maximizar throughput e garantir consistência financeira — exatamente como sistemas reais de big tech.

**Destaques:**
- **5.329 req/s** de throughput sustentado
- **0% de taxa de erro** sob carga intensa
- **p(99) de 459ms** — 99% das requisições respondidas abaixo de 500ms
- **1.28 milhão de transações** processadas em 4 minutos
- **Idempotência** garantida — sem duplicatas, mesmo sob retry
- **Consistência ACID** nas operações financeiras

---

## 🏗 Arquitetura

```
┌─────────────────────────────────────────────────────────────────┐
│                        pix-simulator                            │
│                                                                 │
│   ┌──────────┐    ┌─────────────────┐    ┌──────────────────┐  │
│   │  Cliente │───▶│    API (Go)     │───▶│  Redis           │  │
│   │  (k6)    │    │  :8080          │    │  • Fila (List)   │  │
│   └──────────┘    │                 │◀───│  • Idempotência  │  │
│                   │  • Valida req   │    └────────┬─────────┘  │
│                   │  • Dedup key    │             │            │
│                   │  • Enfileira    │             │ consume    │
│                   └─────────────────┘             ▼            │
│                                        ┌──────────────────┐   │
│                                        │  Worker (Go)     │   │
│                                        │  200 goroutines  │   │
│                                        │                  │   │
│                                        │  • INSERT tx     │   │
│                                        │  • UPDATE sender │   │
│                                        │  • UPDATE recv   │   │
│                                        │  • SET done      │   │
│                                        └────────┬─────────┘   │
│                                                 │             │
│                                                 ▼             │
│                                        ┌──────────────────┐   │
│                                        │  PostgreSQL      │   │
│                                        │  • accounts      │   │
│                                        │  • transactions  │   │
│                                        └──────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Fluxo de uma Transação

```
1. Cliente envia POST /pix
         │
         ▼
2. API verifica Idempotency-Key no Redis
   ├── duplicada? → retorna "duplicated" (200)
   └── nova? → continua
         │
         ▼
3. API gera txID único (UUID)
         │
         ▼
4. API empurra Job na fila Redis
   { id, sender, receiver, amount }
         │
         ▼
5. API salva Idempotency-Key no Redis (TTL: 10min)
         │
         ▼
6. API retorna "queued" instantaneamente ⚡
         │
         ▼ (assíncrono)
7. Worker consome Job da fila
         │
         ▼
8. Worker abre transação no PostgreSQL
   ├── INSERT transactions (status: pending)
   ├── UPDATE accounts SET balance - amount (sender)
   ├── UPDATE accounts SET balance + amount (receiver)
   └── UPDATE transactions (status: done)
         │
         ▼
9. Transação commitada com ACID ✅
```

---

## 🛠 Tecnologias

| Tecnologia | Versão | Uso |
|---|---|---|
| **Go** | 1.22 | API e Worker |
| **PostgreSQL** | 15 | Persistência e consistência financeira |
| **Redis** | 7 | Fila de mensagens e idempotência |
| **Docker** | — | Containerização |
| **Docker Compose** | — | Orquestração local |
| **k6** | — | Teste de carga |

---

## 📁 Estrutura do Projeto

```
pix-simulator/
├── cmd/
│   ├── api/
│   │   └── main.go          # Servidor HTTP
│   └── worker/
│       └── main.go          # Consumidor da fila
│
├── internal/
│   ├── database/
│   │   └── postgres.go      # Conexão e pool do PostgreSQL
│   ├── idempotency/
│   │   └── redis.go         # Verificação de chaves duplicadas
│   ├── queue/
│   │   └── redis.go         # Push/Pop na fila Redis
│   └── service/
│       └── payment.go       # Lógica de processamento do Pix
│
├── docker/
│   ├── Dockerfile.api       # Build da API
│   └── Dockerfile.worker    # Build do Worker
│
├── migrations/
│   ├── 001_create_tables.sql
│   └── 002_seed_accounts.sql
│
├── load-test/
│   └── test.js              # Script k6
│
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## ✅ Pré-requisitos

- [Docker](https://www.docker.com/) >= 24
- [Docker Compose](https://docs.docker.com/compose/) >= 2
- [k6](https://k6.io/docs/get-started/installation/) (para testes de carga)

---

## 🚀 Como Rodar

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/pix-simulator.git
cd pix-simulator
```

### 2. Suba os serviços

```bash
docker compose up --build
```

Isso irá subir:
- **PostgreSQL** na porta `5432`
- **Redis** na porta `6379`
- **Migrate** — aplica as migrations automaticamente
- **API** na porta `8080`
- **Worker** — 200 goroutines processando a fila

### 3. Verifique se está rodando

```bash
curl -X POST http://localhost:8080/pix \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: test-001" \
  -d '{"sender": "1", "receiver": "2", "amount": 10}'

# Resposta esperada: queued
```

### 4. Parar os serviços

```bash
docker compose down
```

---

## 🔍 Inspecionando os Dados

### PostgreSQL (via terminal)

```bash
docker exec -it docker-postgres-1 psql -U pix -d pix
```

```sql
-- Ver transações
SELECT * FROM transactions ORDER BY created_at DESC LIMIT 10;

-- Ver saldos
SELECT * FROM accounts;

-- Contar transações por status
SELECT status, COUNT(*) FROM transactions GROUP BY status;
```

### Redis (via terminal)

```bash
docker exec -it docker-redis-1 redis-cli
```

```bash
KEYS *          # todas as chaves
DBSIZE          # total de chaves
LLEN pix        # jobs restantes na fila
```

### Redis (via interface gráfica)

Recomendado: **[RedisInsight](https://redis.io/insight)** (gratuito e oficial)
- Host: `localhost`
- Porta: `6379`

---

## 🧪 Teste de Carga

### Preparar o ambiente

Antes de rodar o teste, limpe os dados anteriores:

```bash
# Limpar Redis
docker exec docker-redis-1 redis-cli FLUSHALL

# Limpar transações e resetar saldos
docker exec docker-postgres-1 psql -U pix -d pix -c "TRUNCATE TABLE transactions;"
docker exec docker-postgres-1 psql -U pix -d pix -c "
  UPDATE accounts SET balance = 10000 WHERE id = '1';
  UPDATE accounts SET balance = 10000 WHERE id = '2';
"
```

Ou use o script pronto:

```bash
chmod +x run-test.sh
./run-test.sh
```

### Rodar o teste

```bash
k6 run load-test/test.js
```

**Cenários do teste:**

| Estágio | Duração | VUs | Descrição |
|---|---|---|---|
| Aquecimento | 30s | 200 | Warm-up gradual |
| Carga normal | 1m | 1000 | Operação típica |
| Pico | 2m | 2000 | Big tech load |
| Cooldown | 30s | 0 | Encerramento |

> **Nota:** Após o k6 finalizar, o worker pode continuar processando a fila por alguns segundos. Isso é esperado — o sistema é assíncrono. Verifique com `docker exec docker-redis-1 redis-cli LLEN pix` para ver os jobs restantes.

---

## 📊 Performance

Resultado real em ambiente local (Docker):

```
✓ 0%      taxa de erro
✓ 459ms   p(99) de latência
✓ 5.329   requisições por segundo
✓ 1.28M   transações em 4 minutos
✓ 104ms   mediana de resposta
✓ 100%    checks passando
```

| Métrica | pix-simulator | Pix real (BACEN) |
|---|---|---|
| Throughput | 5.329 req/s | ~1.000 req/s |
| p(99) | 459ms | — |
| Taxa de erro | 0% | — |

> O pix-simulator processa **~5x mais** que o sistema real do Pix em ambiente de desenvolvimento.

---

## 🌐 API

### `POST /pix`

Enfileira uma transação Pix para processamento.

**Headers**

| Header | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `Content-Type` | string | ✅ | `application/json` |
| `Idempotency-Key` | string | ✅ | Chave única para evitar duplicatas |

**Body**

```json
{
  "sender": "1",
  "receiver": "2",
  "amount": 10.00
}
```

**Respostas**

| Status | Body | Descrição |
|---|---|---|
| `200` | `queued` | Transação enfileirada com sucesso |
| `200` | `duplicated` | Idempotency-Key já utilizada |
| `400` | `invalid body` | Corpo da requisição inválido |
| `500` | `queue error` | Erro ao enfileirar no Redis |

**Exemplo**

```bash
curl -X POST http://localhost:8080/pix \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(uuidgen)" \
  -d '{
    "sender": "1",
    "receiver": "2",
    "amount": 50.00
  }'
```

---

## 📄 Licença

MIT © pix-simulator
