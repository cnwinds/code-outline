package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"

	"github.com/cnwinds/CodeCartographer/internal/config"
	"github.com/cnwinds/CodeCartographer/internal/models"
)

// TreeSitterParser Tree-sitter 解析器
type TreeSitterParser struct {
	languagesConfig models.LanguagesConfig
	parsers         map[string]*sitter.Parser
}

// NewTreeSitterParser 创建新的 Tree-sitter 解析器
func NewTreeSitterParser(languagesConfig models.LanguagesConfig) (*TreeSitterParser, error) {
	p := &TreeSitterParser{
		languagesConfig: languagesConfig,
		parsers:         make(map[string]*sitter.Parser),
	}

	// 初始化各语言解析器
	if err := p.initParsers(); err != nil {
		return nil, err
	}

	return p, nil
}

// initParsers 初始化语言解析器
func (p *TreeSitterParser) initParsers() error {
	// Go 语言
	goParser := sitter.NewParser()
	goParser.SetLanguage(golang.GetLanguage())
	p.parsers["go"] = goParser

	// JavaScript
	jsParser := sitter.NewParser()
	jsParser.SetLanguage(javascript.GetLanguage())
	p.parsers["javascript"] = jsParser
	p.parsers["typescript"] = jsParser // 暂时共用

	// Python
	pyParser := sitter.NewParser()
	pyParser.SetLanguage(python.GetLanguage())
	p.parsers["python"] = pyParser

	return nil
}

// ParseFile 解析单个文件
func (p *TreeSitterParser) ParseFile(filePath string) (*models.FileInfo, error) {
	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 确定语言
	ext := filepath.Ext(filePath)
	langName, _, found := config.GetLanguageByExtension(p.languagesConfig, ext)
	if !found {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 为每次解析创建新的解析器实例（tree-sitter 不是线程安全的）
	parser := sitter.NewParser()
	var language *sitter.Language
	switch langName {
	case "go":
		language = golang.GetLanguage()
	case "javascript", "typescript":
		language = javascript.GetLanguage()
	case "python":
		language = python.GetLanguage()
	default:
		return nil, fmt.Errorf("未找到 %s 语言的解析器", langName)
	}
	parser.SetLanguage(language)

	// 使用 defer-recover 捕获可能的 panic
	var symbols []models.Symbol
	var parseErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				parseErr = fmt.Errorf("解析文件 %s 时发生错误: %v", filePath, r)
			}
		}()

		// 解析
		tree := parser.Parse(nil, content)
		if tree == nil {
			parseErr = fmt.Errorf("解析失败: tree is nil")
			return
		}
		defer tree.Close()

		rootNode := tree.RootNode()
		if rootNode == nil {
			parseErr = fmt.Errorf("解析失败: root node is nil")
			return
		}

		// 提取符号
		symbols = p.extractSymbols(rootNode, content, langName)
	}()

	if parseErr != nil {
		return nil, parseErr
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &models.FileInfo{
		Purpose:      p.extractFilePurpose(content),
		Symbols:      symbols,
		LastModified: fileInfo.ModTime().Format(time.RFC3339),
		FileSize:     fileInfo.Size(),
	}, nil
}

// extractSymbols 从语法树提取符号
func (p *TreeSitterParser) extractSymbols(node *sitter.Node, content []byte, lang string) []models.Symbol {
	var symbols []models.Symbol

	// 获取查询规则
	langConfig, exists := p.languagesConfig[lang]
	if !exists {
		return symbols
	}

	queries := langConfig.Queries.TopLevelSymbols

	// 获取语言对象
	var language *sitter.Language
	switch lang {
	case "go":
		language = golang.GetLanguage()
	case "javascript", "typescript":
		language = javascript.GetLanguage()
	case "python":
		language = python.GetLanguage()
	default:
		return symbols
	}

	for _, queryStr := range queries {
		query, err := sitter.NewQuery([]byte(queryStr), language)
		if err != nil {
			continue
		}
		defer query.Close()

		cursor := sitter.NewQueryCursor()
		defer cursor.Close()

		cursor.Exec(query, node)

		for {
			match, ok := cursor.NextMatch()
			if !ok {
				break
			}

			for _, capture := range match.Captures {
				symbol := p.nodeToSymbol(capture.Node, content)
				symbols = append(symbols, symbol)
			}
		}
	}

	return symbols
}

// nodeToSymbol 将语法树节点转换为符号
func (p *TreeSitterParser) nodeToSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	return models.Symbol{
		Prototype: string(content[node.StartByte():node.EndByte()]),
		Purpose:   "", // TODO: 从注释提取
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}
}

// extractFilePurpose 提取文件用途
func (p *TreeSitterParser) extractFilePurpose(content []byte) string {
	lines := strings.Split(string(content), "\n")

	// 查找文件顶部的注释
	for _, line := range lines[:minInt(10, len(lines))] {
		trimmed := strings.TrimSpace(line)

		// 跳过package声明等
		if strings.HasPrefix(trimmed, "package") ||
			strings.HasPrefix(trimmed, "import") ||
			trimmed == "" {
			continue
		}

		// 查找注释
		if strings.HasPrefix(trimmed, "//") {
			purpose := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
			if len(purpose) > 10 { // 过滤太短的注释
				return purpose
			}
		}

		// 如果遇到代码行，停止查找
		if !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "#") {
			break
		}
	}

	return "TODO: Describe the purpose of this file."
}

// minInt 返回两个整数中的较小者
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
