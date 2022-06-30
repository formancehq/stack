BINARY_NAME=webhooks
PKG=./...
FAILFAST=-failfast
TIMEOUT=10m
RUN=".*"

all: lint test

build:
	go build -o $(BINARY_NAME)

install: build
	cp $(BINARY_NAME) $(shell go env GOPATH)/bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(shell go env GOPATH)/bin latest
	golangci-lint --version

lint:
	golangci-lint run -v --fix

test:
	go test -v $(FAILFAST) -coverpkg $(PKG) -coverprofile coverage.out -covermode atomic -run $(RUN) -timeout $(TIMEOUT) $(PKG) \
		| sed ''/PASS/s//$(shell printf "\033[32mPASS\033[0m")/'' \
		| sed ''/FAIL/s//$(shell printf "\033[31mFAIL\033[0m")/'' \
		| sed ''/RUN/s//$(shell printf "\033[34mRUN\033[0m")/''
	@go tool cover -html=coverage.out -o coverage.html
	@echo "To open the html coverage file, use one of the following commands:"
	@echo "open coverage.html on mac"
	@echo "xdg-open coverage.html on linux"

bench:
	go test -v -bench=. -run=^a $(PKG)

clean:
	go clean
	rm -f $(BINARY_NAME) $(COVERAGE_FILE)
