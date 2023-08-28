.PHONY: dev
dev:
	air

.PHONY: test-unit
test-unit:
	go test -v -tags=unit ./... -race -timeout=30s

.PHONY: build
build:
	go build -o main ./cmd   

