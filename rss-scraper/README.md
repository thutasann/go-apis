# Golang Web Server and RSS Scraper

- go chi router
- postgres
- sqlc
- goose

## Scripts

**Dev Run**

```bash
make dev
```

**Build**

```bash
make build
```

**Prod Run**

```bash
make run
```

**Vendor**

```bash
go mode vendor
```

```bash
go mode tidy
```

## Sqlc and Goose

**install**

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

brew install sqlc
```

```bash
brew install goose
```

## Goose

```bash
cd sql/schema

goose postgres postgres://postgres:@localhost:5432/rssagg up
```
