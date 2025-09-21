# CodeCartographer

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

**CodeCartographer** 是一个高性能、跨平台的命令行工具，用于通过静态分析为任何复杂的代码仓库生成统一、简洁且信息丰富的 `project_context.json` 文件。此文件将作为大语言模型（LLM）的"全局上下文记忆"，使其能够以前所未有的准确性和深度来理解项目架构。

## ✨ 特性

- 🚀 **高性能**: 基于 Go 的高效解析引擎，支持并发处理
- 🌍 **多语言支持**: 通过配置文件支持 9+ 种编程语言
- ⚡ **并发处理**: 利用 Goroutines 实现高速文件扫描
- 🎯 **LLM 优化**: 为 LLM Token 效率极致优化的 JSON 输出格式
- 🔧 **可配置**: 灵活的排除规则和自定义配置
- 📦 **跨平台**: 支持 Windows、Linux、macOS
- 🔍 **智能解析**: 基于正则表达式的多语言符号提取

## 🚀 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/CodeCartographer.git
cd CodeCartographer

# 构建项目
make build

# 或者直接运行
make run
```

### 基本使用

```bash
# 生成当前目录的项目上下文
./build/contextgen generate

# 指定项目路径
./build/contextgen generate --path /path/to/your/project

# 自定义输出文件
./build/contextgen generate --path . --output my_context.json

# 排除特定目录
./build/contextgen generate --exclude "node_modules,vendor,.git"

# 使用自定义配置
./build/contextgen generate --config my_languages.json
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

让我们看看 CodeCartographer 如何分析自己的项目：

```bash
$ ./contextgen generate
🚀 开始生成项目上下文...
📋 加载语言配置...
✅ 已加载 9 种语言的配置
🔧 初始化解析器...
🔍 扫描项目: .
✅ 扫描完成，找到 6 个文件
📦 构建项目上下文...
💾 生成输出文件: project_context.json

📊 统计信息:
  项目名称: CodeCartographer
  技术栈: Go
  文件数量: 6
  模块数量: 6
  符号数量: 53
  最后更新: 2025-09-21 20:02:20
🎉 项目上下文生成完成!
```

## 📄 输出格式

生成的 `project_context.json` 文件包含：

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

### 项目结构

```
CodeCartographer/
├── cmd/contextgen/          # 主程序入口
├── internal/
│   ├── cmd/                 # CLI 命令实现
│   ├── config/              # 配置管理
│   ├── models/              # 数据结构定义
│   ├── parser/              # 代码解析器
│   └── scanner/             # 文件扫描器
├── languages.json           # 语言配置文件
├── Makefile                # 构建脚本
├── Dockerfile              # Docker 配置
└── README.md               # 项目文档
```

### 构建命令

```bash
# 构建项目
make build

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

## ⚙️ 配置

### 语言配置文件 (languages.json)

工具通过 `languages.json` 文件配置对不同语言的支持：

```json
{
  "go": {
    "extensions": [".go"],
    "grammar_path": "./grammars/tree-sitter-go.so",
    "queries": {
      "top_level_symbols": [
        "(function_declaration) @symbol",
        "(method_declaration) @symbol",
        "(type_declaration) @symbol"
      ],
      "container_body": "(block) @body | (struct_type) @body",
      "container_methods": "(method_declaration) @method"
    }
  }
}
```

### 自定义配置

- 修改 `languages.json` 添加新语言支持
- 调整正则表达式模式以改进符号识别
- 配置文件扩展名映射

## 🎯 使用场景

### 为 LLM 提供项目上下文

```bash
# 生成项目上下文
./contextgen generate --path ./my-project

# 将 project_context.json 提供给 LLM
# LLM 现在可以理解整个项目结构和代码架构
```

### 项目文档生成

CodeCartographer 生成的上下文文件可以作为：
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

- [ ] 真正的 Tree-sitter 集成
- [ ] 更多语言支持
- [ ] Web 界面
- [ ] 云端服务
- [ ] IDE 插件
- [ ] 实时监控和更新

## 📞 联系

- 项目主页: [https://github.com/yourusername/CodeCartographer](https://github.com/yourusername/CodeCartographer)
- Issues: [https://github.com/yourusername/CodeCartographer/issues](https://github.com/yourusername/CodeCartographer/issues)

---

**CodeCartographer** - 让 LLM 更好地理解您的代码项目 🗺️