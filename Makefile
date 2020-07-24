.PHONY: all
all: setup mock.generate test build

.PHONY: setup
setup:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(HOME)/bin latest

.PHONY: build
build:
	go build .

.PHONY: test
test:
	go test ./...

.PHONY: fix
fix:
	$(GOPATH)/bin/golangci-lint run --fix

.PHONY: lint
lint:
	$(GOPATH)/bin/golangci-lint run

.PHONY: mock.generate
mock.generate:
	go get -u github.com/vektra/mockery/cmd/mockery
	mockery -dir=pkg/ -output=pkg/testlib/mocks -case underscore -all
	mockery -dir=util/ -output=util/testlib/mocks -case underscore -all
