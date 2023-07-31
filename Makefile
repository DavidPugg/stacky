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
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" up

migrate-force:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" force 1

migrate-down:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" down 1

mysql:
	docker run --name stacky-mysql -e MYSQL_ROOT_PASSWORD=password -p 3306:3306 -d mysql

mysql-createdb:
	docker exec -it stacky-mysql mysql -uroot -ppassword -e "CREATE DATABASE stacky;"

.PHONY: build-css start migrate-up migrate-force migrate-down watch-css dev yarn mysql mysql-createdb