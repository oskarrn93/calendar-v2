.PHONY: install
install:
	go mod download 

.PHONY: dev
dev:
	godotenv -f .env go run .

.PHONY: build
build:
	go build -o ./bin ./...
	chmod +x ./bin/lambda
	chmod +x ./bin/nba
	

.PHONY: run-lambda
run-lambda: build
	godotenv -f .env ./bin/lambda

.PHONY: run-nba
run-nba: build
	godotenv -f .env ./bin/nba

.PHONY: deploy
deploy:
	cdk deploy

.PHONY: test
test:
	godotenv -f .env.test go test -v ./...

.PHONY: test-ci
test-ci:
	godotenv -f .env.test go test --failfast ./...

.PHONY: update
update:
	go get -u ./...
	go mod tidy

