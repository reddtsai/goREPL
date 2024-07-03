.PHONY:

run:
	go run main.go

gen-mock:
	go generate ./...

test:
	go test -v ./...

test-race:
	go test -race

test-coverage:
	go test -cover ./...

release:
	goreleaser release --rm-dist

local-release:
	goreleaser release --snapshot --clean