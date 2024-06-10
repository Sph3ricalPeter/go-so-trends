all: clean build run

info:
	@echo "make build: build the project"
	@echo "make run: run the project"
	@echo "make clean: remove the binary file"

tidy:
	go mod tidy

build:
	go build -o out/main cmd/go-so-trends/main.go

build-init:
	go build -o out/init scripts/init.go

run: 
	./out/main

run-init:
	./out/init

clean:
	rm -rf out

.PHONY: clean build run