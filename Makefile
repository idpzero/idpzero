test:
	go test ./...

# Format the code
fmt:
	go fmt ./...

# Lint the code
lint:
	golangci-lint run

# generate web resources
web:
	npx tailwindcss -i ./web/css/input.css -o ./pkg/web/assets/css/styles.css
	templ generate

.PHONY: test fmt lint web