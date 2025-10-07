package parser

import (
	"strings"

	"github.com/cnwinds/CodeCartographer/internal/models"
	sitter "github.com/smacker/go-tree-sitter"
)

// JSExtractor JavaScript/TypeScript语言提取器
type JSExtractor struct {
	BaseExtractor
	queries []string
}

// NewJSExtractor 创建JavaScript/TypeScript语言提取器
func NewJSExtractor() *JSExtractor {
	return &JSExtractor{
		queries: []string{
			"(function_declaration) @symbol",
			"(method_definition) @symbol",
			"(class_declaration) @symbol",
			"(interface_declaration) @symbol",
		},
	}
}

// GetQueries 获取JavaScript/TypeScript语言的Tree-sitter查询规则
func (j *JSExtractor) GetQueries() []string {
	return j.queries
}

// ExtractPrototype 提取JS/TS函数/类原型
func (j *JSExtractor) ExtractPrototype(node *sitter.Node, content []byte) string {
	nodeType := node.Type()

	// 对于类声明，只提取类声明部分（不包含类体）
	if nodeType == "class_declaration" {
		return j.extractClassPrototype(node, content)
	}

	// 对于函数定义，提取函数签名
	if nodeType == "function_declaration" || nodeType == "method_definition" {
		return j.extractFunctionPrototype(node, content, j.IsFunctionBodyNode)
	}

	return j.extractFullNode(node, content)
}

// ExtractMethods 提取JS/TS类内部的方法
func (j *JSExtractor) ExtractMethods(classNode *sitter.Node, content []byte) []models.Symbol {
	var methods []models.Symbol

	childCount := int(classNode.ChildCount())
	for i := 0; i < childCount; i++ {
		child := classNode.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找类体
		if childType == "class_body" {
			// 遍历类体中的方法
			bodyChildCount := int(child.ChildCount())
			for k := 0; k < bodyChildCount; k++ {
				bodyChild := child.Child(k)
				if bodyChild == nil {
					continue
				}

				bodyChildType := bodyChild.Type()
				if bodyChildType == "method_definition" {
					method := j.createMethodSymbol(bodyChild, content)
					methods = append(methods, method)
				}
			}
		}
	}

	return methods
}

// IsClassNode 检查是否是类节点
func (j *JSExtractor) IsClassNode(nodeType string) bool {
	return nodeType == "class_declaration"
}

// IsFunctionBodyNode 检查是否是函数体节点
func (j *JSExtractor) IsFunctionBodyNode(nodeType string) bool {
	return nodeType == "block" || nodeType == "function_body" || nodeType == "statement_block"
}

// IsInsideClass 检查节点是否在类内部
func (j *JSExtractor) IsInsideClass(node *sitter.Node) bool {
	current := node.Parent()
	for current != nil {
		nodeType := current.Type()
		if nodeType == "class_declaration" {
			return true
		}
		current = current.Parent()
	}
	return false
}

// ExtractComments 提取JS/TS注释
func (j *JSExtractor) ExtractComments(node *sitter.Node, content []byte) string {
	return j.extractJSComments(node, content)
}

// extractClassPrototype 提取JS/TS类原型
func (j *JSExtractor) extractClassPrototype(node *sitter.Node, content []byte) string {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找类体开始位置
		if childType == "class_body" {
			prototype := string(content[node.StartByte():child.StartByte()])
			return strings.TrimSpace(prototype)
		}
	}
	return ""
}

// createMethodSymbol 创建方法符号
func (j *JSExtractor) createMethodSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	prototype := j.extractFunctionPrototype(node, content, j.IsFunctionBodyNode)

	return models.Symbol{
		Prototype: prototype,
		Purpose:   j.extractJSComments(node, content),
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}
}

// extractJSComments 提取JS/TS注释
func (j *JSExtractor) extractJSComments(node *sitter.Node, content []byte) string {
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
