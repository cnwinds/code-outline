package parser

import (
	"strings"

	"github.com/cnwinds/CodeCartographer/internal/models"
	sitter "github.com/smacker/go-tree-sitter"
)

// JavaExtractor Java语言提取器
type JavaExtractor struct {
	BaseExtractor
	queries []string
}

// NewJavaExtractor 创建Java语言提取器
func NewJavaExtractor() *JavaExtractor {
	return &JavaExtractor{
		queries: []string{
			"(method_declaration) @symbol",
			"(class_declaration) @symbol",
			"(interface_declaration) @symbol",
			"(enum_declaration) @symbol",
			"(constructor_declaration) @symbol",
		},
	}
}

// GetQueries 获取Java语言的Tree-sitter查询规则
func (j *JavaExtractor) GetQueries() []string {
	return j.queries
}

// ExtractPrototype 提取Java类原型
func (j *JavaExtractor) ExtractPrototype(node *sitter.Node, content []byte) string {
	nodeType := node.Type()

	// 对于类声明，只提取类声明部分（不包含类体）
	if nodeType == "class_declaration" {
		return j.extractClassPrototype(node, content)
	}

	// 对于方法声明，提取方法签名
	if nodeType == "method_declaration" {
		return j.extractMethodPrototype(node, content)
	}

	return j.extractFullNode(node, content)
}

// ExtractMethods 提取Java类内部的方法
func (j *JavaExtractor) ExtractMethods(classNode *sitter.Node, content []byte) []models.Symbol {
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
				if bodyChildType == "method_declaration" || bodyChildType == "constructor_declaration" {
					method := j.createMethodSymbol(bodyChild, content)
					methods = append(methods, method)
				}
			}
		}
	}

	return methods
}

// IsClassNode 检查是否是类节点
func (j *JavaExtractor) IsClassNode(nodeType string) bool {
	return nodeType == "class_declaration" || nodeType == "interface_declaration" || nodeType == "enum_declaration"
}

// IsFunctionBodyNode 检查是否是函数体节点
func (j *JavaExtractor) IsFunctionBodyNode(nodeType string) bool {
	return nodeType == "block" || nodeType == "block_statement" || nodeType == "constructor_body"
}

// IsInsideClass 检查节点是否在类内部
func (j *JavaExtractor) IsInsideClass(node *sitter.Node) bool {
	current := node.Parent()
	for current != nil {
		nodeType := current.Type()
		if nodeType == "class_declaration" || nodeType == "interface_declaration" || nodeType == "enum_declaration" {
			return true
		}
		current = current.Parent()
	}
	return false
}

// ExtractComments 提取Java注释
func (j *JavaExtractor) ExtractComments(node *sitter.Node, content []byte) string {
	return extractJavaCommentsFixed(node, content)
}

// extractClassPrototype 提取Java类原型
func (j *JavaExtractor) extractClassPrototype(node *sitter.Node, content []byte) string {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找类体开始位置
		if childType == "class_body" || childType == "{" {
			prototype := string(content[node.StartByte():child.StartByte()])
			return strings.TrimSpace(prototype)
		}
	}
	return ""
}

// extractMethodPrototype 提取Java方法原型
func (j *JavaExtractor) extractMethodPrototype(node *sitter.Node, content []byte) string {
	// 首先尝试找到函数体结束位置
	declarationEnd := j.findFunctionBodyEnd(node, j.IsFunctionBodyNode)
	if declarationEnd > node.StartByte() {
		prototype := string(content[node.StartByte():declarationEnd])
		return j.cleanText(prototype)
	}

	// 如果没有找到函数体，手动查找block节点
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 如果遇到函数体，停止提取
		if childType == "block" {
			prototype := string(content[node.StartByte():child.StartByte()])
			return j.cleanText(prototype)
		}
	}

	// 如果都没有找到，返回整个节点内容
	return j.cleanText(string(content[node.StartByte():node.EndByte()]))
}

// createMethodSymbol 创建方法符号
func (j *JavaExtractor) createMethodSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	prototype := j.extractMethodPrototype(node, content)

	return models.Symbol{
		Prototype: prototype,
		Purpose:   extractJavaCommentsFixed(node, content),
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}
}

// extractJavaCommentsFixed 提取Java注释（修复版）
func extractJavaCommentsFixed(node *sitter.Node, content []byte) string {
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

		// 检查Javadoc注释结束
		if strings.HasSuffix(line, "*/") {
			inMultiLineComment = true
			comment := strings.TrimSpace(strings.TrimSuffix(line, "*/"))
			if comment != "" {
				commentLines = append([]string{comment}, commentLines...)
			}
			continue
		}

		// 检查Javadoc注释中间行
		if inMultiLineComment {
			if strings.HasPrefix(line, "*") {
				comment := strings.TrimSpace(strings.TrimPrefix(line, "*"))
				if comment != "" {
					commentLines = append([]string{comment}, commentLines...)
				}
				continue
			} else if strings.HasPrefix(line, "/**") {
				// Javadoc注释开始
				comment := strings.TrimSpace(strings.TrimPrefix(line, "/**"))
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
		if !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/**") && !strings.HasPrefix(line, "*") && line != "" {
			break
		}
	}

	// 如果有收集到的多行注释，合并它们
	if len(commentLines) > 0 {
		return strings.Join(commentLines, " ")
	}

	return ""
}
