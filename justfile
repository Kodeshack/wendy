_default:
    @just --list

test +flags="-failfast":
	go test {{flags}} ./...

lint:
	staticcheck ./...
	golangci-lint run ./...

fmt:
	@go fmt ./...

clean:
	go clean -cache

