LOCAL_BIN:=$(CURDIR)/bin


.PHONY: goose-create
goose-create:
	 goose create $(name) go

.PHONY: goose-up
goose-up:
	./goose postgres up -dir ./migrations

.PHONY: goose-down
goose-down:
	./goose postgres down -dir ./migrations


.PHONY: build
build:
	go build -v -o ./ ./cmd/...

#.PHONY: test
#test:
#	go test -v ./...
