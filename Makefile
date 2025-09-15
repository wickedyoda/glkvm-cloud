# Makefile

# Go binary name
BINARY_NAME = rttys

# Go build flags
BUILD_FLAGS := -ldflags "-s -w"

# Go build command
GO_BUILD_CMD = go build $(BUILD_FLAGS) -o $(BINARY_NAME)

# Paths
UI_DIR = ui
CONF_FILE = ./rttys.conf

.PHONY: all ui build run build-run full-run

# Build frontend files only
ui:
	cd $(UI_DIR) && npm install && npm run build

# Build Go binary only
build:
	CGO_ENABLED=0 $(GO_BUILD_CMD)

# Run Go program only
run:
	./$(BINARY_NAME) -c $(CONF_FILE)

# Build frontend and Go binary
build-all: ui build

# Build Go binary and run
build-run: build run

# Build frontend, build Go binary, and run
full-run: ui build run

# Build Docker image without updating ui
docker-build: build
	docker build -t glkvm:v1 .

# Full Build Docker image
docker-fullbuild: ui build
	docker build -t glkvm:v1 .