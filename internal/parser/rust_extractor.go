package parser

import (
	"strings"

	"github.com/cnwinds/CodeCartographer/internal/models"
	sitter "github.com/smacker/go-tree-sitter"
)

// RustExtractor Rust语言提取器
type RustExtractor struct {
	BaseExtractor
	queries []string
}

// NewRustExtractor 创建Rust语言提取器
func NewRustExtractor() *RustExtractor {
	return &RustExtractor{
		queries: []string{
			"(function_item) @symbol",
			"(struct_item) @symbol",
			"(enum_item) @symbol",
			"(trait_item) @symbol",
			"(impl_item) @symbol",
		},
	}
}

// GetQueries 获取Rust语言的Tree-sitter查询规则
func (r *RustExtractor) GetQueries() []string {
	return r.queries
}

// ExtractPrototype 提取Rust函数/结构体原型
func (r *RustExtractor) ExtractPrototype(node *sitter.Node, content []byte) string {
	nodeType := node.Type()

	// 对于函数定义，提取函数签名（不包含函数体）
	if nodeType == "function_item" {
		return r.extractFunctionPrototype(node, content, r.IsFunctionBodyNode)
	}

	// 对于结构体，返回完整内容
	if nodeType == "struct_item" {
		return r.extractStructPrototype(node, content)
	}

	// 对于impl块，返回完整内容
	if nodeType == "impl_item" {
		return r.extractImplPrototype(node, content)
	}

	return r.extractFullNode(node, content)
}

// ExtractMethods 提取Rust impl块内部的方法
func (r *RustExtractor) ExtractMethods(implNode *sitter.Node, content []byte) []models.Symbol {
	var methods []models.Symbol

	childCount := int(implNode.ChildCount())
	for i := 0; i < childCount; i++ {
		child := implNode.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找impl体
		if childType == "declaration_list" || childType == "{" {
			// 遍历impl体中的方法
			bodyChildCount := int(child.ChildCount())
			for j := 0; j < bodyChildCount; j++ {
				bodyChild := child.Child(j)
				if bodyChild == nil {
					continue
				}

				bodyChildType := bodyChild.Type()
				if bodyChildType == "function_item" {
					method := r.createMethodSymbol(bodyChild, content)
					methods = append(methods, method)
				}
			}
		}
	}

	return methods
}

// IsClassNode 检查是否是类节点
func (r *RustExtractor) IsClassNode(nodeType string) bool {
	return nodeType == "struct_item" ||
		nodeType == "enum_item" ||
		nodeType == "trait_item" ||
		nodeType == "impl_item"
}

// IsFunctionBodyNode 检查是否是函数体节点
func (r *RustExtractor) IsFunctionBodyNode(nodeType string) bool {
	return nodeType == "block" || nodeType == "function_body"
}

// IsInsideClass 检查节点是否在类内部
func (r *RustExtractor) IsInsideClass(node *sitter.Node) bool {
	current := node.Parent()
	for current != nil {
		nodeType := current.Type()
		if nodeType == "struct_item" ||
			nodeType == "enum_item" ||
			nodeType == "trait_item" ||
			nodeType == "impl_item" {
			return true
		}
		current = current.Parent()
	}
	return false
}

// ExtractComments 提取Rust注释
func (r *RustExtractor) ExtractComments(node *sitter.Node, content []byte) string {
	return r.extractRustComments(node, content)
}

// extractStructPrototype 提取Rust结构体原型
func (r *RustExtractor) extractStructPrototype(node *sitter.Node, content []byte) string {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找结构体体开始位置
		if childType == "field_declaration_list" || childType == "{" {
			prototype := string(content[node.StartByte():child.StartByte()])
			return strings.TrimSpace(prototype)
		}
	}
	return ""
}

// extractImplPrototype 提取Rust impl块原型
func (r *RustExtractor) extractImplPrototype(node *sitter.Node, content []byte) string {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		// 查找impl体开始位置
		if childType == "declaration_list" || childType == "{" {
			prototype := string(content[node.StartByte():child.StartByte()])
			return strings.TrimSpace(prototype)
		}
	}
	return ""
}

// createMethodSymbol 创建方法符号
func (r *RustExtractor) createMethodSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	prototype := r.extractFunctionPrototype(node, content, r.IsFunctionBodyNode)

	return models.Symbol{
		Prototype: prototype,
		Purpose:   r.extractRustComments(node, content),
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}
}

// extractRustComments 提取Rust注释（支持文档注释）
func (r *RustExtractor) extractRustComments(node *sitter.Node, content []byte) string {
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

		// 检查文档注释结束
		if strings.HasSuffix(line, "*/") {
			inMultiLineComment = true
			comment := strings.TrimSpace(strings.TrimSuffix(line, "*/"))
			if comment != "" {
				commentLines = append([]string{comment}, commentLines...)
			}
			continue
		}

		// 检查文档注释中间行
		if inMultiLineComment {
			if strings.HasPrefix(line, "///") {
				comment := strings.TrimSpace(strings.TrimPrefix(line, "///"))
				if comment != "" {
					commentLines = append([]string{comment}, commentLines...)
				}
				continue
			} else if strings.HasPrefix(line, "/**") {
				// 文档注释开始
				comment := strings.TrimSpace(strings.TrimPrefix(line, "/**"))
				if comment != "" {
					commentLines = append([]string{comment}, commentLines...)
				}
				inMultiLineComment = false
				break
			}
		}

		// 检查Rust文档注释 ///
		if strings.HasPrefix(line, "///") {
			comment := strings.TrimSpace(strings.TrimPrefix(line, "///"))
			// 跳过空行和标题行
			if comment != "" && !strings.HasPrefix(comment, "#") {
				commentLines = append([]string{comment}, commentLines...)
			}
			continue
		}

		// 检查普通注释 //
		if strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "///") {
			comment := strings.TrimSpace(strings.TrimPrefix(line, "//"))
			if comment != "" {
				commentLines = append([]string{comment}, commentLines...)
			}
			continue
		}

		// 如果遇到非注释行，停止查找
		if !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/**") && !strings.HasPrefix(line, "///") && line != "" {
			break
		}
	}

	// 如果有收集到的多行注释，合并它们
	if len(commentLines) > 0 {
		return strings.Join(commentLines, " ")
	}

	return ""
}
