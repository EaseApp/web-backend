APP = github.com/easeapp/web-backend

.PHONY: build test $(APP)

all: test build

$(GOPATH)/bin/golint:
	@go get -u github.com/golang/lint/golint

$(GOPATH)/bin/gb:
	@go get github.com/constabulary/gb/...

test: $(APP)

$(APP): $(GOPATH)/bin/golint $(GOPATH)/bin/gb
	gofmt -w=true $(GOPATH)/src/$@/src/
	@echo "Linting..."
	@$(GOPATH)/bin/golint $(GOPATH)/src/$@/src/...
	@echo ""
	@echo "Testing..."
	@gb test -v
	@echo ""

build: $(GOPATH)/bin/gb
	@echo "Building..."
	@gb build


