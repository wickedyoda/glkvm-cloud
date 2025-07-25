# Makefile

# Go binary name
BINARY_NAME = rttys

# Go 构建标志
BUILD_FLAGS := -ldflags "-s -w"

# Go build config
GO_BUILD_CMD = go build $(BUILD_FLAGS) -o $(BINARY_NAME)

# Paths
UI_DIR = ui
CONF_FILE = ./rttys.conf

.PHONY: all ui build run build-run full-run

# 只构建前端文件
ui:
	cd $(UI_DIR) && npm install && npm run build

# 只构建 Go 二进制文件
build:
	CGO_ENABLED=0 $(GO_BUILD_CMD)

# 只运行 Go 程序
run:
	./$(BINARY_NAME) -c $(CONF_FILE)

# 构建 前端和二进制并运行
build-all: ui build

# 构建 Go 二进制并运行
build-run: build run

# 构建前端、构建 Go 并运行
full-run: ui build run

# 构建docker镜像
docker-build: ui build
	docker build -t glkvm:v1 . 
