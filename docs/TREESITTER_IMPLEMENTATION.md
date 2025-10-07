# Tree-sitter 实现指南

本文档提供了在 CodeCartographer 中实现 Tree-sitter 解析器的完整指南。

> **当前状态**: Tree-sitter 功能已集成完成。本指南提供实现细节和使用说明。

---

## 📋 目录

- [1. 概述](#1-概述)
- [2. 技术选型](#2-技术选型)
- [3. 环境准备](#3-环境准备)
- [4. 实现步骤](#4-实现步骤)
- [5. 测试验证](#5-测试验证)
- [6. 故障排除](#6-故障排除)

---

## 1. 概述

### 什么是 Tree-sitter？

Tree-sitter 是一个解析器生成器工具和增量解析库。它可以为任何编程语言构建具体语法树（CST），并在文件编辑时高效地更新这些语法树。

### 为什么使用 Tree-sitter？

相比正则表达式解析器，Tree-sitter 提供：

✅ **更高的准确性**: 真正理解代码的语法结构  
✅ **更好的性能**: 增量解析，只更新变化部分  
✅ **更强的功能**: 支持复杂查询，提取嵌套结构  
✅ **更好的维护性**: 社区维护的语法定义

### 当前状态

- ✅ Tree-sitter 解析器已实现
- ✅ 支持 Go、JavaScript、Python 三种语言
- ✅ 正则表达式解析器作为后备方案
- ✅ 配置结构已完成
- ✅ 接口设计已完成

---

## 2. 技术选型

### 推荐的 Go Binding

我们推荐使用 **`smacker/go-tree-sitter`**：

```bash
go get github.com/smacker/go-tree-sitter
```

**优势**:
- ⭐ 18.5k+ stars
- 📦 预编译的语法包
- 🔄 活跃维护
- 📚 完善的文档
- 🚀 性能优秀

**替代方案**: `tree-sitter/go-tree-sitter`（官方绑定，但需要手动编译语法）

### 支持的语言包

```bash
go get github.com/smacker/go-tree-sitter/golang
go get github.com/smacker/go-tree-sitter/javascript
go get github.com/smacker/go-tree-sitter/typescript/typescript
go get github.com/smacker/go-tree-sitter/python
go get github.com/smacker/go-tree-sitter/java
go get github.com/smacker/go-tree-sitter/c
go get github.com/smacker/go-tree-sitter/cpp
go get github.com/smacker/go-tree-sitter/rust
go get github.com/smacker/go-tree-sitter/csharp
```

---

## 3. 环境准备

### 3.1 系统要求

#### Windows
```powershell
# 安装 MinGW-w64 (推荐使用 MSYS2)
# 下载: https://www.msys2.org/

# 或使用 Chocolatey
choco install mingw

# 验证
gcc --version
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# CentOS/RHEL
sudo yum groupinstall "Development Tools"
```

#### macOS
```bash
# 安装 Xcode Command Line Tools
xcode-select --install
```

### 3.2 启用 CGO

Tree-sitter 需要 CGO 支持：

```bash
# 设置环境变量
export CGO_ENABLED=1  # Linux/macOS
set CGO_ENABLED=1     # Windows CMD
$env:CGO_ENABLED=1    # Windows PowerShell
```

### 3.3 安装依赖

```bash
# 更新 go.mod
go get github.com/smacker/go-tree-sitter
go get github.com/smacker/go-tree-sitter/golang
go get github.com/smacker/go-tree-sitter/javascript
go get github.com/smacker/go-tree-sitter/python
go mod tidy
```

---

## 4. 实现步骤

### 4.1 更新 TreeSitterParser

**文件**: `internal/parser/treesitter_parser.go`

```go
package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	
	"github.com/cnwinds/CodeCartographer/internal/config"
	"github.com/cnwinds/CodeCartographer/internal/models"
)

// TreeSitterParser Tree-sitter 解析器
type TreeSitterParser struct {
	languagesConfig models.LanguagesConfig
	parsers         map[string]*sitter.Parser
}

// NewTreeSitterParser 创建新的 Tree-sitter 解析器
func NewTreeSitterParser(languagesConfig models.LanguagesConfig) (*TreeSitterParser, error) {
	p := &TreeSitterParser{
		languagesConfig: languagesConfig,
		parsers:         make(map[string]*sitter.Parser),
	}

	// 初始化各语言解析器
	if err := p.initParsers(); err != nil {
		return nil, err
	}

	return p, nil
}

// initParsers 初始化语言解析器
func (p *TreeSitterParser) initParsers() error {
	// Go 语言
	goParser := sitter.NewParser()
	goParser.SetLanguage(golang.GetLanguage())
	p.parsers["go"] = goParser

	// JavaScript
	jsParser := sitter.NewParser()
	jsParser.SetLanguage(javascript.GetLanguage())
	p.parsers["javascript"] = jsParser
	p.parsers["typescript"] = jsParser // 暂时共用

	// Python
	pyParser := sitter.NewParser()
	pyParser.SetLanguage(python.GetLanguage())
	p.parsers["python"] = pyParser

	return nil
}

// ParseFile 解析单个文件
func (p *TreeSitterParser) ParseFile(filePath string) (*models.FileInfo, error) {
	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 确定语言
	ext := filepath.Ext(filePath)
	langName, _, found := config.GetLanguageByExtension(p.languagesConfig, ext)
	if !found {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 获取解析器
	parser, ok := p.parsers[langName]
	if !ok {
		return nil, fmt.Errorf("未找到 %s 语言的解析器", langName)
	}

	// 解析
	tree, err := parser.ParseCtx(nil, nil, content)
	if err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}
	defer tree.Close()

	// 提取符号
	symbols := p.extractSymbols(tree.RootNode(), content, langName)

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &models.FileInfo{
		Purpose:      extractFilePurpose(content),
		Symbols:      symbols,
		LastModified: fileInfo.ModTime().Format(time.RFC3339),
		FileSize:     fileInfo.Size(),
	}, nil
}

// extractSymbols 从语法树提取符号
func (p *TreeSitterParser) extractSymbols(node *sitter.Node, content []byte, lang string) []models.Symbol {
	var symbols []models.Symbol

	// 获取查询规则
	queries := p.languagesConfig[lang].Queries.TopLevelSymbols

	for _, queryStr := range queries {
		query, err := sitter.NewQuery([]byte(queryStr), node.Language())
		if err != nil {
			continue
		}
		defer query.Close()

		cursor := sitter.NewQueryCursor()
		defer cursor.Close()

		cursor.Exec(query, node)

		for {
			match, ok := cursor.NextMatch()
			if !ok {
				break
			}

			for _, capture := range match.Captures {
				symbol := p.nodeToSymbol(capture.Node, content)
				symbols = append(symbols, symbol)
			}
		}
	}

	return symbols
}

// nodeToSymbol 将语法树节点转换为符号
func (p *TreeSitterParser) nodeToSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	return models.Symbol{
		Prototype: string(content[node.StartByte():node.EndByte()]),
		Purpose:   "", // TODO: 从注释提取
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}
}

// extractFilePurpose 提取文件用途
func extractFilePurpose(content []byte) string {
	// TODO: 实现注释提取逻辑
	return "TODO: Describe the purpose of this file."
}
```

### 4.2 更新构建配置

**Makefile**:

```makefile
# 添加 CGO 标志
build:
	@echo "🔨 构建 CodeCartographer..."
	@mkdir -p ${BUILD_DIR}
	CGO_ENABLED=1 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PATH}
	@echo "✅ 构建完成: ${BUILD_DIR}/${BINARY_NAME}"

# 添加不使用 Tree-sitter 的构建选项
build-simple:
	@echo "🔨 构建 CodeCartographer (无 Tree-sitter)..."
	@mkdir -p ${BUILD_DIR}
	CGO_ENABLED=0 go build ${LDFLAGS} -tags simple -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PATH}
	@echo "✅ 构建完成: ${BUILD_DIR}/${BINARY_NAME}"
```

**Dockerfile**:

```dockerfile
# 使用支持 CGO 的构建环境
FROM golang:1.21-alpine AS builder

# 安装编译工具
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 启用 CGO 构建
RUN CGO_ENABLED=1 GOOS=linux go build -a -o contextgen ./cmd/contextgen

# 运行时镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /root/

COPY --from=builder /app/contextgen .
COPY --from=builder /app/languages.json .

USER nobody

ENTRYPOINT ["./contextgen"]
CMD ["--help"]
```

### 4.3 添加构建标签（可选）

支持条件编译，同时保留简单解析器：

**treesitter_parser.go**:
```go
//go:build !simple
// +build !simple

package parser
// ... Tree-sitter 实现
```

**simple_parser_fallback.go**:
```go
//go:build simple
// +build simple

package parser
// ... 简单解析器作为后备
```

---

## 5. 测试验证

### 5.1 单元测试

**文件**: `internal/parser/treesitter_parser_test.go`

```go
package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cnwinds/CodeCartographer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTreeSitterParser(t *testing.T) {
	config := getTestConfig()
	
	parser, err := NewTreeSitterParser(config)
	require.NoError(t, err)
	assert.NotNil(t, parser)
}

func TestTreeSitterParseGoFile(t *testing.T) {
	parser, err := NewTreeSitterParser(getTestConfig())
	require.NoError(t, err)

	tmpFile := createTempFile(t, "test.go", goTestCode)
	defer os.Remove(tmpFile)

	result, err := parser.ParseFile(tmpFile)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Symbols), 1)
}
```

### 5.2 集成测试

```bash
# 运行测试
go test ./internal/parser -v

# 测试覆盖率
go test ./internal/parser -cover

# 基准测试
go test ./internal/parser -bench=.
```

### 5.3 手动测试

```bash
# 构建
make build

# 测试单个文件
./build/contextgen generate --path ./testdata/example.go --output test_output.json

# 对比结果
diff test_output.json expected_output.json
```

---

## 6. 故障排除

### 问题 1: CGO 编译错误

**错误信息**:
```
cgo: C compiler "gcc" not found
```

**解决方案**:
- Windows: 安装 MinGW-w64 或 MSYS2
- Linux: `sudo apt-get install build-essential`
- macOS: `xcode-select --install`

### 问题 2: 找不到语法库

**错误信息**:
```
undefined: golang.GetLanguage
```

**解决方案**:
```bash
go get github.com/smacker/go-tree-sitter/golang
go mod tidy
```

### 问题 3: 链接错误

**错误信息**:
```
undefined reference to `ts_parser_new'
```

**解决方案**:
```bash
# 清理并重新构建
go clean -cache
CGO_ENABLED=1 go build -v ./cmd/contextgen
```

### 问题 4: 跨平台编译失败

**问题**: CGO 不支持交叉编译

**解决方案**:
- 方案 1: 在目标平台上构建
- 方案 2: 使用 Docker 多平台构建
- 方案 3: 提供预编译二进制文件

---

## 7. 性能优化

### 7.1 缓存解析器

```go
var parserPool = sync.Pool{
    New: func() interface{} {
        parser := sitter.NewParser()
        parser.SetLanguage(golang.GetLanguage())
        return parser
    },
}
```

### 7.2 并发解析

```go
func (p *TreeSitterParser) ParseFiles(files []string) ([]models.FileInfo, error) {
    results := make([]models.FileInfo, len(files))
    var wg sync.WaitGroup

    for i, file := range files {
        wg.Add(1)
        go func(idx int, path string) {
            defer wg.Done()
            result, err := p.ParseFile(path)
            if err == nil {
                results[idx] = *result
            }
        }(i, file)
    }

    wg.Wait()
    return results, nil
}
```

### 7.3 增量解析

利用 Tree-sitter 的增量解析能力：

```go
// 保存旧的语法树
oldTree := parser.Parse(nil, oldContent)

// 增量解析
newTree := parser.Parse(oldTree, newContent)
```

---

## 8. 参考资源

### 官方文档
- [Tree-sitter 官网](https://tree-sitter.github.io/)
- [Tree-sitter 文档](https://tree-sitter.github.io/tree-sitter/)
- [go-tree-sitter GitHub](https://github.com/smacker/go-tree-sitter)

### 语法仓库
- [tree-sitter-go](https://github.com/tree-sitter/tree-sitter-go)
- [tree-sitter-javascript](https://github.com/tree-sitter/tree-sitter-javascript)
- [tree-sitter-python](https://github.com/tree-sitter/tree-sitter-python)

### 示例项目
- [github/semantic](https://github.com/github/semantic)
- [sourcegraph](https://github.com/sourcegraph/sourcegraph)

---

## 9. 开发路线图

### Phase 1: 基础实现（已完成）
- [x] 环境准备
- [x] Go 语言解析
- [x] JavaScript 解析
- [x] Python 解析
- [x] 基础测试

### Phase 2: 功能增强（2-3 周）
- [ ] 更多语言支持
- [ ] 注释提取
- [ ] 类型信息提取
- [ ] 嵌套结构支持

### Phase 3: 性能优化（1 周）
- [ ] 解析器池
- [ ] 并发优化
- [ ] 增量解析
- [ ] 内存优化

### Phase 4: 生产就绪（1 周）
- [ ] 完整测试覆盖
- [ ] 文档完善
- [ ] 错误处理
- [ ] 发布 v1.0

---

## 10. 常见问题

### Q: 是否必须使用 Tree-sitter？
A: 不是。可以继续使用简单解析器（`--treesitter=false`），但 Tree-sitter 提供更高的准确性。

### Q: 如何在不同平台上构建？
A: 需要在目标平台上构建，或使用 Docker 多平台构建。CGO 不支持交叉编译。

### Q: 性能如何？
A: Tree-sitter 非常快，通常比正则表达式解析器更快，特别是对于大文件和增量更新。

### Q: 如何添加新语言？
A: 
1. 安装对应的 go-tree-sitter 包
2. 在 `initParsers()` 中注册
3. 更新 `languages.json` 配置
4. 添加测试用例

---

**实现完成后，请更新 README.md 移除"开发中"标注！**
