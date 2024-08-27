export CGO_ENABLED=0

.PHONY: run
run:
	go run cmd/main.go

.PHONY: build
build:
	docker build -t 3d-factory .

build:
	docker 3d-factory | bzip2 | ssh user@your-vps-ip 'bunzip2 | docker load'
	ssh user@your-vps-ip "docker run -d -p 8080:8080 --name my-go-project-container my-go-project"

.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $$3}'


.PHONY: goose-up
goose-up:
	goose -dir migrations \
      postgres $(PG_DSN) \
      up

.PHONY: goose-status
goose-status:
	goose -dir migrations \
      postgres $(PG_DSN) \
      status

.PHONY: goose-down
goose-down:
	goose -dir migrations \
      postgres $(PG_DSN) \
      down

.PHONY: goose-create
goose-create:
	goose -dir migrations \
	postgres $(PG_DSN) \
	create $(name) sql

.PHONY: lint
lint:
	golangci-lint run \
		--new-from-rev=origin/master \
		--config=.golangci.pipeline.yaml \
		./...