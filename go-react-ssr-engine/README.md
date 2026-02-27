# Go React SSR Engine (reactgo)

## Run Tests

```bash
# Unit tests (skip engine tests if v8go build is slow)
go test -race -count=1 ./internal/config/ ./internal/cache/ ./internal/router/

# All tests including V8
go test -race -count=1 ./...

# Benchmarks
go test -bench=. -benchmem ./internal/cache/ ./internal/router/
```
