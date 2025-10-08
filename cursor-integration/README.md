# code-outline 声明管理工具

基于 code-outline 的快速项目声明和结构了解工具，专为 Cursor 编辑器集成设计。

## 🎯 核心功能

### 1. 获取所有文件声明
- 获取项目中所有文件的声明内容
- 支持缓存机制，提高性能
- 生成完整的项目声明摘要

### 2. 获取指定文件声明
- 获取单个文件的详细声明信息
- 支持相对路径和绝对路径
- 智能缓存管理

### 3. 创建项目声明文件
- 生成完整的项目声明文档
- 包含文件索引和分类
- 支持多种输出格式

### 4. 更新文件声明
- 增量更新指定文件的声明
- 检测文件变化
- 维护声明一致性

## 🚀 快速开始

### 安装工具

```bash
# 安装到 Cursor
python install-spec-kit.py install

# 检查安装状态
python install-spec-kit.py check

# 卸载工具
python install-spec-kit.py uninstall
```

### 基本使用

#### 1. 获取所有文件声明

```bash
# 命令行使用
python declaration-manager-simple.py get-all --path /path/to/project

# 在 Cursor 中使用
# 按 Ctrl+Shift+P，选择 "获取所有文件声明"
```

#### 2. 获取指定文件声明

```bash
# 命令行使用
python declaration-manager-simple.py get-file --path /path/to/project --file src/main.go

# 在 Cursor 中使用
# 右键文件，选择 "获取文件声明"
```

#### 3. 创建项目声明文件

```bash
# 命令行使用
python declaration-manager-simple.py create-project --path /path/to/project

# 在 Cursor 中使用
# 按 Ctrl+Shift+P，选择 "创建项目声明"
```

#### 4. 更新文件声明

```bash
# 命令行使用
python declaration-manager-simple.py update-file --path /path/to/project --file src/main.go

# 在 Cursor 中使用
# 右键文件，选择 "更新文件声明"
```

## 📊 输出格式

### 所有文件声明输出

```json
{
  "timestamp": "2025-01-07 15:30:00",
  "project_path": "/path/to/project",
  "total_files": 25,
  "declarations": {
    "files": {
      "src/main.go": {
        "purpose": "主程序入口",
        "symbols": [
          {
            "prototype": "func main()",
            "purpose": "程序入口点",
            "range": [10, 15]
          }
        ]
      }
    }
  },
  "summary": {
    "total_files": 25,
    "total_symbols": 150,
    "languages": ["Go", "JavaScript"],
    "file_types": {
      ".go": 15,
      ".js": 10
    }
  }
}
```

### 文件声明输出

```json
{
  "timestamp": "2025-01-07 15:30:00",
  "file_path": "/path/to/project/src/main.go",
  "file_name": "main.go",
  "declarations": {
    "files": {
      "src/main.go": {
        "purpose": "主程序入口",
        "symbols": [...]
      }
    }
  },
  "summary": {
    "file_name": "main.go",
    "total_symbols": 5,
    "symbol_types": {
      "functions": 3,
      "variables": 2
    }
  }
}
```

## ⚙️ 配置选项

### 命令行参数

```bash
python declaration-manager.py [action] [options]

Actions:
  get-all          获取所有文件声明
  get-file         获取指定文件声明
  create-project   创建项目声明文件
  update-file      更新文件声明

Options:
  --path PATH      项目路径 (默认: 当前目录)
  --file FILE      指定文件路径
  --output FILE    输出文件
  --no-cache       不使用缓存
  --verbose        详细输出
```

### 缓存配置

工具支持智能缓存机制：

- **缓存位置**: `.declaration_cache/` 目录
- **缓存有效期**: 24小时
- **缓存策略**: 基于文件修改时间
- **清理机制**: 自动清理过期缓存

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

### 1. 智能缓存

```bash
# 使用缓存（默认）
python declaration-manager.py get-all

# 强制刷新缓存
python declaration-manager.py get-all --no-cache
```

### 2. 批量操作

```bash
# 创建项目声明文件
python declaration-manager.py create-project --output my_project.json

# 更新多个文件
for file in src/*.go; do
    python declaration-manager.py update-file --file "$file"
done
```

### 3. 集成到 CI/CD

```yaml
# GitHub Actions 示例
- name: Generate Project Declarations
  run: |
    python declaration-manager.py create-project
    # 上传到存储或发送到 API
```

## 📈 性能优化

### 缓存策略

- 文件级缓存：基于文件修改时间
- 项目级缓存：基于项目结构变化
- 智能失效：自动检测文件变化

### 性能指标

- 小项目（<100文件）：< 5秒
- 中等项目（100-1000文件）：< 30秒
- 大项目（>1000文件）：< 2分钟

## 🐛 故障排除

### 常见问题

**Q: 找不到 contextgen 可执行文件**
```bash
# 检查 contextgen 是否在 PATH 中
where contextgen  # Windows
which contextgen  # Linux/macOS

# 或指定完整路径
export CONTEXTGEN_PATH="/path/to/contextgen"
```

**Q: 缓存文件损坏**
```bash
# 清理缓存
rm -rf .declaration_cache/
python declaration-manager.py get-all --no-cache
```

**Q: 权限问题**
```bash
# 确保有写入权限
chmod +x declaration-manager.py
chmod +x install-spec-kit.py
```

### 调试模式

```bash
# 启用详细输出
python declaration-manager.py get-all --verbose

# 检查安装状态
python install-spec-kit.py check
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 开发环境设置

```bash
# 克隆项目
git clone <repository-url>
cd code-outline/cursor-integration/spec-driven-tools

# 安装依赖
pip install -r requirements.txt

# 运行测试
python -m pytest tests/
```

## 📄 许可证

MIT License - 详见 LICENSE 文件

---

**code-outline 声明管理工具** - 让您快速了解项目结构！ 🗺️✨

