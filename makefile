.PHONY: build

deps:
	go mod tidy
	go mod vendor

bin:
	mkdir -p bin/

build: deps bin
	go build -o bin/loraficationd cmd/loraficationd/loraficationd.go