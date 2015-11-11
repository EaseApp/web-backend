APP=github.com/EaseApp/web-backend

# Turn on Go 1.5 vendoring.
export GO15VENDOREXPERIMENT=1

.PHONY: build test $(APP)

all: test build

dependencies:
	@echo "Getting dependencies..."
	@go get -t ./...

$(GOPATH)/bin/golint:
	@go get -u github.com/golang/lint/golint

test: $(APP)

$(APP): dependencies $(GOPATH)/bin/golint
	@gofmt -w=true $(GOPATH)/src/$@/src/
	@echo "Linting..."
	@$(GOPATH)/bin/golint ./...
	@echo ""
	@echo "Testing..."
	@if go test -v ./... ; then \
	  echo "TESTS PASSED!!!!!" ; \
	  else echo "TESTS FAILED :'("; fi
	@echo ""

build: dependencies
	@echo "Building executable main..."
	@go build -o bin/main main.go


$(GOPATH)/bin/gin:
	@go get github.com/codegangsta/gin

dev-server: dependencies $(GOPATH)/bin/gin
	@echo "Starting dev server..."
	@$(GOPATH)/bin/gin
