.PHONY: build run test clean rename

APP_NAME ?= app
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
GO_BUILD_CACHE ?= /tmp/templatenextjs-go-cache

build:
	cd webui && bun run build
	rm -rf src/out
	cp -R webui/out src/out
	mkdir -p build
	GOCACHE="$(GO_BUILD_CACHE)" CGO_ENABLED=0 go build -buildvcs=false -trimpath \
		-ldflags '-s -w' \
		-o build/$(APP_NAME) ./src
	@echo "Built single binary: build/$(APP_NAME)"

run:
	@echo "Starting development servers (Ctrl+C to stop both)..."
	@trap 'kill 0' INT TERM EXIT; \
	cd webui && bun dev & \
	DEV=true GOCACHE="$(GO_BUILD_CACHE)" go run -buildvcs=false ./src

test:
	GOCACHE="$(GO_BUILD_CACHE)" go test -buildvcs=false ./src/... ./tools/... -v

clean:
	rm -rf build/
	rm -rf src/out
	rm -rf webui/.next webui/out
	mkdir -p src/out
	touch src/out/.gitkeep

# Rename the Go module after cloning the template.
# Usage: make rename NAME=github.com/you/new_project
rename:
	@test -n "$(NAME)" || (echo "Usage: make rename NAME=github.com/you/new_project" && exit 1)
	GOCACHE="$(GO_BUILD_CACHE)" go run -buildvcs=false ./tools/rename -name "$(NAME)"
	GOCACHE="$(GO_BUILD_CACHE)" go build -buildvcs=false ./...
	@echo "Done. Module is now $(NAME)"
