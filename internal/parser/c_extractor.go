package parser

import (
	"strings"

	"github.com/cnwinds/code-outline/internal/models"
	sitter "github.com/smacker/go-tree-sitter"
)

// CExtractor C语言提取器
type CExtractor struct {
	BaseExtractor
	queries []string
}

// NewCExtractor 创建C语言提取器
func NewCExtractor() *CExtractor {
	return &CExtractor{
		queries: []string{
			"(function_definition) @symbol",
			"(type_definition) @symbol",
		},
	}
}

// GetQueries 获取C语言的Tree-sitter查询规则
func (c *CExtractor) GetQueries() []string {
	return c.queries
}

// ExtractPrototype 提取C函数原型
func (c *CExtractor) ExtractPrototype(node *sitter.Node, content []byte) string {
	nodeType := node.Type()

	// 对于函数定义，提取函数签名（不包含函数体）
	if nodeType == "function_definition" {
		return c.extractFunctionPrototype(node, content, c.IsFunctionBodyNode)
	}

	// 对于类型定义，返回完整内容（typedef）
	if nodeType == "type_definition" {
		return c.extractFullNode(node, content)
	}

	return c.extractFullNode(node, content)
}

// ExtractMethods 提取C方法（C语言没有类，所以返回空）
func (c *CExtractor) ExtractMethods(classNode *sitter.Node, content []byte) []models.Symbol {
	// C语言没有类，所以没有方法
	return []models.Symbol{}
}

// IsClassNode 检查是否是类节点
func (c *CExtractor) IsClassNode(nodeType string) bool {
	// C语言没有类，只有结构体
	return nodeType == "struct_specifier" ||
		nodeType == "union_specifier" ||
		nodeType == "enum_specifier" ||
		nodeType == "type_definition"
}

// IsFunctionBodyNode 检查是否是函数体节点
func (c *CExtractor) IsFunctionBodyNode(nodeType string) bool {
	return nodeType == "compound_statement" || nodeType == "block"
}

// IsInsideClass 检查节点是否在类内部
func (c *CExtractor) IsInsideClass(node *sitter.Node) bool {
	// C语言没有类，所以总是返回false
	return false
}

// ExtractComments 提取C注释
func (c *CExtractor) ExtractComments(node *sitter.Node, content []byte) string {
	return c.extractCComments(node, content)
}

// extractCComments 提取C注释
func (c *CExtractor) extractCComments(node *sitter.Node, content []byte) string {
	startPoint := node.StartPoint()
	startRow := int(startPoint.Row)
	lines := strings.Split(string(content), "\n")

	var commentLines []string
	inMultiLineComment := false

	for i := startRow - 1; i >= 0; i-- {
		if i >= len(lines) {
			continue
		}

		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 检查多行注释结束
		if strings.HasSuffix(line, "*/") {
			inMultiLineComment = true
			comment := strings.TrimSpace(strings.TrimSuffix(line, "*/"))
			if comment != "" {
				commentLines = append([]string{comment}, commentLines...)
			}
			continue
		}

		// 检查多行注释中间行
		if inMultiLineComment {
			if strings.HasPrefix(line, "*") {
				comment := strings.TrimSpace(strings.TrimPrefix(line, "*"))
				if comment != "" {
					commentLines = append([]string{comment}, commentLines...)
				}
				continue
			} else if strings.HasPrefix(line, "/*") {
				// 多行注释开始
				comment := strings.TrimSpace(strings.TrimPrefix(line, "/*"))
				if comment != "" {
					commentLines = append([]string{comment}, commentLines...)
				}
				inMultiLineComment = false
				break
			}
		}

		// 检查单行注释
		if strings.HasPrefix(line, "//") {
			comment := strings.TrimSpace(strings.TrimPrefix(line, "//"))
			if comment != "" {
				return comment
			}
		}

		// 如果遇到非注释行，停止查找
		if !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/*") && !strings.HasPrefix(line, "*") && line != "" {
			break
		}
	}

	// 如果有收集到的多行注释，合并它们
	if len(commentLines) > 0 {
		return strings.Join(commentLines, " ")
	}

	return ""
}
