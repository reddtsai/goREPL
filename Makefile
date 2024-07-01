.PHONY:

run:
	go run main.go

gen-mock:
	go generate ./...

test:
	go test -v ./...

test-coverage:
	go test -cover ./...