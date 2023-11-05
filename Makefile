BINARY_NAME=terradiff
COVERAGE_FILE=coverage.out

.PHONY: build
build:
	go build -o $(BINARY_NAME) ./...

.PHONY: test
test:
	go test ./... -coverprofile=$(COVERAGE_FILE) -race

.PHONY: test-run
test-run: build
	./terradiff -s branch1 -w ~/Desktop/terradiff-workspace -r 'https://github.com/paper2/test-terradiff' --debug --json-log | jq -r '. | "\(.level): \(.msg)"'

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)

.PHONY: lint
lint:
	golangci-lint run