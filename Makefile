postgres:
	docker run --name hamed-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it hamed-postgres createdb --username=root --owner=root hamed_bank

dropdb:
	docker exec -it hamed-postgres dropdb hamed_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/hamed_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/hamed_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go
.PHONY: postgres createdb dropdb migratedown migrateup sqlc