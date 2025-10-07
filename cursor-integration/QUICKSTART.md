# CodeCartographer 声明管理工具 - 快速开始

## 🚀 5分钟快速上手

### 1. 安装工具

```bash
# 进入工具目录
cd cursor-integration

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

## 🎯 使用场景

### 1. 新项目分析
```bash
# 快速了解新项目结构
python declaration-manager.py create-project
```

### 2. 代码审查
```bash
# 获取特定文件的详细信息
python declaration-manager.py get-file --file src/api.go
```

### 3. 项目文档生成
```bash
# 生成完整的项目声明文档
python declaration-manager.py create-project --output project_docs.json
```

### 4. 增量更新
```bash
# 更新修改过的文件
python declaration-manager.py update-file --file src/main.go
```

## ⚡ 性能提示

- 首次运行会较慢，后续使用缓存会很快
- 大项目建议使用 `--no-cache` 强制刷新
- 可以设置 `--timeout` 参数调整超时时间

## 🔧 故障排除

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
chmod +x declaration-manager.py
chmod +x install-spec-kit.py
```

**Q: Cursor 中找不到工具**
```bash
# 重新安装
python install-spec-kit.py uninstall
python install-spec-kit.py install
# 然后重启 Cursor
```

## 📈 高级用法

### 批量处理
```bash
# 处理多个文件
for file in src/*.go; do
    python declaration-manager.py update-file --file "$file"
done
```

### 集成到脚本
```bash
#!/bin/bash
# 项目分析脚本
python declaration-manager.py create-project
python declaration-manager.py get-all --output analysis.json
```

### CI/CD 集成
```yaml
# GitHub Actions 示例
- name: Generate Project Declarations
  run: |
    python declaration-manager.py create-project
    # 上传到存储或发送通知
```

## 🎉 完成！

现在您已经掌握了 CodeCartographer 声明管理工具的基本用法。这个工具将帮助您：

- 🚀 快速了解项目结构
- 📊 分析代码声明和依赖
- 📝 生成项目文档
- 🔄 维护代码一致性

开始使用吧！如果遇到问题，请查看完整的 [README.md](README.md) 文档。
