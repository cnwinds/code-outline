package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yourusername/CodeCartographer/internal/config"
	"github.com/yourusername/CodeCartographer/internal/models"
)

// SimpleParser 简单的基于正则表达式的解析器
// 用于在Tree-sitter语法文件不可用时提供基本功能
type SimpleParser struct {
	languagesConfig models.LanguagesConfig
}

// NewSimpleParser 创建新的简单解析器实例
func NewSimpleParser(languagesConfig models.LanguagesConfig) *SimpleParser {
	return &SimpleParser{
		languagesConfig: languagesConfig,
	}
}

// ParseFile 解析单个文件
func (p *SimpleParser) ParseFile(filePath string) (*models.FileInfo, error) {
	// 获取文件扩展名
	ext := filepath.Ext(filePath)
	
	// 查找对应的语言配置
	langName, _, found := config.GetLanguageByExtension(p.languagesConfig, ext)
	if !found {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 根据语言类型解析符号
	symbols, err := p.parseSymbolsByLanguage(string(content), langName)
	if err != nil {
		return nil, fmt.Errorf("解析符号失败: %v", err)
	}

	return &models.FileInfo{
		Purpose: p.extractFilePurpose(string(content), langName),
		Symbols: symbols,
	}, nil
}

// parseSymbolsByLanguage 根据语言类型解析符号
func (p *SimpleParser) parseSymbolsByLanguage(content, language string) ([]models.Symbol, error) {
	var symbols []models.Symbol
	lines := strings.Split(content, "\n")

	switch language {
	case "go":
		symbols = p.parseGoSymbols(lines)
	case "javascript", "typescript":
		symbols = p.parseJSSymbols(lines)
	case "python":
		symbols = p.parsePythonSymbols(lines)
	case "java":
		symbols = p.parseJavaSymbols(lines)
	default:
		// 通用解析逻辑
		symbols = p.parseGenericSymbols(lines)
	}

	return symbols, nil
}

// parseGoSymbols 解析Go语言符号
func (p *SimpleParser) parseGoSymbols(lines []string) []models.Symbol {
	var symbols []models.Symbol
	
	// 正则表达式模式
	patterns := map[string]*regexp.Regexp{
		"function": regexp.MustCompile(`^func\s+(\w+)?(\([^)]*\))?\s*(\([^)]*\))?\s*(\w.*)?{?`),
		"method":   regexp.MustCompile(`^func\s+\([^)]+\)\s+(\w+)\s*\([^)]*\)\s*(\w.*)?{?`),
		"type":     regexp.MustCompile(`^type\s+(\w+)\s+(struct|interface|[^{]+)`),
		"const":    regexp.MustCompile(`^const\s+(\w+).*`),
		"var":      regexp.MustCompile(`^var\s+(\w+).*`),
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		for symbolType, pattern := range patterns {
			if matches := pattern.FindStringSubmatch(trimmed); len(matches) > 1 {
				purpose := p.extractPurpose(lines, i)
				
				symbol := models.Symbol{
					Prototype: trimmed,
					Purpose:   purpose,
					Range:     []int{i + 1, i + 1}, // 简化的行号范围
				}

				// 对于struct和interface，尝试提取body
				if symbolType == "type" && (strings.Contains(trimmed, "struct") || strings.Contains(trimmed, "interface")) {
					body, endLine := p.extractGoTypeBody(lines, i)
					symbol.Body = body
					symbol.Range[1] = endLine
				}

				symbols = append(symbols, symbol)
				break
			}
		}
	}

	return symbols
}

// parseJSSymbols 解析JavaScript/TypeScript符号
func (p *SimpleParser) parseJSSymbols(lines []string) []models.Symbol {
	var symbols []models.Symbol
	
	patterns := map[string]*regexp.Regexp{
		"function":  regexp.MustCompile(`^(export\s+)?(async\s+)?function\s+(\w+)`),
		"arrow":     regexp.MustCompile(`^(const|let|var)\s+(\w+)\s*=\s*(async\s+)?\([^)]*\)\s*=>`),
		"class":     regexp.MustCompile(`^(export\s+)?(abstract\s+)?class\s+(\w+)`),
		"interface": regexp.MustCompile(`^(export\s+)?interface\s+(\w+)`),
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		for _, pattern := range patterns {
			if matches := pattern.FindStringSubmatch(trimmed); len(matches) > 1 {
				purpose := p.extractPurpose(lines, i)
				
				symbol := models.Symbol{
					Prototype: trimmed,
					Purpose:   purpose,
					Range:     []int{i + 1, i + 1},
				}

				symbols = append(symbols, symbol)
				break
			}
		}
	}

	return symbols
}

// parsePythonSymbols 解析Python符号
func (p *SimpleParser) parsePythonSymbols(lines []string) []models.Symbol {
	var symbols []models.Symbol
	
	patterns := map[string]*regexp.Regexp{
		"function": regexp.MustCompile(`^def\s+(\w+)\s*\(`),
		"class":    regexp.MustCompile(`^class\s+(\w+).*:`),
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		for _, pattern := range patterns {
			if matches := pattern.FindStringSubmatch(trimmed); len(matches) > 1 {
				purpose := p.extractPurpose(lines, i)
				
				symbol := models.Symbol{
					Prototype: trimmed,
					Purpose:   purpose,
					Range:     []int{i + 1, i + 1},
				}

				symbols = append(symbols, symbol)
				break
			}
		}
	}

	return symbols
}

// parseJavaSymbols 解析Java符号
func (p *SimpleParser) parseJavaSymbols(lines []string) []models.Symbol {
	var symbols []models.Symbol
	
	patterns := map[string]*regexp.Regexp{
		"method": regexp.MustCompile(`^\s*(public|private|protected).*\s+(\w+)\s*\(`),
		"class":  regexp.MustCompile(`^\s*(public\s+)?(abstract\s+)?class\s+(\w+)`),
		"interface": regexp.MustCompile(`^\s*(public\s+)?interface\s+(\w+)`),
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		for _, pattern := range patterns {
			if matches := pattern.FindStringSubmatch(trimmed); len(matches) > 1 {
				purpose := p.extractPurpose(lines, i)
				
				symbol := models.Symbol{
					Prototype: trimmed,
					Purpose:   purpose,
					Range:     []int{i + 1, i + 1},
				}

				symbols = append(symbols, symbol)
				break
			}
		}
	}

	return symbols
}

// parseGenericSymbols 通用符号解析
func (p *SimpleParser) parseGenericSymbols(lines []string) []models.Symbol {
	var symbols []models.Symbol
	
	// 简单的通用模式：查找看起来像函数或类定义的行
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^\w+.*\([^)]*\).*{?`),  // 函数模式
		regexp.MustCompile(`^(class|struct|interface)\s+\w+`), // 类型定义模式
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || len(trimmed) < 3 {
			continue
		}

		for _, pattern := range patterns {
			if pattern.MatchString(trimmed) {
				symbol := models.Symbol{
					Prototype: trimmed,
					Purpose:   "Generic symbol detected",
					Range:     []int{i + 1, i + 1},
				}
				symbols = append(symbols, symbol)
				break
			}
		}
	}

	return symbols
}

// extractPurpose 提取符号的用途说明（从前面的注释）
func (p *SimpleParser) extractPurpose(lines []string, currentLine int) string {
	// 向上查找注释
	for i := currentLine - 1; i >= 0 && i >= currentLine-3; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		
		// Go注释
		if strings.HasPrefix(line, "//") {
			return strings.TrimSpace(strings.TrimPrefix(line, "//"))
		}
		
		// Python注释
		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(strings.TrimPrefix(line, "#"))
		}
		
		// 其他语言的注释暂时不处理
		break
	}
	
	return ""
}

// extractGoTypeBody 提取Go类型的主体内容
func (p *SimpleParser) extractGoTypeBody(lines []string, startLine int) (string, int) {
	var bodyLines []string
	braceCount := 0
	started := false
	
	for i := startLine; i < len(lines); i++ {
		line := lines[i]
		
		// 计算大括号
		for _, char := range line {
			if char == '{' {
				braceCount++
				started = true
			} else if char == '}' {
				braceCount--
			}
		}
		
		// 如果已经开始并且找到了内容
		if started && braceCount > 0 {
			// 提取大括号内的内容
			if strings.Contains(line, "{") {
				parts := strings.Split(line, "{")
				if len(parts) > 1 {
					bodyLines = append(bodyLines, parts[1])
				}
			} else {
				bodyLines = append(bodyLines, line)
			}
		}
		
		// 如果大括号闭合，结束
		if started && braceCount == 0 {
			return strings.Join(bodyLines, "\n"), i + 1
		}
	}
	
	return strings.Join(bodyLines, "\n"), len(lines)
}

// extractFilePurpose 提取文件的用途说明
func (p *SimpleParser) extractFilePurpose(content, language string) string {
	lines := strings.Split(content, "\n")
	
	// 查找文件顶部的注释
	for _, line := range lines[:min(10, len(lines))] {
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

// min 返回两个整数中的较小者
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
