# Define a variable for the test function
TARGETOS = linux
APP_NAME = frnt_package
APP_VERSION = `git rev-parse --short HEAD | xargs git describe`
GOLDFLAGS = "-X git.frontiir.net/sa-dev/$(APP_NAME)/internal.version=$(APP_VERSION)"
TEST_FUNC=TestGetCSVFile


.PHONY: format
format:
	go fmt ./...

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build
build: format mod-tidy 
	GOOS=$(TARGETOS) CGO_ENABLED=0 go build -ldflags $(GOLDFLAGS)

# Default target
.PHONY: all
all: 
	@echo "Running unit tests for frnt_package adapters"
	@echo "Running integration simulation tests"
	@go test -v ./tests/...

# Test target
.PHONY: test
test:
	@echo "TEST_FUNC is '$(TEST_FUNC)'"
	@if [ -z "$(TEST_FUNC)" ]; then \
		echo "Please provide the test function name using TEST_FUNC"; \
	else \
		go test -v ./tests -run $(TEST_FUNC); \
	fi

