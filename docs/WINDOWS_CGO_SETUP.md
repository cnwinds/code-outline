# Windows CGO 环境安装指南

本文档详细说明如何在 Windows 环境下安装和配置 CGO 编译环境，以支持 Tree-sitter 解析器的构建。

## 📋 系统要求

- Windows 10/11 (64位)
- 管理员权限（用于安装软件）
- 网络连接（用于下载安装包）

## 🚀 安装步骤

### 步骤 1: 下载并安装 MSYS2

1. **访问 MSYS2 官网**
   - 打开浏览器访问: https://www.msys2.org/
   - 点击 "Download" 下载最新版本的 MSYS2 安装包

2. **运行安装程序**
   - 双击下载的 `.exe` 文件
   - 按照安装向导完成安装
   - 默认安装路径: `C:\msys64\`

3. **启动 MSYS2**
   - 安装完成后，从开始菜单启动 "MSYS2 UCRT64"
   - 或者直接运行: `C:\msys64\msys2_shell.cmd -ucrt64`

### 步骤 2: 安装 MinGW-w64 编译器

在 MSYS2 终端中执行以下命令：

```bash
# 更新包数据库
pacman -Syu

# 安装 MinGW-w64 编译器工具链
pacman -S mingw-w64-ucrt-x86_64-gcc
pacman -S mingw-w64-ucrt-x86_64-gdb
pacman -S mingw-w64-ucrt-x86_64-make

# 安装其他必要的开发工具
pacman -S mingw-w64-ucrt-x86_64-pkg-config
pacman -S mingw-w64-ucrt-x86_64-cmake
```

### 步骤 3: 配置环境变量

#### 方法 1: 通过系统设置（推荐）

1. **打开系统环境变量设置**
   - 按 `Win + R`，输入 `sysdm.cpl`，回车
   - 点击 "高级" 标签页
   - 点击 "环境变量" 按钮

2. **添加 PATH 环境变量**
   - 在 "系统变量" 中找到 `Path` 变量
   - 点击 "编辑"
   - 点击 "新建"，添加: `C:\msys64\ucrt64\bin`
   - 点击 "确定" 保存

#### 方法 2: 通过 PowerShell（临时）

```powershell
# 临时设置环境变量（仅当前会话有效）
$env:PATH += ";C:\msys64\ucrt64\bin"

# 永久设置环境变量（需要管理员权限）
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";C:\msys64\ucrt64\bin", "Machine")
```

### 步骤 4: 设置 CGO 环境变量

在 PowerShell 或命令提示符中设置：

```powershell
# 启用 CGO
$env:CGO_ENABLED = "1"

# 设置 C 编译器
$env:CC = "gcc"

# 设置 C++ 编译器
$env:CXX = "g++"
```

### 步骤 5: 验证安装

打开新的 PowerShell 窗口，执行以下命令验证安装：

```powershell
# 检查 GCC 版本
gcc --version

# 检查 G++ 版本
g++ --version

# 检查 CGO 环境
go env CGO_ENABLED
go env CC
```

预期输出示例：
```
gcc (Rev10, Built by MSYS2 project) 13.2.0
Copyright (C) 2023 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

CGO_ENABLED=1
CC=gcc
```

## 🔧 常见问题排查

### 问题 1: `gcc: command not found`

**原因**: PATH 环境变量未正确设置

**解决方案**:
1. 确认 MSYS2 已正确安装
2. 检查 `C:\msys64\ucrt64\bin` 目录是否存在
3. 重新设置 PATH 环境变量
4. 重启 PowerShell 或命令提示符

### 问题 2: `CGO_ENABLED=0`

**原因**: CGO 未启用

**解决方案**:
```powershell
# 设置环境变量
$env:CGO_ENABLED = "1"

# 验证设置
go env CGO_ENABLED
```

### 问题 3: 编译时出现链接错误

**原因**: 缺少必要的库文件

**解决方案**:
```bash
# 在 MSYS2 中安装额外的开发库
pacman -S mingw-w64-ucrt-x86_64-libffi
pacman -S mingw-w64-ucrt-x86_64-zlib
```

### 问题 4: 权限不足

**原因**: 需要管理员权限

**解决方案**:
1. 以管理员身份运行 PowerShell
2. 或者使用 MSYS2 终端进行开发

## 🧪 测试 CGO 环境

创建一个简单的测试程序验证 CGO 环境：

**文件**: `test_cgo.go`
```go
package main

/*
#include <stdio.h>
void hello() {
    printf("Hello from C!\n");
}
*/
import "C"

func main() {
    C.hello()
}
```

**编译和运行**:
```powershell
# 编译
go build test_cgo.go

# 运行
./test_cgo.exe
```

预期输出: `Hello from C!`

## 📝 开发环境配置

### Visual Studio Code 配置

在 VS Code 中设置 Go 开发环境：

1. **安装 Go 扩展**
2. **配置 settings.json**:
```json
{
    "go.toolsEnvVars": {
        "CGO_ENABLED": "1",
        "CC": "gcc",
        "CXX": "g++"
    }
}
```

### 项目构建配置

在项目根目录创建 `.env` 文件：
```
CGO_ENABLED=1
CC=gcc
CXX=g++
```

## 🎯 下一步

环境配置完成后，您可以：

1. 运行 `go mod tidy` 安装 Tree-sitter 依赖
2. 执行 `make build` 构建项目
3. 运行测试验证 Tree-sitter 功能

## 📚 参考资源

- [MSYS2 官网](https://www.msys2.org/)
- [MinGW-w64 文档](https://www.mingw-w64.org/)
- [Go CGO 文档](https://pkg.go.dev/cmd/cgo)
- [Tree-sitter Go 绑定](https://github.com/smacker/go-tree-sitter)

---

**注意**: 如果在安装过程中遇到问题，请检查：
1. 网络连接是否正常
2. 是否有足够的磁盘空间
3. 是否以管理员权限运行安装程序
