# 项目名称
PROJECT_NAME := k8s-offline-installer

# 可执行文件名称
BINARY_NAME := k8s-installer

# 主程序入口路径
ENTRY_POINT := ./cmd/k8sinstall

# 配置文件路径（默认）
CONFIG_PATH := ./configs/config.yaml

# 目标Linux架构（默认x86_64）
ARCH := amd64

# Go编译参数
GO_BUILD_FLAGS := CGO_ENABLED=0
GOOS ?= linux  # 默认编译为Linux版本

# 帮助信息
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build        编译当前系统的可执行文件"
	@echo "  build-linux  编译Linux系统的可执行文件（默认x86_64）"
	@echo "  build-arm64  编译Linux系统的ARM64架构可执行文件"
	@echo "  run          本地运行程序（使用默认配置）"
	@echo "  run-dev      本地运行程序（显示详细日志）"
	@echo "  clean        清理编译产物"
	@echo "  deps         下载项目依赖"
	@echo "  fmt          格式化代码"
	@echo "  scp          上传可执行文件到目标主机（需修改TARGET_HOST）"

# 编译当前系统的可执行文件
.PHONY: build
build:
	$(GO_BUILD_FLAGS) go build -o $(BINARY_NAME) $(ENTRY_POINT)
	@echo "编译完成：$(BINARY_NAME)（当前系统：$(shell go env GOOS)/$(shell go env GOARCH)）"

# 编译Linux x86_64架构可执行文件
.PHONY: build-linux
build-linux:
	$(GO_BUILD_FLAGS) GOOS=linux GOARCH=$(ARCH) go build -o $(BINARY_NAME)-linux-$(ARCH) $(ENTRY_POINT)
	@echo "编译完成：$(BINARY_NAME)-linux-$(ARCH)"

# 编译Linux ARM64架构可执行文件（适用于ARM服务器）
.PHONY: build-arm64
build-arm64:
	$(GO_BUILD_FLAGS) GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 $(ENTRY_POINT)
	@echo "编译完成：$(BINARY_NAME)-linux-arm64"

# 本地运行程序（使用默认配置）
.PHONY: run
run:
	go run $(ENTRY_POINT) --config $(CONFIG_PATH)

# 本地运行程序（显示调试日志）
.PHONY: run-dev
run-dev:
	go run $(ENTRY_POINT) --config $(CONFIG_PATH) -v=5  # 假设程序支持-v参数控制日志级别

# 清理编译产物
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux-*
	rm -f ./k8s_install.logger  # 清理日志文件
	@echo "清理完成"

# 下载依赖
.PHONY: deps
deps:
	go mod tidy
	@echo "依赖已更新"

# 格式化代码（符合Go规范）
.PHONY: fmt
fmt:
	go fmt ./...
	@echo "代码已格式化"

# 上传可执行文件到目标主机（需修改TARGET_HOST为实际IP）
.PHONY: scp
scp: build-linux
	@if [ -z "$(TARGET_HOST)" ]; then \
		echo "请指定目标主机：make scp TARGET_HOST=root@192.168.110.142"; \
		exit 1; \
	fi
	scp $(BINARY_NAME)-linux-$(ARCH) $(TARGET_HOST):/root/
	scp $(CONFIG_PATH) $(TARGET_HOST):/root/configs/  # 同步配置文件
	@echo "已上传到 $(TARGET_HOST)"
