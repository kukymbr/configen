GOLANGCI_LINT_VERSION := 2.5.0

GIT_VERSION := $(shell git describe --tags --always)
GIT_REVISION := $(shell git rev-parse HEAD)
VERSION_PACKAGE := "github.com/kukymbr/configen/internal/version"

GO_BUILD_LDFLAGS := "-X $(VERSION_PACKAGE).Version=$(GIT_VERSION) \
                -X $(VERSION_PACKAGE).Revision=$(GIT_REVISION) \
                -X $(VERSION_PACKAGE).BuiltAt=$(shell date -u +%Y%m%d%H%M%S)"
GO_BUILD_ARGS := $(build_arguments) --ldflags $(GO_BUILD_LDFLAGS)

all:
	$(MAKE) clean
	$(MAKE) prepare
	$(MAKE) validate
	$(MAKE) build

prepare:
	go mod tidy
	go mod download
	$(MAKE) fmt

validate:
	go vet ./...
	$(MAKE) lint
	$(MAKE) test

build:
	go build $(GO_BUILD_ARGS) -o bin/configen ./cmd/configen

./bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "./bin" "v${GOLANGCI_LINT_VERSION}";

lint:
	$(MAKE) ./bin/golangci-lint
	./bin/golangci-lint run ./...

fmt:
	$(MAKE) ./bin/golangci-lint
	./bin/golangci-lint fmt ./...

test:
	go test -coverprofile=coverage_out ./...
	go tool cover -func=coverage_out
	go tool cover -html=coverage_out -o coverage.html
	rm -f coverage_out

test_report:
	go test -coverprofile=coverage_out -v 2>&1 ./... | go-junit-report -set-exit-code -iocopy -out junit.report.xml
	go tool cover -func=coverage_out
	go tool cover -html=coverage_out -o coverage.html

test_short:
	go test -short ./...

clean:
	go clean

generate_example:
	go run cmd/configen/main.go --source=example --struct=config --yaml=example/config.yaml --env=example/config.env --go=example/config.gen.go
