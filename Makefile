.PHONY: dev
dev:
	go run .

.PHONY: build
build:
	go build -o bin/main
	
.PHONY: run
run: build
	./bin/main

.PHONY: test
test:
	go test -v ./...

.PHONY: update
update:
	go get -u ./...
	go mod tidy
