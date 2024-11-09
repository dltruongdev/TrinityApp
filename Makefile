pullpostgres:
	docker pull postgres:15.8-alpine3.20

postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -d postgres:15.8-alpine3.20

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root trinity_app

dropdb:
	docker exec -it postgres15 dropdb trinity_app

migrate-up:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5432/trinity_app?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5432/trinity_app?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: createdb dropdb pullpostgres postgres migrate-up migrate-down sqlc