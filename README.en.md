[🇧🇷 Português](./README.md) | 🇺🇸 English

# ⚡ pix-simulator

> High-performance Pix payment system built in Go — capable of processing **5,000+ transactions per second** with zero errors.

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Architecture](#-architecture)
- [Transaction Flow](#-transaction-flow)
- [Technologies](#-technologies)
- [Project Structure](#-project-structure)
- [Prerequisites](#-prerequisites)
- [How to Run](#-how-to-run)
- [Load Testing](#-load-testing)
- [Performance](#-performance)
- [API](#-api)

---

## 🎯 Overview

**pix-simulator** is an asynchronous payment system that simulates Pix transaction processing at high scale. The architecture is designed to maximize throughput and ensure financial consistency — exactly like real big tech systems.

**Highlights:**
- **5,329 req/s** sustained throughput
- **0% error rate** under intense load
- **p(99) of 459ms** — 99% of requests responded to under 500ms
- **1.28 million transactions** processed in 4 minutes
- **Idempotency** guaranteed — no duplicates, even under retry
- **ACID consistency** on financial operations

---

## 🏗 Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        pix-simulator                            │
│                                                                 │
│   ┌──────────┐    ┌─────────────────┐    ┌──────────────────┐  │
│   │  Client  │───▶│    API (Go)     │───▶│  Redis           │  │
│   │  (k6)    │    │  :8080          │    │  • Queue (List)  │  │
│   └──────────┘    │                 │◀───│  • Idempotency   │  │
│                   │  • Validate req │    └────────┬─────────┘  │
│                   │  • Dedup key    │             │            │
│                   │  • Enqueue      │             │ consume    │
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

## 🔄 Transaction Flow

```
1. Client sends POST /pix
         │
         ▼
2. API checks Idempotency-Key in Redis
   ├── duplicate? → returns "duplicated" (200)
   └── new? → continues
         │
         ▼
3. API generates unique txID (UUID)
         │
         ▼
4. API pushes Job to Redis queue
   { id, sender, receiver, amount }
         │
         ▼
5. API saves Idempotency-Key in Redis (TTL: 10min)
         │
         ▼
6. API returns "queued" instantly ⚡
         │
         ▼ (async)
7. Worker consumes Job from queue
         │
         ▼
8. Worker opens PostgreSQL transaction
   ├── INSERT transactions (status: pending)
   ├── UPDATE accounts SET balance - amount (sender)
   ├── UPDATE accounts SET balance + amount (receiver)
   └── UPDATE transactions (status: done)
         │
         ▼
9. Transaction committed with ACID ✅
```

---

## 🛠 Technologies

| Technology | Version | Usage |
|---|---|---|
| **Go** | 1.22 | API and Worker |
| **PostgreSQL** | 15 | Persistence and financial consistency |
| **Redis** | 7 | Message queue and idempotency |
| **Docker** | — | Containerization |
| **Docker Compose** | — | Local orchestration |
| **k6** | — | Load testing |

---

## 📁 Project Structure

```
pix-simulator/
├── cmd/
│   ├── api/
│   │   └── main.go          # HTTP server
│   └── worker/
│       └── main.go          # Queue consumer
│
├── internal/
│   ├── database/
│   │   └── postgres.go      # PostgreSQL connection and pool
│   ├── idempotency/
│   │   └── redis.go         # Duplicate key verification
│   ├── queue/
│   │   └── redis.go         # Push/Pop on Redis queue
│   └── service/
│       └── payment.go       # Pix processing logic
│
├── docker/
│   ├── Dockerfile.api       # API build
│   └── Dockerfile.worker    # Worker build
│
├── migrations/
│   ├── 001_create_tables.sql
│   └── 002_seed_accounts.sql
│
├── load-test/
│   └── test.js              # k6 script
│
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## ✅ Prerequisites

- [Docker](https://www.docker.com/) >= 24
- [Docker Compose](https://docs.docker.com/compose/) >= 2
- [k6](https://k6.io/docs/get-started/installation/) (for load testing)

---

## 🚀 How to Run

### 1. Clone the repository

```bash
git clone https://github.com/your-username/pix-simulator.git
cd pix-simulator
```

### 2. Start the services

```bash
docker compose up --build
```

This will start:
- **PostgreSQL** on port `5432`
- **Redis** on port `6379`
- **Migrate** — applies migrations automatically
- **API** on port `8080`
- **Worker** — 200 goroutines processing the queue

### 3. Verify it's running

```bash
curl -X POST http://localhost:8080/pix \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: test-001" \
  -d '{"sender": "1", "receiver": "2", "amount": 10}'

# Expected response: queued
```

### 4. Stop the services

```bash
docker compose down
```

---

## 🔍 Inspecting Data

### PostgreSQL (via terminal)

```bash
docker exec -it docker-postgres-1 psql -U pix -d pix
```

```sql
-- View transactions
SELECT * FROM transactions ORDER BY created_at DESC LIMIT 10;

-- View balances
SELECT * FROM accounts;

-- Count transactions by status
SELECT status, COUNT(*) FROM transactions GROUP BY status;
```

### Redis (via terminal)

```bash
docker exec -it docker-redis-1 redis-cli
```

```bash
KEYS *          # all keys
DBSIZE          # total keys
LLEN pix        # jobs remaining in queue
```

### Redis (via GUI)

Recommended: **[RedisInsight](https://redis.io/insight)** (free and official)
- Host: `localhost`
- Port: `6379`

---

## 🧪 Load Testing

### Prepare the environment

Before running the test, clean up previous data:

```bash
# Clean Redis
docker exec docker-redis-1 redis-cli FLUSHALL

# Clean transactions and reset balances
docker exec docker-postgres-1 psql -U pix -d pix -c "TRUNCATE TABLE transactions;"
docker exec docker-postgres-1 psql -U pix -d pix -c "
  UPDATE accounts SET balance = 10000 WHERE id = '1';
  UPDATE accounts SET balance = 10000 WHERE id = '2';
"
```

Or use the ready-made script:

```bash
chmod +x run-test.sh
./run-test.sh
```

### Run the test

```bash
k6 run load-test/test.js
```

**Test stages:**

| Stage | Duration | VUs | Description |
|---|---|---|---|
| Warm-up | 30s | 200 | Gradual ramp-up |
| Normal load | 1m | 1000 | Typical operation |
| Peak | 2m | 2000 | Big tech load |
| Cooldown | 30s | 0 | Shutdown |

> **Note:** After k6 finishes, the worker may continue processing the queue for a few seconds. This is expected — the system is async. Check with `docker exec docker-redis-1 redis-cli LLEN pix` to see remaining jobs.

---

## 📊 Performance

Real results in a local environment (Docker):

```
✓ 0%      error rate
✓ 459ms   p(99) latency
✓ 5,329   requests per second
✓ 1.28M   transactions in 4 minutes
✓ 104ms   median response time
✓ 100%    checks passing
```

| Metric | pix-simulator | Real Pix (BACEN) |
|---|---|---|
| Throughput | 5,329 req/s | ~1,000 req/s |
| p(99) | 459ms | — |
| Error rate | 0% | — |

> pix-simulator processes **~5x more** than the real Pix system in a development environment.

---

## 🌐 API

### `POST /pix`

Queues a Pix transaction for processing.

**Headers**

| Header | Type | Required | Description |
|---|---|---|---|
| `Content-Type` | string | ✅ | `application/json` |
| `Idempotency-Key` | string | ✅ | Unique key to prevent duplicates |

**Body**

```json
{
  "sender": "1",
  "receiver": "2",
  "amount": 10.00
}
```

**Responses**

| Status | Body | Description |
|---|---|---|
| `200` | `queued` | Transaction successfully queued |
| `200` | `duplicated` | Idempotency-Key already used |
| `400` | `invalid body` | Invalid request body |
| `500` | `queue error` | Error queuing in Redis |

**Example**

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

## 📄 License

MIT © pix-simulator
