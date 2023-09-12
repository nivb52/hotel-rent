build:
  @go build -o bin/api

dev:
  @air -c .air.linux.toml

run: build 
  @go run ./bin/api

test:
  @go test -v ./...

db_up:
  @echo running db
  @echo docker run --name mongodb -d mongo:latest -p 27017:27017

db_down:
  @echo shuting down db
  @cmd docker stop --name mongodb -d mongo:latest -p 27017:27017