# Run unit tests
.PHONY: test
test-unit:
	go test -v -tags=unit ./... -timeout=5s