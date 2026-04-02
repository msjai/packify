# Packify

Order packs calculator for optimal shipment.

## Tech Stack

- **Backend:** Go 1.26, chi router, slog logger
- **Database:** SQLite (mattn/go-sqlite3)
- **Frontend:** Vanilla JS, Pico CSS (classless)
- **Infrastructure:** Docker, multi-stage build, Render

## Quick Start

### Run in Docker

```bash
make docker-run
```

### Run locally

```bash
make local-run
```

The application will be available at `http://localhost:8080`.

## API Endpoints

### GET /api/packs

Returns current pack sizes.

```bash
curl localhost:8080/api/packs
```

```json
{"packs": [250, 500, 1000, 2000, 5000]}
```

### PUT /api/packs

Updates pack sizes. Validates: non-empty list, all values > 0, no duplicates.

```bash
curl -X PUT localhost:8080/api/packs \
  -H "Content-Type: application/json" \
  -d '{"packs": [23, 31, 53]}'
```

### POST /api/calculate

Calculates optimal pack combination for an order.

```bash
curl -X POST localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"order": 501}'
```

```json
{"packs": [{"pack": 500, "quantity": 1}, {"pack": 250, "quantity": 1}]}
```

## Testing

Run all tests:

```bash
go test ./... -v
```

Run benchmark (750,000 items order with prime pack sizes [17, 29, 47], GCD=1):

```bash
go test ./internal/service/ -bench=. -benchmem
```

## Project Structure

```
packify/
├── cmd/                        # Application entry point
│   └── main.go
├── config/                     # Configuration files
│   └── config.yaml
├── internal/
│   ├── app/                    # Application bootstrap and wiring
│   │   └── app.go
│   ├── config/                 # Config loading (cleanenv + YAML)
│   │   └── config.go
│   ├── handler/                # HTTP handlers and routing
│   │   ├── calculate.go
│   │   ├── calculate_test.go
│   │   ├── get.go
│   │   ├── get_test.go
│   │   ├── helpers.go
│   │   ├── mock_test.go
│   │   ├── router.go
│   │   ├── submit.go
│   │   └── submit_test.go
│   ├── httpserver/             # HTTP server with graceful shutdown
│   │   └── server.go
│   ├── logger/                 # slog setup with pretty handler
│   │   ├── logger.go
│   │   └── slogpretty.go
│   ├── service/                # Business logic (DP algorithm)
│   │   ├── calculate.go
│   │   ├── calculate_test.go
│   │   ├── get.go
│   │   ├── get_test.go
│   │   ├── mock_test.go
│   │   ├── submit.go
│   │   └── submit_test.go
│   └── storage/sqlite/         # SQLite persistence
│       └── sqlite.go
├── web/                        # Embedded frontend (go:embed)
│   ├── app.js
│   ├── embed.go
│   └── index.html
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Live Demo

- Application: https://packify.onrender.com
- Uptime status: https://stats.uptimerobot.com/hGUF7HKXwr/802744656

## Algorithm

The problem is essentially a two-level optimization with strict priority ordering:

**Priority 1:** Find the minimum total number of items ≥ order that can be assembled from whole packs. Packs cannot be broken open, so we may need to send slightly more items than ordered.

**Priority 2:** For that same total, use the fewest packs possible.

This is a variant of the classic coin change problem, but instead of hitting an exact target, we need to reach at least the target.

### Dynamic Programming approach

We build an array `dp[t]` representing the minimum number of packs needed to produce exactly `t` items. We fill it from 0 up to an upper bound, then scan forward from the order value to find the first reachable `t`. That single scan resolves both priorities at once: the first reachable `t` guarantees minimum items (Priority 1), and `dp[t]` already holds the minimum pack count for that total (Priority 2).

### GCD optimization

If all pack sizes share a common divisor `g`, we divide every size and the target by `g`, dramatically shrinking the DP array. For the default sizes (250, 500, 1000, 2000, 5000), GCD = 250, so the normalized sizes become [1, 2, 4, 8, 20] and an order of 1,000,000 reduces to a target of just 4,000. For arbitrary sizes like [23, 31, 53] where GCD = 1, no reduction occurs, but the algorithm still runs efficiently.

### Upper bound

The search range extends to `normalizedTarget + max(normalizedPacks)`. If the exact target is unreachable, the answer is guaranteed to lie within one maximum pack size above it.

### Solution recovery

Alongside `dp[t]` we maintain a `used[t]` array recording which pack size was added last to reach `t`. Once we find the optimal `t`, we backtrack through `used[]` to reconstruct the exact combination, then denormalize by multiplying back by `g`.

### Steps

The algorithm first sanitizes input by removing duplicate and non-positive pack sizes. If order is zero, it returns an empty result. Otherwise it computes the GCD of all pack sizes, normalizes sizes and the target using `ceil(order / g)`, fills the DP table up to the upper bound, scans forward from the normalized target to find the first reachable value, backtracks to recover the pack combination, and returns the result as a map of `{packSize: quantity}`.

### Complexity

- **Time:** O(T × N), where T = normalizedTarget + maxNormPack and N = number of pack sizes. For the default sizes with order = 1,000,000: T ≈ 4,020, N = 5, yielding roughly 20K operations.
- **Space:** O(T) for the `dp` and `used` arrays.