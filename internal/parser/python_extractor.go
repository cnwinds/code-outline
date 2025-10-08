package parser

import (
	"strings"

	"github.com/cnwinds/code-outline/internal/models"
	sitter "github.com/smacker/go-tree-sitter"
)

// PythonExtractor Python语言提取器
type PythonExtractor struct {
	BaseExtractor
	queries []string
}

// NewPythonExtractor 创建Python语言提取器
func NewPythonExtractor() *PythonExtractor {
	return &PythonExtractor{
		queries: []string{
			"(function_definition) @symbol",
			"(class_definition) @symbol",
		},
	}
}

// GetQueries 获取Python语言的Tree-sitter查询规则
func (p *PythonExtractor) GetQueries() []string {
	return p.queries
}

// ExtractPrototype 提取Python函数/类原型
func (p *PythonExtractor) ExtractPrototype(node *sitter.Node, content []byte) string {
	nodeType := node.Type()

	// 对于类定义，只提取类声明部分（不包含类体）
	if nodeType == "class_definition" {
		return p.extractClassPrototype(node, content)
	}

	// 对于函数定义，提取函数签名
	if nodeType == "function_definition" {
		return p.extractFunctionPrototype(node, content, p.IsFunctionBodyNode)
	}

	return p.extractFullNode(node, content)
}

// ExtractMethods 提取Python类内部的方法
func (p *PythonExtractor) ExtractMethods(classNode *sitter.Node, content []byte) []models.Symbol {
	var methods []models.Symbol

	childCount := int(classNode.ChildCount())
	for i := 0; i < childCount; i++ {
		child := classNode.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找类体
		if childType == "block" {
			// 遍历类体中的方法
			bodyChildCount := int(child.ChildCount())
			for j := 0; j < bodyChildCount; j++ {
				bodyChild := child.Child(j)
				if bodyChild == nil {
					continue
				}

				bodyChildType := bodyChild.Type()
				if bodyChildType == "function_definition" {
					method := p.createMethodSymbol(bodyChild, content)
					methods = append(methods, method)
				}
			}
		}
	}

	return methods
}

// IsClassNode 检查是否是类节点
func (p *PythonExtractor) IsClassNode(nodeType string) bool {
	return nodeType == "class_definition"
}

// IsFunctionBodyNode 检查是否是函数体节点
func (p *PythonExtractor) IsFunctionBodyNode(nodeType string) bool {
	return nodeType == "block" || nodeType == "function_body"
}

// IsInsideClass 检查节点是否在类内部
func (p *PythonExtractor) IsInsideClass(node *sitter.Node) bool {
	current := node.Parent()
	for current != nil {
		nodeType := current.Type()
		if nodeType == "class_definition" {
			return true
		}
		current = current.Parent()
	}
	return false
}

// ExtractComments 提取Python注释
func (p *PythonExtractor) ExtractComments(node *sitter.Node, content []byte) string {
	return p.extractPythonComments(node, content)
}

// extractClassPrototype 提取Python类原型
func (p *PythonExtractor) extractClassPrototype(node *sitter.Node, content []byte) string {
	fullText := string(content[node.StartByte():node.EndByte()])
	lines := strings.Split(fullText, "\n")

	for i, line := range lines {
		if strings.Contains(line, ":") && !strings.Contains(line, "::") {
			idx := strings.Index(line, ":")
			if i == 0 {
				return strings.TrimSpace(line[:idx+1])
			}
			result := strings.Join(lines[:i], "\n") + "\n" + line[:idx+1]
			return strings.TrimSpace(result)
		}
	}
	return ""
}

// createMethodSymbol 创建方法符号
func (p *PythonExtractor) createMethodSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	prototype := p.extractFunctionPrototype(node, content, p.IsFunctionBodyNode)

	return models.Symbol{
		Prototype: prototype,
		Purpose:   p.extractPythonComments(node, content),
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}
}

// extractPythonComments 提取Python注释（支持docstring）
func (p *PythonExtractor) extractPythonComments(node *sitter.Node, content []byte) string {
	startPoint := node.StartPoint()
	startRow := int(startPoint.Row)
	lines := strings.Split(string(content), "\n")

	// 首先检查docstring
	for i := startRow; i < len(lines) && i < startRow+5; i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "\"\"\"") || strings.HasPrefix(line, "'''") {
			// 找到docstring开始
			quote := "\"\"\""
			if strings.HasPrefix(line, "'''") {
				quote = "'''"
			}

			// 检查是否是单行docstring
			if strings.HasSuffix(line, quote) && len(line) > 6 {
				docstring := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(line, quote), quote))
				if docstring != "" {
					return docstring
				}
			}

			// 多行docstring
			var docLines []string
			firstLine := strings.TrimSpace(strings.TrimPrefix(line, quote))
			if firstLine != "" {
				docLines = append(docLines, firstLine)
			}

			for j := i + 1; j < len(lines) && j < i+10; j++ {
				docLine := strings.TrimSpace(lines[j])
				if strings.HasSuffix(docLine, quote) {
					lastLine := strings.TrimSpace(strings.TrimSuffix(docLine, quote))
					if lastLine != "" {
						docLines = append(docLines, lastLine)
					}
					break
				}
				if docLine != "" {
					docLines = append(docLines, docLine)
				}
			}

			if len(docLines) > 0 {
				return strings.Join(docLines, " ")
			}
		}
	}

	// 检查普通注释
	for i := startRow - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 检查单行注释
		if strings.HasPrefix(line, "#") {
			comment := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			if comment != "" {
				return comment
			}
		}

		// 如果遇到非注释行，停止查找
		if !strings.HasPrefix(line, "#") && line != "" {
			break
		}
	}

	return ""
}
