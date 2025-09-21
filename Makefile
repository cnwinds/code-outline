# CodeCartographer Makefile

# 变量定义
BINARY_NAME=contextgen
MAIN_PATH=./cmd/contextgen
BUILD_DIR=./build
VERSION=v1.0.0
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# 默认目标
.PHONY: all
all: clean build

# 构建二进制文件
.PHONY: build
build:
	@echo "🔨 构建 CodeCartographer..."
	@mkdir -p ${BUILD_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PATH}
	@echo "✅ 构建完成: ${BUILD_DIR}/${BINARY_NAME}"

# 跨平台构建
.PHONY: build-all
build-all:
	@echo "🔨 开始跨平台构建..."
	@mkdir -p ${BUILD_DIR}
	
	# Windows
	@echo "构建 Windows 版本..."
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ${MAIN_PATH}
	
	# Linux
	@echo "构建 Linux 版本..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ${MAIN_PATH}
	
	# macOS
	@echo "构建 macOS 版本..."
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 ${MAIN_PATH}
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 ${MAIN_PATH}
	
	@echo "✅ 跨平台构建完成"

# 运行程序
.PHONY: run
run: build
	@echo "🚀 运行 CodeCartographer..."
	${BUILD_DIR}/${BINARY_NAME} generate --path .

# 测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 基准测试
.PHONY: bench
bench:
	@echo "⚡ 运行基准测试..."
	go test -bench=. -benchmem ./...

# 代码格式化
.PHONY: fmt
fmt:
	@echo "📝 格式化代码..."
	go fmt ./...

# 代码检查
.PHONY: lint
lint:
	@echo "🔍 运行代码检查..."
	golangci-lint run

# 代码整理
.PHONY: tidy
tidy:
	@echo "🧹 整理依赖..."
	go mod tidy

# 清理构建文件
.PHONY: clean
clean:
	@echo "🧽 清理构建文件..."
	rm -rf ${BUILD_DIR}
	rm -f project_context.json

# 安装到系统
.PHONY: install
install: build
	@echo "📦 安装到系统..."
	cp ${BUILD_DIR}/${BINARY_NAME} /usr/local/bin/
	@echo "✅ 安装完成"

# 卸载
.PHONY: uninstall
uninstall:
	@echo "🗑️  卸载..."
	rm -f /usr/local/bin/${BINARY_NAME}
	@echo "✅ 卸载完成"

# 创建语法目录
.PHONY: setup-grammars
setup-grammars:
	@echo "📁 创建语法目录..."
	mkdir -p grammars
	@echo "⚠️  请手动下载并编译Tree-sitter语法文件到grammars目录"
	@echo "   参考: https://github.com/tree-sitter/tree-sitter"

# 生成示例项目上下文
.PHONY: example
example: build
	@echo "📋 生成示例项目上下文..."
	${BUILD_DIR}/${BINARY_NAME} generate --path . --output example_context.json
	@echo "✅ 示例文件生成完成: example_context.json"

# 显示帮助
.PHONY: help
help:
	@echo "CodeCartographer Makefile 命令:"
	@echo "  build        - 构建二进制文件"
	@echo "  build-all    - 跨平台构建"
	@echo "  run          - 构建并运行程序"
	@echo "  test         - 运行测试"
	@echo "  bench        - 运行基准测试"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 运行代码检查"
	@echo "  tidy         - 整理依赖"
	@echo "  clean        - 清理构建文件"
	@echo "  install      - 安装到系统"
	@echo "  uninstall    - 从系统卸载"
	@echo "  setup-grammars - 创建语法目录"
	@echo "  example      - 生成示例项目上下文"
	@echo "  help         - 显示此帮助信息"

# Docker 相关目标
.PHONY: docker-build
docker-build:
	@echo "🐳 构建Docker镜像..."
	docker build -t codecartographer:${VERSION} .

.PHONY: docker-run
docker-run:
	@echo "🐳 运行Docker容器..."
	docker run --rm -v $(PWD):/workspace codecartographer:${VERSION} generate --path /workspace
