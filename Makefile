## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the command line tool celeritas and copies it to myapp
build_cli_copy:
	@go build -o ../myapp/adele ./cmd/cli

## builds the command line tool /dist dir
build_cli:
	@go build -o ./dist/adele ./cmd/cli
