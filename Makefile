
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
	make test:logger test:mailer test:middleware test:mux test:session
test\:cache:
	@go test ./cache/...
test\:database:
	@go test ./database/...
test\:filesystem:
	@go test ./filesystem/...
test\:helpers:
	@go test ./helpers
test\:logger:
	@go test ./logger
test\:middleware:
	@go test ./middleware
test\:mailer:
	@go test ./mailer
test\:mux:
	@go test ./mux
test\:session:
	@go test ./session
test\:render:
	@go test ./render
