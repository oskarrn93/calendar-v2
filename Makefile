.PHONY: dev
dev:
	go run .

.PHONY: build
build:
	go build -o bin/main
	
.PHONY: run
run: build
	./bin/main
