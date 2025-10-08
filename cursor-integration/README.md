# code-outline Cursor 集成工具

基于 code-outline 的智能项目上下文生成工具，专为 Cursor 编辑器集成设计。通过 Tree-sitter 解析器提供高精度的代码结构分析，为 LLM 提供完整的项目理解能力。

## 🎯 核心功能

### 1. 生成项目上下文
- 首次使用或需要完整项目上下文时使用
- 生成整个项目的结构、文件和方法定义
- 创建 `code-outline.json` 文件
- 适用于：新项目分析、完整重构、架构理解

### 2. 更新项目上下文
- 当文件发生变化后使用
- 增量更新 `code-outline.json`，只重新解析已修改的文件
- 适用于：文件修改后的上下文更新、保持项目状态同步

### 3. 查询特定文件
- 需要了解特定文件结构时使用
- 查询指定文件的方法、类、函数定义
- 支持多文件查询，一次返回多个文件内容
- 适用于：代码审查、特定模块分析、函数调用关系理解

### 4. 查询目录结构
- 需要了解特定目录结构时使用
- 查询指定目录下所有文件的方法和类定义
- 适用于：模块分析、目录结构理解、相关文件批量分析

## 🚀 快速开始

### 安装工具

```bash
# 进入工具目录
cd cursor-integration

# 安装到 Cursor
python install-code-outline.py install

# 检查安装状态
python install-code-outline.py check

# 卸载工具
python install-code-outline.py uninstall
```

### 重启 Cursor

安装完成后，重启 Cursor 编辑器以加载新工具。

### 使用工具

#### 方法一：使用 Tasks（推荐）

1. 打开您的项目
2. 按 `Ctrl+Shift+P` (Windows/Linux) 或 `Cmd+Shift+P` (macOS)
3. 输入 "Tasks: Run Task" 或直接输入任务名称：
   - **生成项目上下文** - 分析整个项目
   - **更新项目上下文** - 更新已修改的文件
   - **查询特定文件** - 分析单个文件
   - **查询目录结构** - 分析目录结构

#### 方法二：使用快捷键

- `Ctrl+Shift+G` - 生成项目上下文
- `Ctrl+Shift+U` - 更新项目上下文
- `Ctrl+Shift+Q` - 查询特定文件
- `Ctrl+Shift+D` - 查询目录结构

#### 方法三：命令行使用

```bash
# 生成项目上下文
./build/code-outline.exe generate --path .

# 更新项目上下文
./build/code-outline.exe update --path .

# 查询特定文件
./build/code-outline.exe query --files "main.go,config.go" --path . --ouput query_result.json

# 查询目录结构
./build/code-outline.exe query --dirs "src/,internal/" --path . --ouput directory_query_result.json
```

## 📊 输出文件说明

工具会生成以下文件：

- `code-outline.json` - 完整的项目上下文文件
- `query_result.json` - 特定文件查询结果
- `directory_query_result.json` - 目录结构查询结果

## 🎯 使用场景

### 1. 新项目分析
在 Cursor 中按 `Ctrl+Shift+P`，输入 "Tasks: Run Task"，选择 "生成项目上下文"

### 2. 代码审查
在 Cursor 中按 `Ctrl+Shift+P`，输入 "Tasks: Run Task"，选择 "查询特定文件"

### 3. 项目文档生成
在 Cursor 中按 `Ctrl+Shift+P`，输入 "Tasks: Run Task"，选择 "生成项目上下文"

### 4. 增量更新
在 Cursor 中按 `Ctrl+Shift+P`，输入 "Tasks: Run Task"，选择 "更新项目上下文"

## ⚙️ 配置选项

### 命令行参数

```bash
./build/code-outline.exe [command] [options]

Commands:
  generate          生成项目上下文
  update            更新项目上下文
  query             查询文件或目录
  version           显示版本信息

Options:
  --path PATH       项目路径 (默认: 当前目录)
  --output FILE     输出文件路径
  --exclude PATTERN 排除目录或文件模式
  --files FILES     指定要查询的文件（逗号分隔）
  --dirs DIRS       指定要查询的目录（逗号分隔）
```

### 支持的语言

- Go (.go)
- JavaScript (.js, .jsx)
- TypeScript (.ts, .tsx)
- Python (.py)
- Java (.java)
- C# (.cs)
- Rust (.rs)
- C/C++ (.c, .cpp, .h, .hpp)

## 🔧 高级功能

### 1. 智能解析

- 基于 Tree-sitter 的高精度语法解析
- 支持复杂嵌套结构
- 自动识别函数、类、方法定义
- 提取注释和文档字符串

### 2. 增量更新

```bash
# 更新指定文件
./build/code-outline.exe update --files "main.go,config.go"

# 更新指定目录
./build/code-outline.exe update --dirs "src/,internal/"
```

### 3. 排除模式

```bash
# 排除特定目录
./build/code-outline.exe generate --exclude "node_modules,vendor,.git"
```

## 📈 性能优化

### 性能指标

- 小项目（<100文件）：< 5秒
- 中等项目（100-1000文件）：< 30秒
- 大项目（>1000文件）：< 2分钟

### 优化策略

- 并发文件处理
- 智能缓存机制
- 增量更新支持
- 内存优化

## 🐛 故障排除

### 常见问题

**Q: 找不到 code-outline 可执行文件**
```bash
# 确保 code-outline 已构建
cd ../../  # 回到项目根目录
make build
```

**Q: 权限问题**
```bash
# 给脚本执行权限
chmod +x install-code-outline.py
```

**Q: Cursor 中找不到工具**
```bash
# 重新安装
python install-code-outline.py uninstall
python install-code-outline.py install
# 然后重启 Cursor
```

**Q: Tree-sitter 解析器初始化失败**
```bash
# 检查 CGO 是否启用
go env CGO_ENABLED

# 如果为 false，启用 CGO
export CGO_ENABLED=1  # Linux/macOS
set CGO_ENABLED=1     # Windows CMD
$env:CGO_ENABLED=1    # Windows PowerShell
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 开发环境设置

```bash
# 克隆项目
git clone <repository-url>
cd code-outline

# 构建项目
make build

# 运行测试
make test
```

## 📄 许可证

MIT License - 详见 LICENSE 文件

---

**code-outline Cursor 集成工具** - 让您快速了解项目结构！ 🗺️✨