# code-outline

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

**code-outline** 是一个高性能、跨平台的命令行工具，用于通过静态分析为任何复杂的代码仓库生成统一、简洁且信息丰富的 `code-outline.json` 文件。此文件将作为大语言模型（LLM）的"全局上下文记忆"，使其能够以前所未有的准确性和深度来理解项目架构。

## ✨ 特性

- 🚀 **高性能**: 基于 Go 的高效解析引擎，支持并发处理
- 🌍 **多语言支持**: 内置支持 9+ 种编程语言
- ⚡ **并发处理**: 利用 Goroutines 实现高速文件扫描
- 🎯 **LLM 优化**: 为 LLM Token 效率极致优化的 JSON 输出格式
- 🔧 **可配置**: 灵活的排除规则和自定义配置
- 📦 **跨平台**: 支持 Windows、Linux、macOS
- 🔍 **智能解析**: 基于 Tree-sitter 的高精度语法解析，支持复杂嵌套结构

## 🚀 快速开始

### 安装

#### 方式一：下载预编译版本（推荐）

从 [Releases](https://github.com/yourusername/code-outline/releases) 页面下载适合您系统的预编译版本：

- **Windows**: `code-outline-windows-amd64.exe`
- **Linux**: `code-outline-linux-amd64`
- **macOS**: `code-outline-darwin-amd64`

#### 方式二：从源码构建

```bash
# 克隆仓库
git clone https://github.com/yourusername/code-outline.git
cd code-outline

# 构建项目（自动检测平台）
make build

# Windows 专用构建（64 位架构）
make build-windows

# 跨平台构建（所有平台）
make build-all

# 或者直接运行
make run
```

### 基本使用

```bash
# 生成当前目录的项目上下文
./build/code-outline generate

# 指定项目路径
./build/code-outline generate --path /path/to/your/project

# 自定义输出文件
./build/code-outline generate --path . --output my_context.json

# 排除特定目录
./build/code-outline generate --exclude "node_modules,vendor,.git"

# 增量更新项目上下文
./build/code-outline update

# 更新指定文件
./build/code-outline update --files "main.go,config.go"

# 更新指定目录
./build/code-outline update --dirs "internal/,cmd/"

# 同时更新指定文件和目录
./build/code-outline update --files "main.go" --dirs "internal/"

# 查询所有文件和方法定义
./build/code-outline query

# 查询指定文件的数据
./build/code-outline query --files "main.go,config.go"

# 查询指定目录的数据
./build/code-outline query --dirs "internal/,cmd/"

# 保存查询结果到文件
./build/code-outline query --files "main.go" --output data.json
```

## 📋 支持的语言

当前支持的编程语言：

| 语言 | 扩展名 | 符号类型 |
|------|--------|----------|
| Go | `.go` | 函数、方法、结构体、常量、变量 |
| JavaScript | `.js`, `.jsx` | 函数、类、箭头函数、声明 |
| TypeScript | `.ts`, `.tsx` | 函数、类、接口、类型别名 |
| Python | `.py` | 函数、类、赋值 |
| Java | `.java` | 方法、类、接口、字段 |
| C# | `.cs` | 方法、类、接口、结构体、属性 |
| Rust | `.rs` | 函数、结构体、枚举、特征、实现 |
| C++ | `.cpp`, `.cc`, `.cxx`, `.hpp` | 函数、类、结构体、命名空间 |
| C | `.c`, `.h` | 函数、结构体、枚举 |

## 🎯 演示

让我们看看 code-outline 如何分析自己的项目：

```bash
$ ./code-outline generate
🚀 开始生成项目上下文...
📋 加载语言配置...
✅ 已加载 9 种语言的配置
🔧 初始化解析器...
🔍 扫描项目: .
✅ 扫描完成，找到 6 个文件
📦 构建项目上下文...
💾 生成输出文件: code-outline.json

📊 统计信息:
  项目名称: code-outline
  技术栈: Go
  文件数量: 6
  模块数量: 6
  符号数量: 53
  最后更新: 2025-09-21 20:02:20
🎉 项目上下文生成完成!
```

## 📄 输出格式

生成的 `code-outline.json` 文件包含：

```json
{
  "projectName": "项目名称",
  "projectGoal": "项目目标描述", 
  "techStack": ["Go", "JavaScript"],
  "lastUpdated": "2025-09-21T20:02:20Z",
  "architecture": {
    "overview": "架构概述",
    "moduleSummary": {
      "module_path": "模块描述"
    }
  },
  "files": {
    "path/to/file.go": {
      "purpose": "文件用途",
      "symbols": [
        {
          "prototype": "func Example() error",
          "purpose": "函数说明",
          "range": [10, 15],
          "body": "函数体内容（适用于结构体等）",
          "methods": []
        }
      ]
    }
  }
}
```

## 🛠️ 开发

### 环境要求

**Tree-sitter 解析器需要 C 编译器支持：**

- **Windows**: 安装 [MSYS2](https://www.msys2.org/) 和 MinGW-w64
- **Linux**: 安装 `build-essential` 包
- **macOS**: 安装 Xcode Command Line Tools

详细安装指南请参考：[Windows CGO 环境安装文档](docs/WINDOWS_CGO_SETUP.md)

#### Windows 环境 GCC 安装

**方法一：使用 Chocolatey（推荐）**
```bash
# 安装 Chocolatey（如果未安装）
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 安装 MinGW
choco install mingw -y

# 验证安装
gcc --version
```

**方法二：使用 MSYS2**
```bash
# 1. 下载并安装 MSYS2: https://www.msys2.org/
# 2. 打开 MSYS2 终端，运行：
pacman -S mingw-w64-x86_64-gcc
pacman -S mingw-w64-x86_64-pkg-config

# 3. 将 MSYS2 的 bin 目录添加到 PATH
# 通常路径为: C:\msys64\mingw64\bin
```

**方法三：使用 TDM-GCC**
```bash
# 1. 下载 TDM-GCC: https://jmeubank.github.io/tdm-gcc/
# 2. 安装时选择 "Add to PATH"
# 3. 重启命令行验证
gcc --version
```

#### Linux 环境 GCC 安装

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install -y build-essential
```

**CentOS/RHEL:**
```bash
sudo yum groupinstall "Development Tools"
# 或者
sudo dnf groupinstall "Development Tools"
```

**Arch Linux:**
```bash
sudo pacman -S base-devel
```

#### macOS 环境 GCC 安装

```bash
# 安装 Xcode Command Line Tools
xcode-select --install

# 或者使用 Homebrew
brew install gcc
```

#### 验证 CGO 环境

```bash
# 设置环境变量
export CGO_ENABLED=1

# 验证 Go 可以找到 C 编译器
go env CGO_ENABLED
go env CC
```

#### 代码质量检查

**安装 staticcheck:**

```bash
# 使用 go install（推荐）
go install honnef.co/go/tools/cmd/staticcheck@latest

# 或者使用包管理器
# Windows (Chocolatey)
choco install staticcheck

# macOS (Homebrew)
brew install staticcheck

# Ubuntu/Debian
sudo apt-get install staticcheck
```

**运行代码质量检查:**

```bash
# 运行所有检查
staticcheck ./...

# 快速检查（忽略包注释）
staticcheck -checks="all,-ST1000" ./...

# 检查未使用代码
staticcheck -checks=U1000 ./...

# 检查性能问题
staticcheck -checks="S1000,S1001,S1002,S1003,S1004,S1005,S1006,S1007,S1008,S1009,S1010,S1011,S1012,S1016,S1017,S1018,S1019,S1020,S1021,S1023,S1024,S1025,S1028,S1029,S1030,S1031,S1032,S1033,S1034,S1035,S1036,S1037,S1038,S1039,S1040" ./...

# 检查特定目录
staticcheck ./internal/...

# 生成 JSON 报告
staticcheck -f json ./... > lint-report.json
```

**验证安装和运行：**

```bash
# 检查 staticcheck 版本
staticcheck --version

# 运行代码检查
staticcheck ./internal/config ./internal/scanner
```

**使用 Makefile 命令：**

```bash
# 快速检查
make lint

# 完整检查
make lint-full

# 检查未使用代码
make lint-unused

# 检查性能问题
make lint-performance
```

### 项目结构

```
code-outline/
├── cmd/code-outline/          # 主程序入口
├── internal/
│   ├── cmd/                 # CLI 命令实现
│   ├── config/              # 配置管理
│   ├── models/              # 数据结构定义
│   ├── parser/              # 代码解析器
│   └── scanner/             # 文件扫描器
├── Makefile                # 构建脚本
├── Dockerfile              # Docker 配置
└── README.md               # 项目文档
```

### 构建命令

```bash
# 构建项目（自动检测平台）
make build

# Windows 专用构建（64 位架构）
make build-windows

# 跨平台构建
make build-all

# 运行测试
make test

# 代码格式化
make fmt

# 清理构建文件
make clean

# 生成示例
make example
```

### Docker 使用

```bash
# 构建镜像
make docker-build

# 使用 Docker 运行
make docker-run
```

## 🔄 更新模式

code-outline 支持增量更新模式，可以只更新指定的文件或目录，大大提高更新效率：

### 基本更新命令

```bash
# 检测所有文件变更并更新
./build/code-outline update

# 指定项目路径和输出文件
./build/code-outline update --path /path/to/project --output my_context.json
```

### 指定文件更新

```bash
# 更新单个文件
./build/code-outline update --files "main.go"

# 更新多个文件
./build/code-outline update --files "main.go,config.go,utils.go"

# 更新指定路径的文件
./build/code-outline update --files "cmd/main.go,internal/config/config.go"
```

### 指定目录更新

```bash
# 更新单个目录
./build/code-outline update --dirs "internal/"

# 更新多个目录
./build/code-outline update --dirs "internal/,cmd/,pkg/"

# 更新子目录
./build/code-outline update --dirs "internal/parser/,internal/scanner/"
```

### 混合更新模式

```bash
# 同时更新指定文件和目录
./build/code-outline update --files "main.go" --dirs "internal/"

# 结合排除规则
./build/code-outline update --files "main.go" --exclude "*.test.go"
```

### 更新模式的优势

- **高效**: 只解析指定的文件，避免全量扫描
- **精确**: 针对特定文件或目录进行更新
- **快速**: 大幅减少更新时间和资源消耗
- **灵活**: 支持文件和目录的任意组合

## 🔍 查询模式

code-outline 支持查询模式，可以查询指定文件或目录中的所有文件和方法定义，返回结构化的JSON数据：

### 基本查询命令

```bash
# 查询所有文件和方法定义
./build/code-outline query

# 指定项目路径
./build/code-outline query --path /path/to/project

# 输出到文件
./build/code-outline query --output data.json
```

### 指定文件查询

```bash
# 查询单个文件的数据
./build/code-outline query --files "main.go"

# 查询多个文件的数据
./build/code-outline query --files "main.go,config.go,utils.go"

# 查询指定路径的文件数据
./build/code-outline query --files "cmd/main.go,internal/config/config.go"
```

### 指定目录查询

```bash
# 查询单个目录的数据
./build/code-outline query --dirs "internal/"

# 查询多个目录的数据
./build/code-outline query --dirs "internal/,cmd/,pkg/"

# 查询子目录的数据
./build/code-outline query --dirs "internal/parser/,internal/scanner/"
```

### 混合查询模式

```bash
# 同时查询指定文件和目录的数据
./build/code-outline query --files "main.go" --dirs "internal/"

# 结合排除规则
./build/code-outline query --files "main.go" --exclude "*.test.go"

# 输出到标准输出（不指定output参数）
./build/code-outline query --files "main.go"
```

### 查询模式的优势

- **结构化**: 返回标准化的JSON格式数据
- **精确**: 可以指定特定的文件或目录
- **完整**: 包含所有文件和方法定义信息
- **灵活**: 支持多种输出方式（文件或标准输出）

### 输出格式

查询模式返回的JSON格式包含：

```json
{
  "files": {
    "path/to/file.go": {
      "purpose": "文件用途描述",
      "symbols": [
        {
          "prototype": "func Example() error",
          "purpose": "函数说明",
          "range": [10, 15],
          "body": "函数体内容",
          "methods": []
        }
      ],
      "lastModified": "2025-01-01T12:00:00Z",
      "fileSize": 1024
    }
  },
  "stats": {
    "totalFiles": 10,
    "totalSymbols": 50,
    "languages": ["Go", "JavaScript"]
  }
}
```

## 🎯 使用场景

### 为 LLM 提供项目上下文

```bash
# 生成项目上下文
./code-outline generate --path ./my-project

# 将 code-outline.json 提供给 LLM
# LLM 现在可以理解整个项目结构和代码架构
```

### 项目文档生成

code-outline 生成的上下文文件可以作为：
- 项目架构文档的基础
- 新成员入职的参考资料
- 代码审查的辅助工具
- 重构规划的依据

### 代码分析

- 快速了解大型项目的结构
- 识别关键模块和依赖关系
- 分析代码质量和复杂度

## 📊 性能

- **并发处理**: 多 Goroutine 并行扫描文件
- **内存效率**: 流式处理大型文件
- **速度优化**: 智能文件过滤和缓存

典型性能指标：
- 1000 个文件的项目：~2-5 秒
- 10000 个文件的项目：~10-30 秒

## 🤝 贡献

欢迎贡献代码！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详细信息。

### 开发流程

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 License

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🔮 未来计划

- [x] Tree-sitter 集成（已完成）
- [ ] 更多语言支持
- [ ] 注释提取优化
- [ ] Web 界面
- [ ] 云端服务
- [ ] IDE 插件
- [ ] 实时监控和更新

## 🚀 自动发布

### 创建 Release

项目支持自动发布功能。当您推送以 `v` 开头的 tag 时，会自动触发 release 工作流：

```bash
# 创建带消息的 tag
git tag -a v1.0.0 -m "发布版本 1.0.0

## 新功能
- 添加了新的解析器
- 优化了性能
- 修复了已知问题"

# 推送 tag 触发自动发布
git push origin v1.0.0
```

### 测试 Release 工作流

使用提供的测试脚本：

```bash
# Linux/macOS
./scripts/test-release.sh

# Windows PowerShell
.\scripts\test-release.ps1
```

### Release 内容

自动生成的 release 将包含：

- 📦 **Windows**: `code-outline-windows-amd64.exe`
- 🐧 **Linux**: `code-outline-linux-amd64`  
- 🍎 **macOS**: `code-outline-darwin-amd64`

每个 release 都会自动生成下载链接和使用说明。

## 🛠️ 故障排除

### 常见问题

**Q: Tree-sitter 解析器无法使用？**
A: 请确保已安装 C 编译器。Windows 用户请参考 [Windows CGO 环境安装文档](docs/WINDOWS_CGO_SETUP.md)。如果仍有问题，可以使用 Docker 构建方式。

**Q: Windows 下编译时出现链接器错误（如 "cannot find -lmingwex"）？**
A: 这通常是因为 Go 使用了 32 位架构。解决方法：
```bash
# 设置 64 位架构
$env:GOARCH="amd64"
$env:CGO_ENABLED=1
$env:CC="gcc"

# 然后重新构建
go build -o build/code-outline.exe ./cmd/code-outline
```

**Q: 扫描大项目时内存占用过高？**
A: 这是已知问题，建议使用 `--exclude` 参数排除不必要的目录，如 `node_modules`、`vendor` 等。

**Q: 生成的 JSON 文件过大？**
A: 可以调整排除规则，或考虑分模块生成上下文文件。

### 性能优化建议

1. 使用 `--exclude` 排除大型依赖目录
2. 对于大型项目，考虑分模块处理
3. 定期清理生成的上下文文件

---

**code-outline** - 让 LLM 更好地理解您的代码项目 🗺️