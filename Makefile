.PHONY: test-unit
test-unit:
	go test -v -tags=unit ./... -race -timeout=5s

.PHONY: build
build:
	go build -o main ./cmd   
