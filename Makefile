postgres: 
	docker run --name postgres13 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:alpine3.19
createdb: 
	docker exec -it postgres13 createdb --username=root --owner=root simple_bank
dropdb: 
	docker exec -it postgres13 dropdb simple_bank
migrateup: 
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateup1: 
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown: 
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migratedown1: 
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
test: 
	go test -v -cover ./...
sqlc: 
	sqlc generate 
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
mock : 
	mockgen --build_flags=--mod=mod -package mockdb -destination ./db/mock/store.go sgithub.com/techschool/simplebank/db/sqlc Store 
	mockgen --build_flags=--mod=mod -package mockwk -destination ./worker/mock/distributor.go sgithub.com/techschool/simplebank/worker TaskDistributor
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb \
	--grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	./proto/*.proto
	statik -src=./doc/swagger -dest=./doc
evans: 
	evans --host localhost --port 8000 -r repl
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
server: 
	go run main.go
.PHONY: postgres createdb dropdb migrateup migratedown migratedown1 migrateup1 test server db_docs db_schema proto new_migration