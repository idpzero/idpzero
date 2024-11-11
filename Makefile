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

db:
	sqlc generate --file ./pkg/store/sqlc.yaml

.PHONY: test fmt lint web db