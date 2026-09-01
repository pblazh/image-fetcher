GO_CMD := "go"

build: build-darwin-arm64

clean:
	rm -rf bin
	rm -f coverage.out coverage.html
	rm -f plugins/tabula.vscode/node_modules

build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 {{GO_CMD}} build -o bin/darwin/arm64/image-fetcher ./cmd/cli

setup:
  {{GO_CMD}} mod download

test:
	{{GO_CMD}} test ./...

coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Opening coverage report in browser..."
	@open coverage.html || xdg-open coverage.html || start coverage.html


install:
	{{GO_CMD}} install ./cmd/cli

lint:
	golangci-lint run
