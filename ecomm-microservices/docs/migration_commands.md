# Migration Commands

### Init

```bash
docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.17.0 create -ext sql -dir db/migrations init_schema
```

### Interactive

```bash
docker exec -i ecomm-mysql mysql -uroot -ppassword <<< "CREATE DATABASE ecom;"
```

```bash
docker run -it --rm --network host --volume ./db:/db migrate/migrate:v4.17.0 -path=/db/migrations -database "mysql://root:password@tcp(localhost:3306)/ecom" up
```