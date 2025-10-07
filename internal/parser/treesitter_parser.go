package parser

import (
	"context"
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

const (
	langGo         = "go"
	langJavaScript = "javascript"
	langTypeScript = "typescript"
	langPython     = "python"
)

// 节点类型常量
const (
	nodeClassDeclaration    = "class_declaration"
	nodeClassDefinition     = "class_definition"
	nodeMethodDefinition    = "method_definition"
	nodeFunctionDefinition  = "function_definition"
	nodeExpressionStatement = "expression_statement"
	nodeClassBody           = "class_body"
	nodeBlock               = "block"
	nodeString              = "string"
	nodeUnknown             = "Unknown"
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
	p.initParsers()

	return p, nil
}

// initParsers 初始化语言解析器
func (p *TreeSitterParser) initParsers() {
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
	case langGo:
		language = golang.GetLanguage()
	case langJavaScript, langTypeScript:
		language = javascript.GetLanguage()
	case langPython:
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
		tree, _ := parser.ParseCtx(context.TODO(), nil, content)
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
	case langGo:
		language = golang.GetLanguage()
	case langJavaScript, langTypeScript:
		language = javascript.GetLanguage()
	case langPython:
		language = python.GetLanguage()
	default:
		return symbols
	}

	for _, queryStr := range queries {
		query, err := sitter.NewQuery([]byte(queryStr), language)
		if err != nil {
			continue
		}

		cursor := sitter.NewQueryCursor()
		cursor.Exec(query, node)

		for {
			match, ok := cursor.NextMatch()
			if !ok {
				break
			}

			for _, capture := range match.Captures {
				// 检查函数/方法是否在类内部
				// 如果是类内部的方法，跳过（它们会作为类的 methods 提取）
				nodeType := capture.Node.Type()
				if (nodeType == nodeFunctionDefinition || nodeType == nodeMethodDefinition) && p.isInsideClass(capture.Node) {
					continue
				}

				symbol := p.nodeToSymbol(capture.Node, content)
				symbols = append(symbols, symbol)
			}
		}

		query.Close()
		cursor.Close()
	}

	return symbols
}

// isInsideClass 检查节点是否在类内部
func (p *TreeSitterParser) isInsideClass(node *sitter.Node) bool {
	// 向上遍历父节点，查找是否在类定义中
	current := node.Parent()
	for current != nil {
		nodeType := current.Type()
		if nodeType == nodeClassDefinition || nodeType == nodeClassDeclaration {
			return true
		}
		current = current.Parent()
	}
	return false
}

// nodeToSymbol 将语法树节点转换为符号
func (p *TreeSitterParser) nodeToSymbol(node *sitter.Node, content []byte) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()
	nodeType := node.Type()

	// 提取函数原型（不包括函数体）
	prototype := p.extractPrototype(node, content)

	symbol := models.Symbol{
		Prototype: prototype,
		Purpose:   p.extractPurpose(node, content),
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}

	// 如果是类节点，提取类内部的方法
	if nodeType == nodeClassDeclaration || nodeType == nodeClassDefinition {
		symbol.Methods = p.extractClassMethods(node, content)
	}

	return symbol
}

// createMethodSymbol 创建方法符号
func (p *TreeSitterParser) createMethodSymbol(node *sitter.Node, content []byte, className string) models.Symbol {
	start := node.StartPoint()
	end := node.EndPoint()

	// 提取方法原型
	methodPrototype := p.extractPrototype(node, content)
	// 在方法名前面加上类名
	methodPrototype = p.addClassNameToMethod(className, methodPrototype)

	return models.Symbol{
		Prototype: methodPrototype,
		Purpose:   p.extractPurpose(node, content),
		Range:     []int{int(start.Row) + 1, int(end.Row) + 1},
	}
}

// extractClassMethods 提取类内部的方法
func (p *TreeSitterParser) extractClassMethods(classNode *sitter.Node, content []byte) []models.Symbol {
	var methods []models.Symbol

	// 首先提取类名
	className := p.extractClassName(classNode, content)

	// 遍历类的子节点，查找方法
	childCount := int(classNode.ChildCount())
	for i := 0; i < childCount; i++ {
		child := classNode.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()

		// 对于 JavaScript，方法在 class_body 中
		if childType == nodeClassBody {
			// 遍历 class_body 的子节点查找方法
			bodyChildCount := int(child.ChildCount())
			for j := 0; j < bodyChildCount; j++ {
				bodyChild := child.Child(j)
				if bodyChild == nil {
					continue
				}

				bodyChildType := bodyChild.Type()
				if bodyChildType == nodeMethodDefinition {
					method := p.createMethodSymbol(bodyChild, content, className)
					methods = append(methods, method)
				}
			}
		}

		// 对于 Python，方法可能在类节点下或 block 子节点中
		// 优先检查 block 子节点，避免重复提取
		if childType == nodeBlock {
			// 遍历 block 的子节点查找方法
			blockChildCount := int(child.ChildCount())
			for k := 0; k < blockChildCount; k++ {
				blockChild := child.Child(k)
				if blockChild == nil {
					continue
				}

				blockChildType := blockChild.Type()
				if blockChildType == nodeFunctionDefinition {
					method := p.createMethodSymbol(blockChild, content, className)
					methods = append(methods, method)
				}
			}
		} else if childType == nodeFunctionDefinition {
			// 如果方法直接在类节点下（非 block 中）
			method := p.createMethodSymbol(child, content, className)
			methods = append(methods, method)
		}
	}

	return methods
}

// extractClassName 提取类名
func (p *TreeSitterParser) extractClassName(classNode *sitter.Node, content []byte) string {
	// 遍历类的子节点，查找类名
	childCount := int(classNode.ChildCount())
	for i := 0; i < childCount; i++ {
		child := classNode.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()

		// JavaScript: identifier (类名)
		// Python: identifier (类名)
		if childType == "identifier" {
			className := string(content[child.StartByte():child.EndByte()])
			return className
		}
	}

	return nodeUnknown
}

// addClassNameToMethod 在方法名前面加上类名
func (p *TreeSitterParser) addClassNameToMethod(className, methodPrototype string) string {
	if className == nodeUnknown {
		return methodPrototype
	}

	// 对于不同的语言，使用不同的格式
	// JavaScript: ClassName.methodName()
	// Python: ClassName.methodName()

	// 查找方法名的开始位置
	// 对于 JavaScript: methodName() 或 methodName(params)
	// 对于 Python: def methodName() 或 def methodName(params)

	lines := strings.Split(methodPrototype, "\n")
	if len(lines) == 0 {
		return methodPrototype
	}

	// 处理第一行（通常包含方法名）
	firstLine := lines[0]

	// JavaScript: 直接是方法名
	if strings.Contains(firstLine, "(") && !strings.Contains(firstLine, "def ") {
		// 找到方法名
		parenIndex := strings.Index(firstLine, "(")
		if parenIndex > 0 {
			methodName := strings.TrimSpace(firstLine[:parenIndex])
			// 替换方法名
			newFirstLine := strings.Replace(firstLine, methodName, className+"."+methodName, 1)
			lines[0] = newFirstLine
		}
	}

	// Python: def methodName
	if strings.HasPrefix(firstLine, "def ") {
		// 找到 def 后面的方法名
		defIndex := strings.Index(firstLine, "def ")
		if defIndex >= 0 {
			afterDef := firstLine[defIndex+4:]
			// 找到方法名（到空格或冒号为止）
			spaceIndex := strings.Index(afterDef, " ")
			colonIndex := strings.Index(afterDef, ":")

			var methodName string
			if spaceIndex > 0 && (colonIndex == -1 || spaceIndex < colonIndex) {
				methodName = strings.TrimSpace(afterDef[:spaceIndex])
			} else if colonIndex > 0 {
				methodName = strings.TrimSpace(afterDef[:colonIndex])
			} else {
				methodName = strings.TrimSpace(afterDef)
			}

			// 替换方法名
			newFirstLine := strings.Replace(firstLine, "def "+methodName, "def "+className+"."+methodName, 1)
			lines[0] = newFirstLine
		}
	}

	return strings.Join(lines, "\n")
}

// extractPrototype 提取函数原型（不包括函数体）
func (p *TreeSitterParser) extractPrototype(node *sitter.Node, content []byte) string {
	nodeType := node.Type()

	// 处理类声明
	if prototype := p.extractClassPrototype(node, content, nodeType); prototype != "" {
		return prototype
	}

	// 处理函数体排除
	if prototype := p.extractFunctionPrototype(node, content); prototype != "" {
		return prototype
	}

	// 处理 Python 函数定义
	if prototype := p.extractPythonFunctionPrototype(node, content, nodeType); prototype != "" {
		return prototype
	}

	// 默认返回整个节点
	return p.extractFullNode(node, content)
}

// extractClassPrototype 提取类原型
func (p *TreeSitterParser) extractClassPrototype(node *sitter.Node, content []byte, nodeType string) string {
	// 对于 JavaScript class_declaration
	if nodeType == "class_declaration" {
		return p.extractJavaScriptClassPrototype(node, content)
	}

	// 对于 Python class_definition
	if nodeType == nodeClassDefinition {
		return p.extractPythonClassPrototype(node, content)
	}

	return ""
}

// extractJavaScriptClassPrototype 提取 JavaScript 类原型
func (p *TreeSitterParser) extractJavaScriptClassPrototype(node *sitter.Node, content []byte) string {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		if child.Type() == "class_body" {
			prototype := string(content[node.StartByte():child.StartByte()])
			return strings.TrimSpace(prototype)
		}
	}
	return ""
}

// extractPythonClassPrototype 提取 Python 类原型
func (p *TreeSitterParser) extractPythonClassPrototype(node *sitter.Node, content []byte) string {
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

// extractFunctionPrototype 提取函数原型（排除函数体）
func (p *TreeSitterParser) extractFunctionPrototype(node *sitter.Node, content []byte) string {
	declarationEnd := p.findFunctionBodyEnd(node)
	if declarationEnd > node.StartByte() {
		prototype := string(content[node.StartByte():declarationEnd])
		return strings.TrimSpace(prototype)
	}
	return ""
}

// findFunctionBodyEnd 查找函数体结束位置
func (p *TreeSitterParser) findFunctionBodyEnd(node *sitter.Node) uint32 {
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type()
		if p.isFunctionBodyNode(childType) {
			return child.StartByte()
		}
	}
	return 0
}

// isFunctionBodyNode 检查是否是函数体节点
func (p *TreeSitterParser) isFunctionBodyNode(nodeType string) bool {
	return nodeType == nodeBlock || nodeType == "statement_block" ||
		nodeType == "body" || nodeType == "function_body"
}

// extractPythonFunctionPrototype 提取 Python 函数原型
func (p *TreeSitterParser) extractPythonFunctionPrototype(node *sitter.Node, content []byte, nodeType string) string {
	if nodeType != nodeFunctionDefinition {
		return ""
	}

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

// extractFullNode 提取完整节点内容
func (p *TreeSitterParser) extractFullNode(node *sitter.Node, content []byte) string {
	fullText := string(content[node.StartByte():node.EndByte()])
	return strings.TrimSpace(fullText)
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

// extractPurpose 提取符号的注释/说明
func (p *TreeSitterParser) extractPurpose(node *sitter.Node, content []byte) string {
	// 首先检查节点前的注释
	if purpose := p.extractPrecedingComments(node, content); purpose != "" {
		return purpose
	}

	// 然后检查节点内部的注释（对于类和方法）
	return p.extractInternalComments(node, content)
}

// extractPrecedingComments 提取节点前的注释
func (p *TreeSitterParser) extractPrecedingComments(node *sitter.Node, content []byte) string {
	startPoint := node.StartPoint()
	startRow := int(startPoint.Row)
	lines := strings.Split(string(content), "\n")

	for i := startRow - 1; i >= 0; i-- {
		if i >= len(lines) {
			continue
		}

		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 检查 Python docstring
		if purpose := p.extractDocstringFromLine(line); purpose != "" {
			return purpose
		}

		// 检查单行注释
		if purpose := p.extractCommentFromLine(line); purpose != "" {
			return purpose
		}

		// 如果遇到非注释行，停止查找
		if !p.isCommentLine(line) && line != "" {
			break
		}
	}

	return ""
}

// extractInternalComments 提取节点内部的注释
func (p *TreeSitterParser) extractInternalComments(node *sitter.Node, content []byte) string {
	nodeType := node.Type()
	if !p.isCommentableNode(nodeType) {
		return ""
	}

	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		if purpose := p.extractCommentFromChild(child, content); purpose != "" {
			return purpose
		}
	}

	return ""
}

// extractDocstringFromLine 从行中提取 docstring
func (p *TreeSitterParser) extractDocstringFromLine(line string) string {
	if strings.HasPrefix(line, "\"\"\"") {
		docstring := strings.TrimPrefix(line, "\"\"\"")
		docstring = strings.TrimSuffix(docstring, "\"\"\"")
		docstring = strings.TrimSpace(docstring)
		if docstring != "" {
			return docstring
		}
	}
	return ""
}

// extractCommentFromLine 从行中提取注释
func (p *TreeSitterParser) extractCommentFromLine(line string) string {
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
		comment := strings.TrimSpace(strings.TrimPrefix(line, "#"))
		comment = strings.TrimSpace(strings.TrimPrefix(comment, "//"))
		comment = strings.TrimSpace(strings.TrimPrefix(comment, "/*"))
		comment = strings.TrimSpace(strings.TrimSuffix(comment, "*/"))
		if comment != "" {
			return comment
		}
	}
	return ""
}

// isCommentLine 检查是否是注释行
func (p *TreeSitterParser) isCommentLine(line string) bool {
	return strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*")
}

// isCommentableNode 检查是否是可注释的节点类型
func (p *TreeSitterParser) isCommentableNode(nodeType string) bool {
	return nodeType == nodeClassDefinition || nodeType == nodeFunctionDefinition || nodeType == nodeMethodDefinition
}

// extractCommentFromChild 从子节点中提取注释
func (p *TreeSitterParser) extractCommentFromChild(child *sitter.Node, content []byte) string {
	childType := child.Type()

	// 检查是否是 docstring 节点
	if p.isDocstringNode(childType) {
		if purpose := p.extractDocstringFromNode(child, content); purpose != "" {
			return purpose
		}
	}

	// 如果是 block，检查其子节点
	if childType == nodeBlock {
		return p.extractCommentFromBlock(child, content)
	}

	return ""
}

// isDocstringNode 检查是否是 docstring 节点
func (p *TreeSitterParser) isDocstringNode(nodeType string) bool {
	return nodeType == nodeExpressionStatement || nodeType == nodeString || nodeType == nodeBlock
}

// extractDocstringFromNode 从节点中提取 docstring
func (p *TreeSitterParser) extractDocstringFromNode(node *sitter.Node, content []byte) string {
	childText := string(content[node.StartByte():node.EndByte()])
	childText = strings.TrimSpace(childText)

	// 检查 Python docstring: """...""" 或 '''...'''
	if (strings.HasPrefix(childText, "\"\"\"") && strings.HasSuffix(childText, "\"\"\"")) ||
		(strings.HasPrefix(childText, "'''") && strings.HasSuffix(childText, "'''")) {
		docstring := strings.TrimPrefix(childText, "\"\"\"")
		docstring = strings.TrimSuffix(docstring, "\"\"\"")
		docstring = strings.TrimPrefix(docstring, "'''")
		docstring = strings.TrimSuffix(docstring, "'''")
		docstring = strings.TrimSpace(docstring)
		if docstring != "" {
			return docstring
		}
	}

	return ""
}

// extractCommentFromBlock 从 block 节点中提取注释
func (p *TreeSitterParser) extractCommentFromBlock(blockNode *sitter.Node, content []byte) string {
	blockChildCount := int(blockNode.ChildCount())
	for j := 0; j < blockChildCount; j++ {
		blockChild := blockNode.Child(j)
		if blockChild == nil {
			continue
		}

		blockChildType := blockChild.Type()
		if blockChildType == nodeExpressionStatement || blockChildType == nodeString {
			if purpose := p.extractDocstringFromNode(blockChild, content); purpose != "" {
				return purpose
			}
		}
	}

	return ""
}

// minInt 返回两个整数中的较小者
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
