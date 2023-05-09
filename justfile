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

changelog:
	git-chglog -o CHANGELOG.md

release tag: changelog
    git add CHANGELOG.md
    git commit -m "Generated changelog for {{tag}}"
    git tag {{tag}}
    git push
    git push origin {{tag}}

