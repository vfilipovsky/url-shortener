include .env
export

GO=$(shell which go)
DB_URL=postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}
.DEFAULT_GOAL := build-and-run

.PHONY: install
install:
	${GO} get ./...

.PHONY: dev
dev:
	${GO} run cmd/server/main.go

.PHONY: build
build:
	${GO} build -o .bin/server ./cmd/server/main.go

.PHONY: build-and-run
build-and-run:
	${GO} build -o .bin/server ./cmd/server/main.go && ./.bin/server

.PHONY: access
access:
	${GO} run ./cmd/access/main.go

.PHONY: mocks
mocks:
	${GO} generate ./...

.PHONY: cron
cron:
	${GO} build -o .bin/cron ./cmd/cron/main.go && ./.bin/cron

.PHONY: new
new:
	migrate create -ext sql -dir ./migrations $(name)

.PHONY: up
up:
	migrate -source file://migrations -database ${DB_URL} up

.PHONY: down
down:
	migrate -source file://migrations -database ${DB_URL} down

.PHONY: fmt
fmt:
	${GO} fmt ./...

.PHONY: test
test:
	$(eval include .env.test)
	$(eval export $(shell sed 's/=.*//' .env.test))
	${GO} generate ./...
	${GO} test -v ./...

.PHONY: cover
cover:
	$(eval include .env.test)
	$(eval export $(shell sed 's/=.*//' .env.test))
	${GO} generate ./...
	${GO} test -coverprofile=coverage.out ./...; go tool cover -func=coverage.out