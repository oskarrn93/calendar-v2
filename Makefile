.PHONY: install
install:
	go mod download 

.PHONY: build
build:
	go build -o ./bin ./...

# e.g. make run-nba; the lambda binary only runs inside the Lambda runtime.
run-%: build
	godotenv -f .env ./bin/$*

.PHONY: deploy
deploy:
	cd infra && terraform apply

.PHONY: test
test:
	go test -v ./...

.PHONY: test-ci
test-ci:
	go test -failfast ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: update
update:
	go get -u ./...
	go mod tidy

