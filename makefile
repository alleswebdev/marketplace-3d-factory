PG_DSN="user=postgres password=example host=localhost port=5432 database=test sslmode=disable"

.PHONY: run
run:
	go run cmd/grpc-server/main.go

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