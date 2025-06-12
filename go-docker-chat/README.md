# Go Docker Chat app

## Postgres

**DB URI**

```bash
postgresql://root:password@localhost:5433/go-chat?sslmode=disable
```

```bash
docker pull postgres:15-alpine
```

```bash
docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine
```

**Create DB**

```bash
docker exec -it postgres15 createdb --username=root --owner=root go-chat
```

**PSQL Playground**

```bash
docker exec -it postgres15 psql

\l

\c go-chat

\d

exit
```

## Golang Migrate

[Golang Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

```bash
$ brew install golang-migrate
```

```bash
migrate create -ext sql -dir db/migrations add_users_table
```

```bash
migrate -path db/migrations -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose up
```
