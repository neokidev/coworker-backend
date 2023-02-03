postgres:
	docker run --name management-app-demo-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it management-app-demo-db createdb --username=root --owner=root zeal_dev

dropdb:
	docker exec -it management-app-demo-db dropdb zeal_dev

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/zeal_dev?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/zeal_dev?sslmode=disable" -verbose down

migrateforce:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/zeal_dev?sslmode=disable" force 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown migrateforce sqlc test server
