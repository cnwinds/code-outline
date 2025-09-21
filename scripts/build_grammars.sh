#!/bin/bash

# Tree-sitter 语法库构建脚本
# 根据 tree-sitter.md 文档实现

set -e

GRAMMAR_DIR="./grammars"
TEMP_DIR="./temp_grammars"

echo "🔧 开始构建 Tree-sitter 语法库..."

# 创建目录
mkdir -p "$GRAMMAR_DIR"
mkdir -p "$TEMP_DIR"

cd "$TEMP_DIR"

# 支持的语言列表
declare -A LANGUAGES=(
    ["go"]="https://github.com/tree-sitter/tree-sitter-go"
    ["javascript"]="https://github.com/tree-sitter/tree-sitter-javascript"
    ["python"]="https://github.com/tree-sitter/tree-sitter-python"
    ["typescript"]="https://github.com/tree-sitter/tree-sitter-typescript"
    ["java"]="https://github.com/tree-sitter/tree-sitter-java"
    ["rust"]="https://github.com/tree-sitter/tree-sitter-rust"
    ["cpp"]="https://github.com/tree-sitter/tree-sitter-cpp"
    ["c"]="https://github.com/tree-sitter/tree-sitter-c"
)

# 编译每种语言的语法库
for lang in "${!LANGUAGES[@]}"; do
    echo "📦 构建 $lang 语法库..."
    
    # 克隆仓库
    if [ ! -d "$lang" ]; then
        git clone "${LANGUAGES[$lang]}" "$lang"
    fi
    
    cd "$lang"
    
    # 检查必要的源文件
    if [ ! -f "src/parser.c" ]; then
        echo "❌ 错误: src/parser.c 不存在于 $lang"
        cd ..
        continue
    fi
    
    # 编译命令
    echo "🔨 编译 $lang..."
    
    # 根据操作系统选择编译参数
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        gcc -shared -fPIC -Isrc src/parser.c -o "../$GRAMMAR_DIR/tree-sitter-$lang.so"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        gcc -shared -fPIC -Isrc src/parser.c -o "../$GRAMMAR_DIR/tree-sitter-$lang.dylib"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        # Windows
        gcc -shared -fPIC -Isrc src/parser.c -o "../$GRAMMAR_DIR/tree-sitter-$lang.dll"
    fi
    
    # 检查是否有 scanner.c 或 scanner.cc
    if [ -f "src/scanner.c" ]; then
        echo "🔨 包含 scanner.c..."
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            gcc -shared -fPIC -Isrc src/parser.c src/scanner.c -o "../$GRAMMAR_DIR/tree-sitter-$lang.so"
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            gcc -shared -fPIC -Isrc src/parser.c src/scanner.c -o "../$GRAMMAR_DIR/tree-sitter-$lang.dylib"
        elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
            gcc -shared -fPIC -Isrc src/parser.c src/scanner.c -o "../$GRAMMAR_DIR/tree-sitter-$lang.dll"
        fi
    elif [ -f "src/scanner.cc" ]; then
        echo "🔨 包含 scanner.cc..."
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            g++ -shared -fPIC -Isrc src/parser.c src/scanner.cc -o "../$GRAMMAR_DIR/tree-sitter-$lang.so"
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            g++ -shared -fPIC -Isrc src/parser.c src/scanner.cc -o "../$GRAMMAR_DIR/tree-sitter-$lang.dylib"
        elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
            g++ -shared -fPIC -Isrc src/parser.c src/scanner.cc -o "../$GRAMMAR_DIR/tree-sitter-$lang.dll"
        fi
    fi
    
    if [ $? -eq 0 ]; then
        echo "✅ $lang 语法库构建成功"
    else
        echo "❌ $lang 语法库构建失败"
    fi
    
    cd ..
done

cd ..

# 清理临时文件
echo "🧹 清理临时文件..."
rm -rf "$TEMP_DIR"

echo "🎉 语法库构建完成！"
echo "📁 语法库位置: $GRAMMAR_DIR/"
ls -la "$GRAMMAR_DIR/"
