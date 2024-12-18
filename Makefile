test:
	go test ./...

# Format the code
fmt:
	go fmt ./...

# Lint the code
lint:
	golangci-lint run

db:
	sqlc generate --file ./pkg/store/sqlc.yaml

vhs:
	rm -rf .idpzero
	vhs ./scripts/vhs/init.tape
	vhs ./scripts/vhs/serve.tape
	rm -rf .idpzero

watch/tailwind:
	npx --yes tailwindcss -i ./pkg/web/assets/input.css -o ./pkg/web/assets/static/styles.css --minify --watch

watch/templ:
	templ generate --watch --proxy="http://localhost:4379" --open-browser=true -v

watch/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/main" \
	--build.bin "tmp/main" \
	--build.args_bin "serve --debug" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,css" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

watch/assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "pkg/web/assets" \

watch/dbquery:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "sqlc generate --file ./pkg/store/sqlc.yaml" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_ext "sql" \
	--build.include_dir "pkg/store" \

watch:
	make -j5 watch/templ watch/server watch/tailwind watch/assets watch/dbquery

.PHONY: test fmt lint db