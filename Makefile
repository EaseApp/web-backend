APP = github.com/EaseApp/web-backend

.PHONY: build test $(APP)

all: test build

$(GOPATH)/bin/golint:
	@go get -u github.com/golang/lint/golint

$(GOPATH)/bin/gb:
	@go get github.com/constabulary/gb/...

test: $(APP)

$(APP): $(GOPATH)/bin/golint $(GOPATH)/bin/gb
	@gofmt -w=true $(GOPATH)/src/$@/src/
	@echo "Linting..."
	@$(GOPATH)/bin/golint $(GOPATH)/src/$@/src/...
	@echo ""
	@echo "Testing..."
	@$(GOPATH)/bin/gb test -v
	@echo ""

build: $(GOPATH)/bin/gb
	@echo "Building..."
	@$(GOPATH)/bin/gb build


