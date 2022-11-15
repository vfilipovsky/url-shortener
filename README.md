# Url Shortener

## Installation and running

```
docker-compose up -d
go install github.com/golang/mock/mockgen@v1.6.0
make install
make up

make
```

## Tests

```
make test
```

## go fmt

```
make fmt
```

## Migrations

Install: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

To create a new migration use

```
make new name=your_migration_name
```

To apply your migration use

```
make up
```

To downgrade the latest migration use

```
make down
```
