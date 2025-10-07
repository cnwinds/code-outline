package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cnwinds/CodeCartographer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSimpleParser(t *testing.T) {
	config := models.LanguagesConfig{
		"go": {
			Extensions: []string{".go"},
		},
	}

	parser := NewSimpleParser(config)
	assert.NotNil(t, parser)
	assert.Equal(t, config, parser.languagesConfig)
}

func TestParseGoFile(t *testing.T) {
	parser := NewSimpleParser(getTestConfig())

	// 创建临时测试文件
	tmpFile := createTempFile(t, "test.go", goTestCode)
	defer os.Remove(tmpFile)

	result, err := parser.ParseFile(tmpFile)

	require.NoError(t, err)
	require.NotNil(t, result)

	// 验证符号数量
	assert.GreaterOrEqual(t, len(result.Symbols), 1)

	// 验证第一个符号
	if len(result.Symbols) > 0 {
		symbol := result.Symbols[0]
		assert.NotEmpty(t, symbol.Prototype)
		assert.NotEmpty(t, symbol.Range)
		assert.Equal(t, 2, len(symbol.Range))
	}
}

func TestParseGoSymbols(t *testing.T) {
	parser := NewSimpleParser(getTestConfig())

	lines := []string{
		"package main",
		"",
		"// main 函数",
		"func main() {",
		"    fmt.Println(\"Hello\")",
		"}",
	}

	symbols := parser.parseGoSymbols(lines)

	assert.Equal(t, 1, len(symbols))
	assert.Contains(t, symbols[0].Prototype, "func main()")
	assert.Equal(t, "main 函数", symbols[0].Purpose)
}

func TestParseJSSymbols(t *testing.T) {
	parser := NewSimpleParser(getTestConfig())

	lines := []string{
		"// 用户类",
		"class User {",
		"    constructor(name) {",
		"        this.name = name;",
		"    }",
		"}",
	}

	symbols := parser.parseJSSymbols(lines)

	assert.GreaterOrEqual(t, len(symbols), 1)
}

func TestParsePythonSymbols(t *testing.T) {
	parser := NewSimpleParser(getTestConfig())

	lines := []string{
		"# 用户类",
		"class User:",
		"    def __init__(self, name):",
		"        self.name = name",
	}

	symbols := parser.parsePythonSymbols(lines)

	assert.GreaterOrEqual(t, len(symbols), 1)
}

func TestExtractPurpose(t *testing.T) {
	parser := NewSimpleParser(getTestConfig())

	testCases := []struct {
		name     string
		lines    []string
		lineNum  int
		expected string
	}{
		{
			name: "Go单行注释",
			lines: []string{
				"// 这是一个测试函数",
				"func test() {}",
			},
			lineNum:  1,
			expected: "这是一个测试函数",
		},
		{
			name: "Python注释",
			lines: []string{
				"# 这是Python函数",
				"def test():",
			},
			lineNum:  1,
			expected: "这是Python函数",
		},
		{
			name: "无注释",
			lines: []string{
				"func test() {}",
			},
			lineNum:  0,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parser.extractPurpose(tc.lines, tc.lineNum)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExtractFilePurpose(t *testing.T) {
	parser := NewSimpleParser(getTestConfig())

	// 测试有注释的文件
	contentWithComment := `package main

// 这是一个测试程序
import "fmt"

func main() {
    fmt.Println("Hello")
}`

	purpose := parser.extractFilePurpose(contentWithComment, "go")
	assert.Contains(t, purpose, "测试程序")

	// 测试无注释的文件
	contentWithoutComment := `package main

import "fmt"

func main() {
    fmt.Println("Hello")
}`

	purpose = parser.extractFilePurpose(contentWithoutComment, "go")
	assert.Equal(t, "TODO: Describe the purpose of this file.", purpose)
}

func TestParseFileWithUnsupportedExtension(t *testing.T) {
	parser := NewSimpleParser(getTestConfig())

	// 创建临时测试文件
	tmpFile := createTempFile(t, "test.txt", "some content")
	defer os.Remove(tmpFile)

	_, err := parser.ParseFile(tmpFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不支持的文件类型")
}

func TestParseFileNotFound(t *testing.T) {
	parser := NewSimpleParser(getTestConfig())

	_, err := parser.ParseFile("nonexistent.go")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "获取文件信息失败")
}

// 辅助函数

func getTestConfig() models.LanguagesConfig {
	return models.LanguagesConfig{
		"go": {
			Extensions: []string{".go"},
		},
		"javascript": {
			Extensions: []string{".js", ".jsx"},
		},
		"python": {
			Extensions: []string{".py"},
		},
	}
}

func createTempFile(t *testing.T, name, content string) string {
	tmpFile := filepath.Join(t.TempDir(), name)
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)
	return tmpFile
}

const goTestCode = `package main

// main 函数
func main() {
    println("test")
}

// helper 函数
func helper() string {
    return "help"
}
`
