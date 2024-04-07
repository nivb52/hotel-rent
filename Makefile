# WIN
SHELL=cmd
DOCKER_PORT=5000
DOCKER_TAG=hotel_rent_app

dev:
	@air -c .air.conf
serve : dev

run: build 
	@echo Starting app..
	go run main.go

start: db_up
	serve

debug_back:
	@echo Debugging...
	@ D:\d-dev\goworkspace\bin\dlv.exe dap --listen=127.0.0.1:62712 from ${pwd}

########### TEST ###########
test:
	go test -v ./...

########### DB ###########
# SEED
seed: 
	go run ./scripts/seed.go

# CONSINSTENT DB
db_up: 
	@echo "Docker compose up: db image..."
	@docker compose --env-file .env --env-file .env.local -p hotel_rent  up -d  
	@echo "Docker db up!"

db_down:
	@echo "Docker compose down: db image..."
	@docker-compose  -p hotel_rent down
	@echo "Docker db down!"

# TEST DB
dbtest_up: 
	@echo "Docker compose up: dbtest image..."
	@docker compose --env-file .env --env-file .env.local --env-file .env.test.local -p hotel_rent_test  up -d  
	@echo "Docker db up!"

dbtest_down:
	@echo "Docker compose down: dbtest image..."
	@docker-compose  -p hotel_rent_test down
	@echo "Docker dbtest down!"
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

docker:
	@echo Building Docker
	@docker build -t $(DOCKER_TAG) .
	@echo "Running App inside Docker at $(DOCKER_PORT)"
	@docker run -p $(DOCKER_PORT):5000 $(DOCKER_TAG)

########### Make Rest File  ###########
rest: 
	$(shell powershell  ./generate_main_rest_file.ps1)
