package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cnwinds/CodeCartographer/internal/models"
)

// mockParser 用于测试的模拟解析器
type mockParser struct{}

func (m *mockParser) ParseFile(filePath string) (*models.FileInfo, error) {
	// 模拟解析结果
	return &models.FileInfo{
		Purpose: "Mock file",
		Symbols: []models.Symbol{
			{
				Prototype: "func mock()",
				Purpose:   "Mock function",
				Range:     []int{1, 1},
			},
		},
		LastModified: "2025-10-07T00:00:00Z",
		FileSize:     100,
	}, nil
}

func TestNewScanner(t *testing.T) {
	parser := &mockParser{}
	excludePatterns := []string{"test", "temp"}

	scanner := NewScanner(parser, excludePatterns)
	assert.NotNil(t, scanner)
	assert.Equal(t, parser, scanner.parser)
	assert.Equal(t, excludePatterns, scanner.excludePatterns)
}

func TestScanProject(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()

	// 创建测试文件
	createTestFile(t, tmpDir, "main.go", goTestCode)
	createTestFile(t, tmpDir, "helper.js", jsTestCode)
	createTestFile(t, tmpDir, "README.md", "# Test")

	// 创建子目录
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("创建目录失败: %v", err)
	}
	createTestFile(t, subDir, "utils.go", goTestCode)

	// 创建解析器
	parser := &mockParser{}
	scanner := NewScanner(parser, nil)

	// 扫描项目
	files, techStack, err := scanner.ScanProject(tmpDir)

	require.NoError(t, err)
	assert.NotNil(t, files)
	assert.NotNil(t, techStack)

	// 验证结果
	assert.GreaterOrEqual(t, len(files), 2) // 至少包含 .go 和 .js 文件
	assert.Contains(t, techStack, "Go")
	assert.Contains(t, techStack, "JavaScript")
}

func TestScanProjectWithExcludePatterns(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()

	// 创建测试文件
	createTestFile(t, tmpDir, "main.go", goTestCode)
	createTestFile(t, tmpDir, "test.go", goTestCode)
	createTestFile(t, tmpDir, "temp.js", jsTestCode)

	// 创建解析器，排除包含 "test" 和 "temp" 的文件
	parser := &mockParser{}
	excludePatterns := []string{"test", "temp"}
	scanner := NewScanner(parser, excludePatterns)

	// 扫描项目
	files, _, err := scanner.ScanProject(tmpDir)

	require.NoError(t, err)
	assert.NotNil(t, files)

	// 验证排除的文件不在结果中
	for filePath := range files {
		assert.NotContains(t, filePath, "test")
		assert.NotContains(t, filePath, "temp")
	}
}

func TestShouldExclude(t *testing.T) {
	scanner := &Scanner{}

	// 测试默认排除模式
	testCases := []struct {
		path     string
		expected bool
	}{
		{".git/config", true},
		{"node_modules/package", true},
		{"vendor/dependency", true},
		{"main.go", false},
		{"src/main.go", false},
		{".DS_Store", true},
		{"temp.log", true},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := scanner.shouldExclude(tc.path)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestShouldExcludeWithCustomPatterns(t *testing.T) {
	excludePatterns := []string{"test", "temp"}
	scanner := NewScanner(&mockParser{}, excludePatterns)

	testCases := []struct {
		path     string
		expected bool
	}{
		{"test.go", true},
		{"temp.js", true},
		{"main.go", false},
		{"src/main.go", false},
		{"latest.go", true}, // 包含 "test"
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := scanner.shouldExclude(tc.path)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetLanguageFromExtension(t *testing.T) {
	scanner := &Scanner{}

	testCases := []struct {
		ext      string
		expected string
	}{
		{".go", "Go"},
		{".js", "JavaScript"},
		{".jsx", "JavaScript"},
		{".ts", "TypeScript"},
		{".py", "Python"},
		{".java", "Java"},
		{".cs", "C#"},
		{".rs", "Rust"},
		{".cpp", "C++"},
		{".c", "C"},
		{".txt", ""},
		{"", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.ext, func(t *testing.T) {
			result := scanner.getLanguageFromExtension(tc.ext)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"go", "js", "py"}, "go", true},
		{[]string{"go", "js", "py"}, "js", true},
		{[]string{"go", "js", "py"}, "java", false},
		{[]string{}, "go", false},
		{nil, "go", false},
	}

	for _, tc := range testCases {
		t.Run(tc.item, func(t *testing.T) {
			result := contains(tc.slice, tc.item)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestScanProjectEmptyDirectory(t *testing.T) {
	// 创建空目录
	tmpDir := t.TempDir()

	parser := &mockParser{}
	scanner := NewScanner(parser, nil)

	files, techStack, err := scanner.ScanProject(tmpDir)

	require.NoError(t, err)
	assert.Empty(t, files)
	assert.Empty(t, techStack)
}

func TestScanProjectWithErrors(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()

	// 创建测试文件
	createTestFile(t, tmpDir, "main.go", goTestCode)

	// 创建会失败的解析器
	failingParser := &failingParser{}
	scanner := NewScanner(failingParser, nil)

	// 扫描项目（应该处理错误但不中断）
	files, techStack, err := scanner.ScanProject(tmpDir)

	// 应该成功，但文件可能为空
	require.NoError(t, err)

	// 验证结果：由于解析器失败，files 应该为空
	assert.Empty(t, files)
	assert.Empty(t, techStack)
}

// 辅助函数

func createTestFile(t *testing.T, dir, name, content string) {
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0600)
	require.NoError(t, err)
}

// failingParser 总是返回错误的解析器
type failingParser struct{}

func (f *failingParser) ParseFile(filePath string) (*models.FileInfo, error) {
	return nil, assert.AnError
}

const goTestCode = `package main

func main() {
    println("test")
}
`

const jsTestCode = `function test() {
    console.log("test");
}
`
