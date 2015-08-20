APP=github.com/EaseApp/web-backend
EXECUTABLES:=$(basename $(notdir $(wildcard src/cmd/*/*.go)))

# Turn on Go 1.5 vendoring.
export GO15VENDOREXPERIMENT=1

.PHONY: build test $(APP)

all: test build

dependencies:
	@echo "Getting dependencies..."
	@go get ./...

$(GOPATH)/bin/golint:
	@go get -u github.com/golang/lint/golint

test: $(APP)

$(APP): dependencies $(GOPATH)/bin/golint
	@gofmt -w=true $(GOPATH)src/$@/src/
	@echo "Linting..."
	@$(GOPATH)/bin/golint ./...
	@echo ""
	@echo "Testing..."
	@go test -v ./...
	@echo ""

build: $(EXECUTABLES)

$(EXECUTABLES): dependencies
	@echo "Building executable $@..."
	@go build -o $(addprefix bin/, $@) ./src/cmd/$@

