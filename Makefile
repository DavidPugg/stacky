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

seed:
	go run cmd/seed/main.go

migrate:
	atlas schema apply -u "${DB_URL}" --to file://schema.hcl

db:
	docker run --name some-postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres

db-start:
	docker start some-postgres

db-stop:
	docker stop some-postgres

db-create:
	docker exec -it some-postgres psql -U postgres -c "CREATE DATABASE stacky;"

.PHONY: build-css start migrate-up migrate-force migrate-down watch-css dev yarn build db db-start db-stop db-create