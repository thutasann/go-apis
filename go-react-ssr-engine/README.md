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

## Health check

```bash
# Health check
curl http://localhost:3000/_health

# Normal page
curl -v http://localhost:3000/

# Gzip
curl -H "Accept-Encoding: gzip" --compressed http://localhost:3000/

# ETag 304
ETAG=$(curl -sI http://localhost:3000/ | grep ETag | awk '{print $2}' | tr -d '\r')
curl -H "If-None-Match: $ETAG" -o /dev/null -w "%{http_code}" http://localhost:3000/

# Skip cache
curl http://localhost:3000/?nocache=1
```
