# WIN
SHELL=cmd

dev:
	@air -c .air.conf
serve : dev

run: build 
	@echo Starting app..
	go run main.go

start: db_up
	run

debug_back:
	@echo Debugging...
	@ D:\d-dev\goworkspace\bin\dlv.exe dap --listen=127.0.0.1:62712 from ${pwd}

########### TEST ###########
test:
	go test -v ./...

########### DB ###########
# CONSINSTENT DB
db_up: db_down
	@echo "Docker compose up: db image..."
	@docker-compose --env-file .env --env-file local.env -p hotel_rent  up -d  
	@echo "Docker db up!"

db_down:
	@echo "Docker compose down: db image..."
	@docker-compose  -p hotel_rent down
	@echo "Docker db down!"
# TEMP DB
dbtemp_up:
	@echo running db (temp mode)
	@cmd docker run --name mongodb -d mongo:latest -p 27017:27017

dbtemp_down:
	@echo shuting down db (temp mode)
	@cmd docker stop --name mongodb -d mongo:latest -p 27017:27017

########### BUILD  ###########
build:
	@echo Building app..
	go build -o bin/main main.go
