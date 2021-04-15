binary: go-tidy go-test go-build
quick: go-build

go-build:
	@echo "  >  Building binary..."
	go build -o probr main.go

go-test:
	@echo "  >  Validating code..."
	golint ./...
	go vet ./...
	go test ./...

go-tidy:
	@echo "  >  Tidying go.mod ..."
	go mod tidy

go-test-cov:
	@echo "Running tests and generating coverage output"
	@go test ./... -coverprofile coverage.out -covermode count
	@echo "Current test coverage : $(shell go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+') %"