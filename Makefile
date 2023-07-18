include .env

build-css:
	yarn run build-css

start:
	air

migrate-up:
	docker run -v $(shell pwd)/server/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_URL}" up

migrate-force:
	docker run -v $(shell pwd)/server/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_URL}" force 1

migrate-down:
	docker run -v $(shell pwd)/server/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_URL}" down

.PHONY: build-css start migrate-up migrate-force migrate-down