package parser

import (
	"strings"

	"github.com/cnwinds/code-outline/internal/models"
	sitter "github.com/smacker/go-tree-sitter"
)

// CppExtractor C++语言提取器
type CppExtractor struct {
	BaseExtractor
	queries []string
}

// NewCppExtractor 创建C++语言提取器
func NewCppExtractor() *CppExtractor {
	return &CppExtractor{
		queries: []string{
			"(function_definition) @symbol",
			"(class_specifier) @symbol",
			"(struct_specifier) @symbol",
			"(union_specifier) @symbol",
			"(enum_specifier) @symbol",
			"(namespace_definition) @symbol",
		},
	}
}

// GetQueries 获取C++语言的Tree-sitter查询规则
func (c *CppExtractor) GetQueries() []string {
	return c.queries
}

// ExtractPrototype 提取C++类原型
func (c *CppExtractor) ExtractPrototype(node *sitter.Node, content []byte) string {
	nodeType := node.Type()

	// 对于类声明，只提取类声明部分（不包含类体）
	if nodeType == "class_specifier" || nodeType == "struct_specifier" {
		return c.extractClassPrototype(node, content)
	}

	// 对于命名空间，只提取命名空间声明
	if nodeType == "namespace_definition" {
		return c.extractNamespacePrototype(node, content)
	}

	// 对于函数定义，提取函数签名
	if nodeType == "function_definition" {
		return c.extractFunctionPrototype(node, content, c.IsFunctionBodyNode)
	}

	return c.extractFullNode(node, content)
}

// ExtractMethods 提取C++类内部的方法或命名空间内部的类
func (c *CppExtractor) ExtractMethods(classNode *sitter.Node, content []byte) []models.Symbol {
	var methods []models.Symbol
	nodeType := classNode.Type()

	// 如果是命名空间，提取命名空间内部的类
	if nodeType == "namespace_definition" {
		return c.extractNamespaceClasses(classNode, content)
	}

	// 如果是类，提取类内部的方法
	childCount := int(classNode.ChildCount())
	for i := 0; i < childCount; i++ {
		child := classNode.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找类体
		if childType == "field_declaration_list" {
			// 遍历类体中的方法
			bodyChildCount := int(child.ChildCount())
			for j := 0; j < bodyChildCount; j++ {
				bodyChild := child.Child(j)
				if bodyChild == nil {
					continue
				}

				bodyChildType := bodyChild.Type()
				if bodyChildType == "function_definition" {
					method := c.createMethodSymbol(bodyChild, content)
					methods = append(methods, method)
				}
			}
		}
	}

	return methods
}

// extractNamespaceClasses 提取命名空间内部的类
func (c *CppExtractor) extractNamespaceClasses(namespaceNode *sitter.Node, content []byte) []models.Symbol {
	var classes []models.Symbol

	childCount := int(namespaceNode.ChildCount())
	for i := 0; i < childCount; i++ {
		child := namespaceNode.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找命名空间体
		if childType == "declaration_list" {
			// 遍历命名空间体中的类
			bodyChildCount := int(child.ChildCount())
			for j := 0; j < bodyChildCount; j++ {
				bodyChild := child.Child(j)
				if bodyChild == nil {
					continue
				}

				bodyChildType := bodyChild.Type()
				// C++: class_specifier, struct_specifier
				if bodyChildType == "class_specifier" || bodyChildType == "struct_specifier" {
					class := c.createClassSymbol(bodyChild, content)
					classes = append(classes, class)
				}
			}
		}
	}

	return classes
}

// createClassSymbol 创建类符号
func (c *CppExtractor) createClassSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	prototype := c.extractClassPrototype(node, content)

	return models.Symbol{
		Prototype: prototype,
		Purpose:   extractMultiLineComments(node, content),
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
		Methods:   c.extractClassMethods(node, content),
	}
}

// extractClassMethods 提取类内部的方法
func (c *CppExtractor) extractClassMethods(classNode *sitter.Node, content []byte) []models.Symbol {
	var methods []models.Symbol

	childCount := int(classNode.ChildCount())
	for i := 0; i < childCount; i++ {
		child := classNode.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找类体
		if childType == "field_declaration_list" {
			// 遍历类体中的方法
			bodyChildCount := int(child.ChildCount())
			for j := 0; j < bodyChildCount; j++ {
				bodyChild := child.Child(j)
				if bodyChild == nil {
					continue
				}

				bodyChildType := bodyChild.Type()
				if bodyChildType == "function_definition" {
					method := c.createMethodSymbol(bodyChild, content)
					methods = append(methods, method)
				}
			}
		}
	}

	return methods
}

// IsClassNode 检查是否是类节点或命名空间节点
func (c *CppExtractor) IsClassNode(nodeType string) bool {
	return nodeType == "class_specifier" ||
		nodeType == "struct_specifier" ||
		nodeType == "union_specifier" ||
		nodeType == "enum_specifier" ||
		nodeType == "namespace_definition"
}

// IsFunctionBodyNode 检查是否是函数体节点
func (c *CppExtractor) IsFunctionBodyNode(nodeType string) bool {
	return nodeType == "compound_statement" || nodeType == "block"
}

// IsInsideClass 检查节点是否在类内部
func (c *CppExtractor) IsInsideClass(node *sitter.Node) bool {
	current := node.Parent()
	for current != nil {
		nodeType := current.Type()
		if nodeType == "class_specifier" ||
			nodeType == "struct_specifier" ||
			nodeType == "union_specifier" ||
			nodeType == "enum_specifier" ||
			nodeType == "namespace_definition" {
			return true
		}
		current = current.Parent()
	}
	return false
}

// ExtractComments 提取C++注释
func (c *CppExtractor) ExtractComments(node *sitter.Node, content []byte) string {
	return extractMultiLineComments(node, content)
}

// extractClassPrototype 提取C++类原型
func (c *CppExtractor) extractClassPrototype(node *sitter.Node, content []byte) string {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找类体开始位置
		if childType == "field_declaration_list" || childType == "{" {
			prototype := string(content[node.StartByte():child.StartByte()])
			return strings.TrimSpace(prototype)
		}
	}
	return ""
}

// extractNamespacePrototype 提取C++命名空间原型
func (c *CppExtractor) extractNamespacePrototype(node *sitter.Node, content []byte) string {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找命名空间体开始位置
		if childType == "declaration_list" || childType == "{" {
			prototype := string(content[node.StartByte():child.StartByte()])
			return strings.TrimSpace(prototype)
		}
	}
	return ""
}

// createMethodSymbol 创建方法符号
func (c *CppExtractor) createMethodSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	prototype := c.extractFunctionPrototype(node, content, c.IsFunctionBodyNode)

	return models.Symbol{
		Prototype: prototype,
		Purpose:   c.extractCppComments(node, content),
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}
}

// extractCppComments 提取C++注释
func (c *CppExtractor) extractCppComments(node *sitter.Node, content []byte) string {
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
