package parser

import (
	"strings"

	"github.com/cnwinds/code-outline/internal/models"
	sitter "github.com/smacker/go-tree-sitter"
)

// GoExtractor Go语言提取器
type GoExtractor struct {
	BaseExtractor
	queries []string
}

// NewGoExtractor 创建Go语言提取器
func NewGoExtractor() *GoExtractor {
	return &GoExtractor{
		queries: []string{
			"(function_declaration) @symbol",
			"(method_declaration) @symbol",
			"(type_declaration) @symbol",
			"(var_declaration) @symbol",
			"(const_declaration) @symbol",
		},
	}
}

// GetQueries 获取Go语言的Tree-sitter查询规则
func (g *GoExtractor) GetQueries() []string {
	return g.queries
}

// ExtractPrototype 提取Go函数/类型原型
func (g *GoExtractor) ExtractPrototype(node *sitter.Node, content []byte) string {
	nodeType := node.Type()

	// 对于函数定义，提取函数签名（不包含函数体）
	if nodeType == "function_declaration" || nodeType == "method_declaration" {
		return g.extractFunctionPrototype(node, content, g.IsFunctionBodyNode)
	}

	// 对于类型定义，返回完整内容
	return g.extractFullNode(node, content)
}

// ExtractMethods 提取Go方法（Go没有类，方法通过receiver关联）
func (g *GoExtractor) ExtractMethods(classNode *sitter.Node, content []byte) []models.Symbol {
	// Go语言没有传统意义上的类，方法通过receiver关联到类型
	// 这里返回空，因为Go的方法会作为独立的函数提取
	return []models.Symbol{}
}

// IsClassNode 检查是否是类节点
func (g *GoExtractor) IsClassNode(nodeType string) bool {
	// Go语言没有类，只有结构体
	return nodeType == "type_declaration" || nodeType == "struct_type"
}

// IsFunctionBodyNode 检查是否是函数体节点
func (g *GoExtractor) IsFunctionBodyNode(nodeType string) bool {
	return nodeType == "block" || nodeType == "function_body"
}

// IsInsideClass 检查节点是否在类内部
func (g *GoExtractor) IsInsideClass(node *sitter.Node) bool {
	// Go语言没有类，所以总是返回false
	return false
}

// ExtractComments 提取Go注释
func (g *GoExtractor) ExtractComments(node *sitter.Node, content []byte) string {
	// Go语言使用 // 和 /* */ 注释
	return g.extractGoComments(node, content)
}

// extractGoComments 提取Go注释
func (g *GoExtractor) extractGoComments(node *sitter.Node, content []byte) string {
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
