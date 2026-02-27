# Go React SSR Engine (reactgo)

## Run

```bash
# Install deps
npm install
go mod tidy

# Start dev
make dev

# In another terminal, test the routes
curl http://localhost:3000/
curl http://localhost:3000/about
curl http://localhost:3000/posts/42
```

## Run Tests

```bash
# Unit tests (skip engine tests if v8go build is slow)
go test -race -count=1 ./internal/config/ ./internal/cache/ ./internal/router/

# All tests including V8
go test -race -count=1 ./...

# Benchmarks
go test -bench=. -benchmem ./internal/cache/ ./internal/router/
```
