package parser

import (
	"fmt"

	"github.com/cnwinds/CodeCartographer/internal/models"
)

// TreeSitterParser Tree-sitter 解析器（暂时禁用 CGO 版本）
type TreeSitterParser struct {
	languagesConfig models.LanguagesConfig
}

// NewTreeSitterParser 创建新的 Tree-sitter 解析器
func NewTreeSitterParser(languagesConfig models.LanguagesConfig) (*TreeSitterParser, error) {
	// 当前因为 CGO 依赖问题，暂时返回错误促使回退到简单解析器
	return nil, fmt.Errorf("Tree-sitter 需要 CGO 支持，当前环境未配置 C 编译器。请安装 MinGW-w64 或使用 --treesitter=false")
}

// ParseFile 解析单个文件（占位符实现）
func (p *TreeSitterParser) ParseFile(filePath string) (*models.FileInfo, error) {
	return nil, fmt.Errorf("Tree-sitter 解析器需要 CGO 支持，当前环境未配置")
}
