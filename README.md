# GoLang
<!-- docker -->
postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable
0. Build : docker  build -t "simplebank:latest"  
1. docker run --name simplebank --network bank-network -p 8080:8080 -e DB_SOURCE="postgresql://root:123456@postgres13:5432/simple_bank?sslmode=disable" -e SERVER_ADDRESS="0.0.0.0:8080" simplebank:latest
3. Connect container to database container through network bridge 
 -  docker network create bank-network
 - docker network connect bank-network postgres13
4. Sometime you want to inspect : 
 - Network :  docker network inspect bank-network. 
 - Container :  docker container inspect <docker-name>
5. Create a new db migration : 
  - migrate create -ext sql -dir db/migration -seq 
6. docker exec -it <container-name-or-id> <shell-executable> => psql -U username -d database_name (to connect to databse)

<!-- Todo list -->
Make UpdateUserTestAPI for both Gin and gRPC
