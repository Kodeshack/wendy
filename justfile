_default:
    @just --list

test +flags="-failfast -race":
	go test {{flags}} ./...

test-watch *flags:
	gotestsum --watch -- {{flags}}

lint:
	staticcheck ./...
	golangci-lint run ./...

fmt:
	@go fmt ./...

clean:
	go clean -cache

changelog:
	git-chglog -o CHANGELOG.md

release tag:
    just changelog
    git add CHANGELOG.md
    git commit -m "Generated changelog for {{tag}}"
    git tag {{tag}}
    git push
    git push origin {{tag}}

