# code-outline Makefile
# 
# 常用命令:
#   make build        - 构建项目
#   make lint         - 运行代码检查
#   make lint-quick   - 快速代码检查（忽略包注释）
#   make lint-unused  - 检查未使用的代码
#   make test         - 运行测试
#   make clean        - 清理构建文件

# 变量定义
BINARY_NAME=code-outline
MAIN_PATH=./cmd/code-outline
BUILD_DIR=./build
VERSION=v1.0.0
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# 默认目标
.PHONY: all
all: clean build

# 构建二进制文件
.PHONY: build
build:
	@echo "🔨 构建 code-outline..."
	@mkdir -p ${BUILD_DIR}
	@if [ "$(OS)" = "Windows_NT" ]; then \
		echo "🪟 检测到 Windows 环境，设置 64 位架构..."; \
		set GOARCH=amd64 && set CGO_ENABLED=1 && go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}.exe ${MAIN_PATH}; \
	else \
		CGO_ENABLED=1 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PATH}; \
	fi
	@echo "✅ 构建完成: ${BUILD_DIR}/${BINARY_NAME}"

# 跨平台构建
.PHONY: build-all
build-all:
	@echo "🔨 开始跨平台构建..."
	@mkdir -p ${BUILD_DIR}
	
	# Windows
	@echo "构建 Windows 版本..."
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ${MAIN_PATH}
	
	# Linux
	@echo "构建 Linux 版本..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ${MAIN_PATH}
	
	# macOS
	@echo "构建 macOS 版本..."
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 ${MAIN_PATH}
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 ${MAIN_PATH}
	
	@echo "✅ 跨平台构建完成"

# Windows 专用构建（解决 64 位架构问题）
.PHONY: build-windows
build-windows:
	@echo "🪟 构建 Windows 版本（64 位）..."
	@mkdir -p ${BUILD_DIR}
	@echo "设置 64 位架构和 CGO..."
	set GOARCH=amd64 && set CGO_ENABLED=1 && go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}.exe ${MAIN_PATH}
	@echo "✅ Windows 构建完成: ${BUILD_DIR}/${BINARY_NAME}.exe"


# 运行程序
.PHONY: run
run: build
	@echo "🚀 运行 code-outline..."
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

# 代码检查（默认：快速检查）
.PHONY: lint
lint:
	@echo "🔍 运行快速代码检查..."
	staticcheck -checks="all,-ST1000" ./...

# 代码检查（完整检查）
.PHONY: lint-full
lint-full:
	@echo "🔍 运行完整代码检查..."
	staticcheck ./...

# 代码检查（忽略包注释）
.PHONY: lint-quick
lint-quick:
	@echo "🔍 运行快速代码检查（忽略包注释）..."
	staticcheck -checks="all,-ST1000" ./...

# 代码检查（仅未使用代码）
.PHONY: lint-unused
lint-unused:
	@echo "🔍 检查未使用的代码..."
	staticcheck -checks=U1000 ./...

# 代码检查（仅性能问题）
.PHONY: lint-performance
lint-performance:
	@echo "🔍 检查性能问题..."
	staticcheck -checks=S1000,S1001,S1002,S1003,S1004,S1005,S1006,S1007,S1008,S1009,S1010,S1011,S1012,S1016,S1017,S1018,S1019,S1020,S1021,S1023,S1024,S1025,S1028,S1029,S1030,S1031,S1032,S1033,S1034,S1035,S1036,S1037,S1038,S1039,S1040 ./...

# 代码检查（仅错误和警告）
.PHONY: lint-errors
lint-errors:
	@echo "🔍 检查错误和警告..."
	staticcheck -checks="SA,ST" ./...

# 代码检查（特定目录）
.PHONY: lint-internal
lint-internal:
	@echo "🔍 检查 internal 目录..."
	staticcheck -checks="all,-ST1000" ./internal/...

# 代码检查（生成报告）
.PHONY: lint-report
lint-report:
	@echo "📊 生成代码检查报告..."
	staticcheck -checks="all,-ST1000" -f json ./... > lint-report.json
	@echo "✅ 报告已生成: lint-report.json"

# 安装 staticcheck
.PHONY: install-lint
install-lint:
	@echo "📦 安装 staticcheck..."
	@if command -v staticcheck >/dev/null 2>&1; then \
		echo "✅ staticcheck 已安装"; \
	else \
		echo "正在安装 staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
		echo "✅ staticcheck 安装完成"; \
	fi


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
	rm -f code-outline.json

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
	@echo "code-outline Makefile 命令:"
	@echo "  build        - 构建二进制文件 (启用 CGO，自动检测平台)"
	@echo "  build-windows- 构建 Windows 版本 (64 位架构)"
	@echo "  build-all    - 跨平台构建"
	@echo "  run          - 构建并运行程序"
	@echo "  test         - 运行测试"
	@echo "  bench        - 运行基准测试"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 运行代码检查"
	@echo "  lint-verbose - 运行代码检查（详细输出）"
	@echo "  lint-fix     - 运行代码检查并自动修复"
	@echo "  lint-internal- 检查 internal 目录"
	@echo "  lint-report  - 生成代码检查报告"
	@echo "  install-lint - 安装 staticcheck"
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
