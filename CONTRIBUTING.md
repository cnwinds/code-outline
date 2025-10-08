# 贡献指南

感谢您对 code-outline 项目的关注！我们欢迎任何形式的贡献。

## 🚀 快速开始

### 环境要求

- Go 1.21 或更高版本
- Git
- 基本的命令行工具

### 开发环境设置

1. **Fork 并克隆项目**
   ```bash
   git clone https://github.com/yourusername/code-outline.git
   cd code-outline
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **构建项目**
   ```bash
   make build
   ```

4. **运行测试**
   ```bash
   make test
   ```

## 📝 贡献类型

### 🐛 Bug 报告

在提交 Bug 报告前，请：

1. 检查 [Issues](https://github.com/cnwinds/code-outline/issues) 是否已存在相同问题
2. 使用最新版本进行测试
3. 提供详细的复现步骤

**Bug 报告模板**：
- 操作系统和版本
- Go 版本
- 复现步骤
- 期望行为 vs 实际行为
- 相关日志或错误信息

### ✨ 功能请求

我们欢迎新功能建议！请：

1. 检查现有功能是否已满足需求
2. 详细描述使用场景
3. 说明预期效果

### 🔧 代码贡献

#### 开发流程

1. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

2. **编写代码**
   - 遵循 Go 编码规范
   - 添加适当的注释
   - 编写单元测试

3. **测试**
   ```bash
   make test
   make lint
   ```

4. **提交代码**
   ```bash
   git add .
   git commit -m "feat: 添加新功能描述"
   git push origin feature/your-feature-name
   ```

5. **创建 Pull Request**

#### 代码规范

- **命名**: 使用驼峰命名法
- **注释**: 公共函数必须有注释
- **错误处理**: 使用 `%w` 包装错误
- **测试**: 新功能必须有对应测试

#### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**类型**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**示例**:
```
feat(parser): 添加 Python 类方法解析支持

- 支持解析 Python 类中的方法
- 提取方法参数和返回类型
- 添加相应的单元测试

Closes #123
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/parser

# 运行测试并显示覆盖率
go test -cover ./...
```

### 测试覆盖率

我们要求新代码的测试覆盖率至少达到 70%。

### 基准测试

```bash
make bench
```

## 📚 文档

### 更新文档

- README.md: 主要功能说明
- CONTRIBUTING.md: 贡献指南（本文件）
- CHANGELOG.md: 版本变更记录
- 代码注释: 函数和类型说明

### 文档规范

- 使用中文编写
- 保持格式一致
- 提供示例代码
- 及时更新过时信息

## 🔍 代码审查

### 审查流程

1. 所有 PR 都需要经过代码审查
2. 至少需要一位维护者批准
3. 通过所有 CI 检查
4. 解决所有审查意见

### 审查要点

- 代码质量和可读性
- 测试覆盖率
- 性能影响
- 安全性考虑
- 向后兼容性

## 🏗️ 项目结构

```
code-outline/
├── cmd/code-outline/          # 主程序入口
├── internal/                # 内部包
│   ├── cmd/                # CLI 命令
│   ├── config/             # 配置管理
│   ├── models/             # 数据模型
│   ├── parser/             # 代码解析器
│   ├── scanner/            # 文件扫描器
│   └── updater/            # 增量更新
├── scripts/                # 构建脚本
├── docs/                   # 文档
└── testdata/              # 测试数据
```

## 🚨 安全

### 安全漏洞报告

如果您发现了安全漏洞，请：

1. **不要** 在公开 Issues 中报告
2. 发送邮件至: security@example.com
3. 提供详细的漏洞描述和复现步骤

### 安全编码实践

- 验证所有用户输入
- 避免路径遍历攻击
- 限制资源使用
- 使用安全的文件操作

## 📞 获取帮助

### 社区支持

- GitHub Issues: 技术问题和 Bug 报告
- GitHub Discussions: 功能讨论和问题咨询
- 邮件: support@example.com

### 维护者

- @maintainer1 - 项目负责人
- @maintainer2 - 核心开发者

## 📄 许可证

本项目采用 MIT 许可证。贡献代码即表示您同意将代码以相同许可证发布。

## 🙏 致谢

感谢所有为 code-outline 项目做出贡献的开发者！

---

**Happy Coding! 🎉**
