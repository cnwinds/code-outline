# Tree-sitter 语法库构建脚本 (Windows PowerShell)
# 根据 tree-sitter.md 文档实现

param(
    [switch]$Clean
)

$ErrorActionPreference = "Stop"

$GRAMMAR_DIR = ".\grammars"
$TEMP_DIR = ".\temp_grammars"

Write-Host "🔧 开始构建 Tree-sitter 语法库..." -ForegroundColor Green

# 检查是否安装了必要的工具
$gccInstalled = Get-Command gcc -ErrorAction SilentlyContinue
if (-not $gccInstalled) {
    Write-Host "❌ 错误: 需要安装 GCC 编译器" -ForegroundColor Red
    Write-Host "请安装 MinGW-w64 或 MSYS2" -ForegroundColor Yellow
    exit 1
}

$gitInstalled = Get-Command git -ErrorAction SilentlyContinue
if (-not $gitInstalled) {
    Write-Host "❌ 错误: 需要安装 Git" -ForegroundColor Red
    exit 1
}

# 创建目录
if (Test-Path $GRAMMAR_DIR) {
    if ($Clean) {
        Remove-Item $GRAMMAR_DIR -Recurse -Force
    }
}
New-Item -ItemType Directory -Force -Path $GRAMMAR_DIR | Out-Null

if (Test-Path $TEMP_DIR) {
    Remove-Item $TEMP_DIR -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $TEMP_DIR | Out-Null

Set-Location $TEMP_DIR

# 支持的语言列表
$LANGUAGES = @{
    "go" = "https://github.com/tree-sitter/tree-sitter-go"
    "javascript" = "https://github.com/tree-sitter/tree-sitter-javascript"
    "python" = "https://github.com/tree-sitter/tree-sitter-python"
    "typescript" = "https://github.com/tree-sitter/tree-sitter-typescript"
    "java" = "https://github.com/tree-sitter/tree-sitter-java"
    "rust" = "https://github.com/tree-sitter/tree-sitter-rust"
    "cpp" = "https://github.com/tree-sitter/tree-sitter-cpp"
    "c" = "https://github.com/tree-sitter/tree-sitter-c"
}

# 编译每种语言的语法库
foreach ($lang in $LANGUAGES.Keys) {
    Write-Host "📦 构建 $lang 语法库..." -ForegroundColor Cyan
    
    # 克隆仓库
    if (-not (Test-Path $lang)) {
        Write-Host "📥 克隆 $lang 仓库..."
        git clone $LANGUAGES[$lang] $lang
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ 克隆 $lang 失败" -ForegroundColor Red
            continue
        }
    }
    
    Set-Location $lang
    
    # 检查必要的源文件
    if (-not (Test-Path "src\parser.c")) {
        Write-Host "❌ 错误: src\parser.c 不存在于 $lang" -ForegroundColor Red
        Set-Location ..
        continue
    }
    
    # 编译命令
    Write-Host "🔨 编译 $lang..."
    
    $outputFile = "..\$GRAMMAR_DIR\tree-sitter-$lang.dll"
    
    # 检查是否有 scanner 文件
    $hasScanner = $false
    $scannerFile = ""
    
    if (Test-Path "src\scanner.c") {
        $hasScanner = $true
        $scannerFile = "src\scanner.c"
        Write-Host "🔨 包含 scanner.c..."
    } elseif (Test-Path "src\scanner.cc") {
        $hasScanner = $true
        $scannerFile = "src\scanner.cc"
        Write-Host "🔨 包含 scanner.cc..."
    }
    
    try {
        if ($hasScanner) {
            if ($scannerFile.EndsWith(".cc")) {
                # 使用 g++ 编译 C++
                & g++ -shared -fPIC -Isrc src\parser.c $scannerFile -o $outputFile
            } else {
                # 使用 gcc 编译 C
                & gcc -shared -fPIC -Isrc src\parser.c $scannerFile -o $outputFile
            }
        } else {
            # 只编译 parser.c
            & gcc -shared -fPIC -Isrc src\parser.c -o $outputFile
        }
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ $lang 语法库构建成功" -ForegroundColor Green
        } else {
            Write-Host "❌ $lang 语法库构建失败" -ForegroundColor Red
        }
    } catch {
        Write-Host "❌ $lang 编译时发生错误: $_" -ForegroundColor Red
    }
    
    Set-Location ..
}

Set-Location ..

# 清理临时文件
Write-Host "🧹 清理临时文件..." -ForegroundColor Yellow
Remove-Item $TEMP_DIR -Recurse -Force

Write-Host "🎉 语法库构建完成！" -ForegroundColor Green
Write-Host "📁 语法库位置: $GRAMMAR_DIR\" -ForegroundColor Cyan

if (Test-Path $GRAMMAR_DIR) {
    Get-ChildItem $GRAMMAR_DIR -File | ForEach-Object {
        $sizeKB = [math]::Round($_.Length / 1KB, 2)
        Write-Host "  $($_.Name) ($sizeKB KB)" -ForegroundColor Gray
    }
}
