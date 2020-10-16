# Local Development Maintenance/Building Targets

.PHONY: deps
deps:
	go mod tidy
	go mod download

.PHONY: bin
bin:
	mkdir -p bin/

.PHONY: build
build: deps bin
	go build -o bin/loraficationd cmd/loraficationd/loraficationd.go

# Docker/Docker Compose Targets

.PHONY: run
run: stop up

.PHONY: up
up:
	docker-compose -f docker-compose.yml up -d --build

.PHONY: stop
stop:
	docker-compose -f docker-compose.yml stop

.PHONY: down
down:
	docker-compose -f docker-compose.yml down

.PHONY: logs
logs:
	docker-compose -f docker-compose.yml logs