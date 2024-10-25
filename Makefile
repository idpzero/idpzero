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
	npx tailwindcss -i ./pkg/web/assets/css/input.css -o ./pkg/web/assets/css/output.css
	templ generate

.PHONY: test fmt lint