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

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'