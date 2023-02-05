postgres:
	docker run --name coworker-db --network coworker-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it coworker-db createdb --username=postgres --owner=postgres postgres

dropdb:
	docker exec -it coworker-db dropdb postgres

migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable" -verbose down

migrateforce:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable" force 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ot07/coworker-backend/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateforce sqlc test server mock
