postgres:
	docker run --name coworker-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it coworker-db createdb --username=root --owner=root coworker

dropdb:
	docker exec -it coworker-db dropdb coworker

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/coworker?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/coworker?sslmode=disable" -verbose down

migrateforce:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/coworker?sslmode=disable" force 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ot07/coworker-backend/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateforce sqlc test server mock
