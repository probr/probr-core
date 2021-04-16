binary: go-tidy go-test go-build
quick: go-build
testcov: go-test go-test-cov

go-build:
	@echo "  >  Building binary ..."
	go build -o probr -ldflags="-X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`'"

go-test:
	@echo "  >  Validating code ..."
	golint ./...
	go vet ./...
	go test ./...

go-tidy:
	@echo "  >  Tidying go.mod ..."
	go mod tidy

go-test-cov:
	@echo "Running tests and generating coverage output ..."
	@go test ./... -coverprofile coverage.out -covermode count
	@sleep 2 # Sleeping to allow for coverage.out file to get generated
	@echo "Current test coverage : $(shell go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+') %"

release-candidate: binary
	@echo "  >  Building release candidate ..."
	go build -o probr -ldflags="-X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`' -X 'main.Prerelease=rc'"

release: binary
	@echo "  >  Building release ..."
	go build -o probr -ldflags="-X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`' -X 'main.Prerelease='"