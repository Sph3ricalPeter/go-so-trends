all: clean build run

info:
	@echo "make build: build the project"
	@echo "make run: run the project"
	@echo "make clean: remove the binary file"

tidy:
	go mod tidy

build:
	go build -o out/main cmd/go-so-trends/main.go

run: 
	./out/main

clean:
	rm -rf out

.PHONY: clean build run