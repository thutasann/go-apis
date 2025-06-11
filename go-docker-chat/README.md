# Go Docker Chat app

## Postgres

```bash
docker pull postgres:15-alpine
```

```bash
docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine
```
