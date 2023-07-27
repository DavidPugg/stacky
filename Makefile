ifeq ($(wildcard .env),)
$(warning .env file not found, skipping inclusion.)
else
include .env
endif

watch-css:
	yarn run watch-css

dev:
	air

build-css:
	yarn run build-css

build: build-css
	go build -o web cmd/web/main.go

start: build
	./web

migrate-up:
	docker run -v $(shell pwd)/server/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_URL}" up

migrate-force:
	docker run -v $(shell pwd)/server/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_URL}" force 1

migrate-down:
	docker run -v $(shell pwd)/server/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_URL}" down

.PHONY: build-css start migrate-up migrate-force migrate-down watch-css dev