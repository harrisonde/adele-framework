
## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## package tests
test\:all:
	@go clean -testcache
	make test:logger test:middleware test:mux
test\:logger:
	@go test ./logger
test\:mux:
	@go test ./mux
