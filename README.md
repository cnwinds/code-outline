# CodeCartographer

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

**CodeCartographer** 是一个高性能、跨平台的命令行工具，用于通过静态分析为任何复杂的代码仓库生成统一、简洁且信息丰富的 `project_context.json` 文件。此文件将作为大语言模型（LLM）的"全局上下文记忆"，使其能够以前所未有的准确性和深度来理解项目架构。

## ✨ 特性

- 🚀 **高性能**: 基于 Go 和 Tree-sitter 的高效解析引擎
- 🌍 **多语言支持**: 通过配置文件支持任意编程语言
- ⚡ **并发处理**: 利用 Goroutines 实现高速文件扫描
- 🎯 **LLM 优化**: 为 LLM Token 效率极致优化的 JSON 输出格式
- 🔧 **可配置**: 灵活的排除规则和自定义配置
- 📦 **跨平台**: 支持 Windows、Linux、macOS

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
```

---

### **最终版开发需求文档：通用型项目上下文生成器 (`ContextGen`)**

**1. 项目愿景**

开发一个名为 `ContextGen` 的高性能、跨平台的命令行工具。该工具旨在通过静态分析，为任何复杂的代码仓库生成一个统一、简洁且信息丰富的 `project_context.json` 文件。此文件将作为大语言模型（LLM）的“全局上下文记忆”，使其能够以前所未有的准确性和深度来理解项目架构，从而革命性地提升其在代码生成、需求变更、重构和调试等任务上的表现。

**2. 核心技术栈**

*   **开发语言:** **Go**。利用其卓越的性能、强大的并发能力、简单的跨平台编译和部署。
*   **代码解析框架:** **Tree-sitter**。利用其高效的增量解析能力和丰富的社区语法包，实现对多种编程语言的精确、健壮的解析。

**3. `project_context.json` 输出格式规范 (最终版)**

这是工具的核心产出物，其设计在 **LLM 理解能力** 和 **Token 效率** 之间达到了最佳平衡。

#### a. 顶层结构

```json
{
  "projectName": "...",
  "projectGoal": "TODO: ...",
  "techStack": ["Go", "JavaScript", "..."],
  "lastUpdated": "...",
  "architecture": {
    "overview": "TODO: ...",
    "moduleSummary": {
      "cmd/contextgen": "TODO: ..."
    }
  },
  "files": {
    "path/to/file.go": {
      "purpose": "TODO: ...",
      "symbols": [ /* Symbol Object Array */ ]
    }
  }
}
```

#### b. `Symbol` 对象结构

这是描述代码中一个“符号”（如结构体、函数、常量等）的统一格式。

```go
// Go Struct Definition for a Symbol
type Symbol struct {
    Prototype string   `json:"prototype"`
    Purpose   string   `json:"purpose"`
    Range     []int    `json:"range"`
    Body      string   `json:"body,omitempty"` // 用于类/结构体/接口等容器类型
    Methods   []Symbol `json:"methods,omitempty"` // 用于类/结构体的方法
}
```

*   `prototype`: 符号的完整声明行，从源代码原样复制。
*   `purpose`: 从符号上方或旁边的文档注释中提取的说明。
*   `range`: `[start_line, end_line]`，符号在文件中的行号范围。
*   `body`: **(关键优化)** 对于结构体、类、接口等，此字段包含其内部所有内容的**原始多行字符串**，保留缩进和注释。这极大地节省了 token。
*   `methods`: 对于可以拥有方法的符号（如结构体），此数组包含其所有方法，每个方法也是一个 `Symbol` 对象。

#### c. 完整示例

对于以下 Go 源码 (`database/models.go`):

```go
package database

const DefaultRole = "user"

// User defines the user model.
type User struct {
    ID    int    `json:"id"`
    Email string `json:"email"`
}

// IsAdmin checks user privileges.
func (u *User) IsAdmin() bool {
    return u.Email == "admin@example.com"
}
```

生成的 JSON 部分应为：

```json
"database/models.go": {
  "purpose": "TODO: Describe the purpose of this file.",
  "symbols": [
    {
      "prototype": "const DefaultRole = \"user\"",
      "purpose": "",
      "range": [3, 3]
    },
    {
      "prototype": "type User struct",
      "purpose": "User defines the user model.",
      "range": [6, 9],
      "body": "    ID    int    `json:\"id\"`\n    Email string `json:\"email\"`",
      "methods": [
        {
          "prototype": "func (u *User) IsAdmin() bool",
          "purpose": "IsAdmin checks user privileges.",
          "range": [12, 14]
        }
      ]
    }
  ]
}
```

**4. 功能与技术实现需求**

1.  **可扩展的多语言支持:**
    *   工具必须通过一个外部配置文件 `languages.json` 来管理对不同语言的支持，无需重新编译程序即可添加新语言。
    *   该配置文件定义了语言与文件扩展名的映射、预编译的 Tree-sitter 语法库 (`.so`/`.dll`) 的路径，以及用于提取各种符号的 Tree-sitter 查询。

2.  **基于 Tree-sitter 的核心解析引擎:**
    *   **动态加载语法**: 程序能根据要解析的文件类型，动态加载对应的 Tree-sitter 语法库。
    *   **分层查询**: 解析逻辑应分层。首先，使用查询找到文件中的所有顶级符号。然后，对于容器类型的符号（如 `struct`），在其对应的语法树节点上**递归执行**方法查询和**提取**主体文本，以构建嵌套的 `Symbol` 结构。
    *   **文本提取**: `prototype`, `body` 和 `purpose` 必须直接从源文件文本中精确提取，保留原始格式。

3.  **高性能并发处理:**
    *   必须利用 Go 的 Goroutines 对文件进行并发扫描和解析，以显著加快在大型代码库上的运行速度。

4.  **健壮的命令行接口 (CLI):**
    *   使用 Go 的 `cobra` 或类似库构建。
    *   `contextgen generate --path <project_path>`: 主命令。
    *   `--output <file_path>` (可选): 指定输出文件路径。
    *   `--exclude <dir1,dir2>` (可选): 指定要排除的目录或文件模式。
    *   `--config <config_path>` (可选): 指定 `languages.json` 的路径。

**5. 语言配置文件 (`languages.json`) 规范**

```json
{
  "go": {
    "extensions": [".go"],
    "grammar_path": "./grammars/tree-sitter-go.so",
    "queries": {
      "top_level_symbols": [
        "(function_declaration) @symbol",
        "(method_declaration) @symbol",
        "(type_declaration) @symbol",
        "(const_declaration) @symbol",
        "(var_declaration) @symbol"
      ],
      "container_body": "(block) @body | (struct_type) @body | (interface_type) @body",
      "container_methods": "(method_declaration) @method"
    }
  }
}
```*   `queries` 定义了从语法树中捕获目标节点的规则。

**6. 开发实施计划**

1.  **项目初始化:** 设置 Go 项目，引入 `go-tree-sitter` 和 CLI 库。
2.  **定义数据结构:** 在 Go 中创建与 `project_context.json` 格式完全匹配的 `struct`。
3.  **配置模块:** 实现 `languages.json` 的加载和解析逻辑。
4.  **核心解析器 (`Parser`):** 这是项目的核心。封装 Tree-sitter 的所有交互：加载语法、解析代码、执行查询，并将查询结果转换为我们的 `Symbol` 结构。
5.  **文件处理与并发控制:** 实现文件遍历、过滤逻辑，并使用 Goroutine 池来调度 `Parser` 对文件进行并发处理。
6.  **CLI 实现:** 构建用户友好的命令行接口。
7.  **主程序:** 整合所有模块，编排从参数解析到最终文件生成的完整流程。

这份文档为 `ContextGen` 的开发提供了完整的蓝图。请开始构建这个强大而高效的开发辅助工具。