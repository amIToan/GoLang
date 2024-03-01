postgres: 
	docker run --name postgres13 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:alpine3.19
createdb: 
	docker exec -it postgres13 createdb --username=root --owner=root simple_bank
dropdb: 
	docker exec -it postgres13 dropdb simple_bank
migrateup: 
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1: 
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown: 
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1: 
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
test: 
	go test -v -cover ./...
sqlc: 
	sqlc generate 
mock : 
	mockgen --build_flags=--mod=mod -package mockdb -destination ./db/mock/store.go sgithub.com/techschool/simplebank/db/sqlc Store 
server: 
	go run main.go
.PHONY: postgres createdb dropdb migrateup migratedown migratedown1 migrateup1 test server