ifeq ($(wildcard .env),)
$(warning .env file not found, skipping inclusion.)
else
include .env
endif

watch-css:
	yarn run watch-css

dev:
	air

yarn:
	yarn

build-css: yarn
	yarn run build-css

build:
	go build -o web cmd/web/main.go

start: build
	./web

migrate-up:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-force:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" force 1

migrate-down:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down 1

db:
	docker run --name some-postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres

db-start:
	docker start some-postgres

db-stop:
	docker stop some-postgres

db-create:
	docker exec -it some-postgres psql -U postgres -c "CREATE DATABASE stacky;"

.PHONY: build-css start migrate-up migrate-force migrate-down watch-css dev yarn build db db-start db-stop db-create