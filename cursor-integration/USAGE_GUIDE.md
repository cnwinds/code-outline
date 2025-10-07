# CodeCartographer 声明管理工具使用指南

## 🎯 工具概述

这是一个基于 CodeCartographer 的快速项目声明和结构了解工具，专为 Cursor 编辑器集成设计。工具提供4个核心功能：

1. **获取所有文件声明** - 分析整个项目的所有文件声明
2. **获取指定文件声明** - 分析单个文件的声明内容  
3. **创建项目声明文件** - 生成完整的项目声明文档
4. **更新文件声明** - 增量更新指定文件的声明

## 🚀 快速开始

### 1. 安装工具

```bash
# 进入工具目录
cd cursor-integration/spec-driven-tools

# 安装到 Cursor
python install-spec-kit.py install

# 检查安装状态
python install-spec-kit.py check
```

### 2. 重启 Cursor

安装完成后，重启 Cursor 编辑器以加载新工具。

### 3. 使用工具

#### 方法一：在 Cursor 中使用

1. 打开您的项目
2. 按 `Ctrl+Shift+P` (Windows/Linux) 或 `Cmd+Shift+P` (macOS)
3. 输入 "External Tools" 或 "声明管理"
4. 选择相应的工具：
   - **获取所有文件声明** - 分析整个项目
   - **获取指定文件声明** - 分析单个文件
   - **创建项目声明** - 生成项目文档
   - **更新文件声明** - 更新文件信息

#### 方法二：命令行使用

```bash
# 获取所有文件声明
python declaration-manager-simple.py get-all --path /path/to/your/project

# 获取指定文件声明
python declaration-manager-simple.py get-file --path /path/to/your/project --file src/main.go

# 创建项目声明文件
python declaration-manager-simple.py create-project --path /path/to/your/project

# 更新文件声明
python declaration-manager-simple.py update-file --path /path/to/your/project --file src/main.go
```

## 📊 输出文件说明

工具会生成以下文件：

- `all_declarations.json` - 所有文件声明
- `file_declarations.json` - 指定文件声明
- `project_declarations.json` - 项目声明文档
- `updated_declarations.json` - 更新记录

### 输出格式示例

#### 所有文件声明输出

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

#### 文件声明输出

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

## 🎯 使用场景

### 1. 新项目分析
```bash
# 快速了解新项目结构
python declaration-manager-simple.py create-project
```

### 2. 代码审查
```bash
# 获取特定文件的详细信息
python declaration-manager-simple.py get-file --file src/api.go
```

### 3. 项目文档生成
```bash
# 生成完整的项目声明文档
python declaration-manager-simple.py create-project --output project_docs.json
```

### 4. 增量更新
```bash
# 更新修改过的文件
python declaration-manager-simple.py update-file --file src/main.go
```

## ⚙️ 配置选项

### 命令行参数

```bash
python declaration-manager-simple.py [action] [options]

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

### 1. 批量处理
```bash
# 处理多个文件
for file in src/*.go; do
    python declaration-manager-simple.py update-file --file "$file"
done
```

### 2. 集成到脚本
```bash
#!/bin/bash
# 项目分析脚本
python declaration-manager-simple.py create-project
python declaration-manager-simple.py get-all --output analysis.json
```

### 3. CI/CD 集成
```yaml
# GitHub Actions 示例
- name: Generate Project Declarations
  run: |
    python declaration-manager-simple.py create-project
    # 上传到存储或发送通知
```

## 📈 性能优化

### 缓存机制

工具支持智能缓存机制：

- **缓存位置**: 临时文件缓存
- **缓存策略**: 基于文件修改时间
- **清理机制**: 自动清理临时文件

### 性能指标

- 小项目（<100文件）：< 10秒
- 中等项目（100-1000文件）：< 1分钟
- 大项目（>1000文件）：< 5分钟

## 🐛 故障排除

### 常见问题

**Q: 找不到 contextgen 可执行文件**
```bash
# 确保 CodeCartographer 已构建
cd ../../  # 回到项目根目录
make build
```

**Q: 权限问题**
```bash
# 给脚本执行权限
chmod +x declaration-manager-simple.py
chmod +x install-spec-kit.py
```

**Q: Cursor 中找不到工具**
```bash
# 重新安装
python install-spec-kit.py uninstall
python install-spec-kit.py install
# 然后重启 Cursor
```

**Q: 编码问题**
```bash
# 设置环境变量
export PYTHONIOENCODING=utf-8
# 或在 Windows 中
set PYTHONIOENCODING=utf-8
```

### 调试模式

```bash
# 启用详细输出
python declaration-manager-simple.py get-all --verbose

# 检查安装状态
python install-spec-kit.py check
```

## 🎉 完成！

现在您已经掌握了 CodeCartographer 声明管理工具的基本用法。这个工具将帮助您：

- 🚀 快速了解项目结构
- 📊 分析代码声明和依赖
- 📝 生成项目文档
- 🔄 维护代码一致性

开始使用吧！如果遇到问题，请查看完整的 [README.md](README.md) 文档。
