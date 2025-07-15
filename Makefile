build:
	@go build -o bin/htmxlogin cmd/main.go

run: build
	@./bin/htmxlogin

test:
	@go test -v ./...