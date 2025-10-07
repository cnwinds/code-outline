# 变更日志

所有重要的项目变更都会记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 计划中
- Tree-sitter 解析器集成
- 更多编程语言支持
- 依赖关系分析
- 函数调用图生成

## [0.1.0] - 2025-10-07

### 新增
- 🎉 初始版本发布
- ✨ 支持 9+ 种编程语言的代码解析
- ✨ 基于正则表达式的符号提取
- ✨ 并发文件扫描，支持大型项目
- ✨ 增量更新功能，提升性能
- ✨ 灵活的排除规则配置
- ✨ 为 LLM 优化的 JSON 输出格式
- ✨ 跨平台支持（Windows、Linux、macOS）
- ✨ Docker 容器化支持
- ✨ 完整的 Makefile 构建系统

### 功能
- **CLI 命令**:
  - `generate`: 全量生成项目上下文
  - `update`: 增量更新项目上下文
  - `version`: 显示版本信息

- **支持的语言**:
  - Go: 函数、方法、结构体、接口、常量、变量
  - JavaScript: 函数、箭头函数、类、接口
  - TypeScript: 函数、类、接口、类型别名
  - Python: 函数、类
  - Java: 方法、类、接口
  - C#: 方法、类、接口、结构体、属性
  - Rust: 函数、结构体、枚举、特征、实现
  - C++: 函数、类、结构体、命名空间
  - C: 函数、结构体、枚举

- **输出格式**:
  - 项目基本信息（名称、技术栈、最后更新）
  - 架构概述和模块摘要
  - 文件级别的符号信息
  - 符号原型、用途、行号范围
  - 支持容器类型的内部内容提取

### 技术特性
- **高性能**: 基于 Go 的并发处理
- **可扩展**: 模块化架构设计
- **可配置**: 支持自定义语言配置
- **容错性**: 单个文件解析失败不影响整体流程
- **内存效率**: 流式处理大型文件

### 已知限制
- Tree-sitter 解析器尚未实现（计划在 v0.2.0）
- 测试覆盖率较低（计划在 v0.2.0 改进）
- 缺少依赖关系分析
- 未支持函数调用图生成

### 文档
- 📖 完整的 README.md 使用指南
- 📖 详细的 CONTRIBUTING.md 贡献指南
- 📖 Tree-sitter 集成说明文档
- 📖 项目审查报告

### 构建和部署
- 🔨 Makefile 自动化构建
- 🐳 Docker 多阶段构建
- 📦 跨平台编译脚本
- 🚀 GitHub Actions CI/CD（计划中）

---

## 版本说明

### 版本号格式
我们使用 [语义化版本](https://semver.org/lang/zh-CN/) 格式：`主版本号.次版本号.修订号`

- **主版本号**: 不兼容的 API 修改
- **次版本号**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

### 发布周期
- **Alpha**: 内部测试版本，功能可能不完整
- **Beta**: 公开测试版本，功能基本完整
- **RC**: 候选发布版本，接近正式版
- **Stable**: 正式发布版本，生产就绪

### 支持策略
- **当前版本**: 完全支持，包括新功能和 Bug 修复
- **前一个主版本**: 仅支持关键 Bug 修复
- **更早版本**: 不再支持

---

## 贡献者

感谢所有为 CodeCartographer 项目做出贡献的开发者！

### 核心团队
- @maintainer1 - 项目负责人
- @maintainer2 - 核心开发者

### 贡献者
- 感谢所有提交 Issue 和 Pull Request 的社区成员

---

## 链接

- [项目主页](https://github.com/cnwinds/CodeCartographer)
- [问题报告](https://github.com/cnwinds/CodeCartographer/issues)
- [功能请求](https://github.com/cnwinds/CodeCartographer/discussions)
- [贡献指南](CONTRIBUTING.md)
- [许可证](LICENSE)
