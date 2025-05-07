build:
	@go build -o bin/htmxlogin

run: build
	@./bin/htmxlogin

test:
	@go test -v ./...