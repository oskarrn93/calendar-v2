.PHONY: install
install:
	go mod download 

.PHONY: build
build:
	go build -o ./bin ./...

# Runs all jobs; select specific ones with e.g. make run JOBS=nba,f1.
# The lambda binary only runs inside the Lambda runtime.
.PHONY: run
run: build
	godotenv -f .env ./bin/calendar $(if $(JOBS),-jobs=$(JOBS))

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
	@golangci-lint version --short | grep -qxF "$$(cat .golangci-lint-version)" || { echo "golangci-lint $$(golangci-lint version --short) does not match .golangci-lint-version ($$(cat .golangci-lint-version))"; exit 1; }
	golangci-lint run ./...

.PHONY: install-lint
install-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(HOME)/.local/bin "v$$(cat .golangci-lint-version)"

.PHONY: update
update:
	go get -u ./...
	go mod tidy

