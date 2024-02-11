postgres: 
	docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:alpine3.19
createdb: 
	docker exec -it postgres13 createdb --username=root --owner=root simple_bank
dropdb: 
	docker exec -it postgres13 dropdb simple_bank
migrateup: 
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown: 
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down
test: 
	go test -v -cover ./...
sqlc: 
	sqlc generate 
.PHONY: postgres createdb dropdb migrateup migratedown test