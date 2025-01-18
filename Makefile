dev:
	air

install:
	go mod tidy


test:
	go test -cover ./internal/...