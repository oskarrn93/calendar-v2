.PHONY: install
install:
	go mod download 

.PHONY: dev
dev:
	go run .

.PHONY: build
build:
	go build -o ./bin ./...
	chmod +x ./bin/lambda
	chmod +x ./bin/nba
	
.PHONY: run-lambda
run-lambda: build
	./bin/lambda

.PHONY: run-nba
run-nba: build
	./bin/nba

.PHONY: test
test:
	go test -v ./...

.PHONY: test-ci
test-ci:
	go test --failfast ./...

.PHONY: update
update:
	go get -u ./...
	go mod tidy
