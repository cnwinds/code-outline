package parser

import (
	"regexp"
	"strings"

	"github.com/cnwinds/CodeCartographer/internal/models"
	sitter "github.com/smacker/go-tree-sitter"
)

// extractMultiLineComments 提取多行注释（通用方法，支持 /* */ 和 /** */ 格式）
func extractMultiLineComments(node *sitter.Node, content []byte) string {
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
			} else if strings.HasPrefix(line, "/**") || strings.HasPrefix(line, "/*") {
				// 多行注释开始
				comment := strings.TrimSpace(strings.TrimPrefix(line, "/**"))
				comment = strings.TrimSpace(strings.TrimPrefix(comment, "/*"))
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
		if !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/**") && !strings.HasPrefix(line, "/*") && !strings.HasPrefix(line, "*") && line != "" {
			break
		}
	}

	// 如果有收集到的多行注释，合并它们
	if len(commentLines) > 0 {
		return strings.Join(commentLines, " ")
	}

	return ""
}

// extractXMLDocComments 提取XML文档注释（用于C#的 /// 格式）
func extractXMLDocComments(node *sitter.Node, content []byte) string {
	startPoint := node.StartPoint()
	startRow := int(startPoint.Row)
	lines := strings.Split(string(content), "\n")

	var commentLines []string

	for i := startRow - 1; i >= 0; i-- {
		if i >= len(lines) {
			continue
		}

		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 检查XML文档注释
		if strings.HasPrefix(line, "///") {
			comment := strings.TrimSpace(strings.TrimPrefix(line, "///"))
			// 移除XML标签
			comment = strings.TrimSpace(comment)

			// 跳过纯标签行
			if comment == "<summary>" || comment == "</summary>" || comment == "<param>" || comment == "</param>" || comment == "<returns>" || comment == "</returns>" {
				continue
			}

			if strings.HasPrefix(comment, "<summary>") {
				comment = strings.TrimPrefix(comment, "<summary>")
			}
			if strings.HasSuffix(comment, "</summary>") {
				comment = strings.TrimSuffix(comment, "</summary>")
			}
			if strings.HasPrefix(comment, "<param") {
				// 提取参数描述
				if idx := strings.Index(comment, ">"); idx > 0 {
					comment = comment[idx+1:]
				}
			}
			if strings.HasPrefix(comment, "<returns>") {
				comment = strings.TrimPrefix(comment, "<returns>")
			}
			if strings.HasSuffix(comment, "</returns>") {
				comment = strings.TrimSuffix(comment, "</returns>")
			}
			if strings.HasSuffix(comment, "</param>") {
				comment = strings.TrimSuffix(comment, "</param>")
			}

			comment = strings.TrimSpace(comment)
			if comment != "" {
				commentLines = append([]string{comment}, commentLines...)
			}
			continue
		}

		// 如果遇到非注释行，停止查找
		if !strings.HasPrefix(line, "///") && line != "" {
			break
		}
	}

	// 如果有收集到的XML文档注释，合并它们
	if len(commentLines) > 0 {
		return strings.Join(commentLines, " ")
	}

	return ""
}

// LanguageExtractor 语言提取器接口
type LanguageExtractor interface {
	// GetQueries 获取Tree-sitter查询规则
	GetQueries() []string

	// ExtractPrototype 提取函数/类原型（不包含函数体/类体）
	ExtractPrototype(node *sitter.Node, content []byte) string

	// ExtractMethods 提取类内部的方法
	ExtractMethods(classNode *sitter.Node, content []byte) []models.Symbol

	// IsClassNode 检查是否是类节点
	IsClassNode(nodeType string) bool

	// IsFunctionBodyNode 检查是否是函数体节点
	IsFunctionBodyNode(nodeType string) bool

	// IsInsideClass 检查节点是否在类内部
	IsInsideClass(node *sitter.Node) bool

	// ExtractComments 提取注释
	ExtractComments(node *sitter.Node, content []byte) string
}

// BaseExtractor 基础提取器，提供通用功能
type BaseExtractor struct{}

// cleanText 清理文本，移除多余的空格和换行
func (b *BaseExtractor) cleanText(text string) string {
	text = strings.TrimSpace(text)

	// 移除多余的空格（多个连续空格替换为单个空格）
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	// 移除参数列表中的换行和缩进
	text = strings.ReplaceAll(text, ",\n", ", ")
	text = strings.ReplaceAll(text, "(\n", "(")
	text = strings.ReplaceAll(text, "\n)", ")")

	// 移除函数参数中的换行和缩进
	paramRegex := regexp.MustCompile(`\s*\n\s*`)
	text = paramRegex.ReplaceAllString(text, " ")

	// 清理模板参数中的换行
	text = strings.ReplaceAll(text, "<\n", "<")
	text = strings.ReplaceAll(text, "\n>", ">")

	// 最终清理：移除多余空格
	text = strings.TrimSpace(text)

	return text
}

// extractFullNode 提取完整节点内容
func (b *BaseExtractor) extractFullNode(node *sitter.Node, content []byte) string {
	fullText := string(content[node.StartByte():node.EndByte()])
	return b.cleanText(fullText)
}

// findFunctionBodyEnd 查找函数体结束位置
func (b *BaseExtractor) findFunctionBodyEnd(node *sitter.Node, isFunctionBodyNode func(string) bool) uint32 {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		if isFunctionBodyNode(childType) {
			return child.StartByte()
		}
	}
	return 0
}

// extractFunctionPrototype 提取函数原型（通用实现）
func (b *BaseExtractor) extractFunctionPrototype(node *sitter.Node, content []byte, isFunctionBodyNode func(string) bool) string {
	declarationEnd := b.findFunctionBodyEnd(node, isFunctionBodyNode)
	if declarationEnd > node.StartByte() {
		prototype := string(content[node.StartByte():declarationEnd])
		return b.cleanText(prototype)
	}
	return ""
}
