package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cnwinds/CodeCartographer/internal/models"
)

// FileParser 文件解析器接口
type FileParser interface {
	ParseFile(filePath string) (*models.FileInfo, error)
}

// Scanner 文件扫描器
type Scanner struct {
	parser          FileParser
	excludePatterns []string
}

// NewScanner 创建新的扫描器实例
func NewScanner(parser FileParser, excludePatterns []string) *Scanner {
	return &Scanner{
		parser:          parser,
		excludePatterns: excludePatterns,
	}
}

// ScanProject 扫描整个项目
func (s *Scanner) ScanProject(projectPath string) (map[string]models.FileInfo, []string, error) {
	files := make(map[string]models.FileInfo)
	var techStack []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 用于收集错误的channel
	errorChan := make(chan error, 100)
	var scanErrors []error

	// 启动错误收集goroutine
	go func() {
		for err := range errorChan {
			scanErrors = append(scanErrors, err)
		}
	}()

	// 遍历项目文件
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			// 检查是否是排除的目录
			if s.shouldExclude(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过排除的文件
		if s.shouldExclude(path) {
			return nil
		}

		// 检查文件扩展名
		ext := filepath.Ext(path)
		if ext == "" {
			return nil
		}

		// 获取相对路径
		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			relPath = path
		}

		// 并发解析文件
		wg.Add(1)
		go func(filePath, relativePath string) {
			defer wg.Done()

			fileInfo, err := s.parser.ParseFile(filePath)
			if err != nil {
				errorChan <- fmt.Errorf("解析文件 %s 失败: %w", relativePath, err)
				return
			}

			// 安全地更新结果
			mu.Lock()
			files[relativePath] = *fileInfo

			// 收集技术栈信息
			lang := s.getLanguageFromExtension(ext)
			if lang != "" && !contains(techStack, lang) {
				techStack = append(techStack, lang)
			}
			mu.Unlock()
		}(path, relPath)

		return nil
	})

	// 等待所有goroutine完成
	wg.Wait()
	close(errorChan)

	if err != nil {
		return nil, nil, fmt.Errorf("扫描项目失败: %w", err)
	}

	// 如果有扫描错误，记录但不中断
	if len(scanErrors) > 0 {
		fmt.Printf("警告: 扫描过程中遇到 %d 个错误:\n", len(scanErrors))
		for i, err := range scanErrors {
			if i < 5 { // 只显示前5个错误
				fmt.Printf("  - %v\n", err)
			}
		}
		if len(scanErrors) > 5 {
			fmt.Printf("  ... 还有 %d 个错误\n", len(scanErrors)-5)
		}
	}

	return files, techStack, nil
}

// shouldExclude 检查路径是否应该被排除
func (s *Scanner) shouldExclude(path string) bool {
	// 默认排除模式
	defaultExcludes := []string{
		".git",
		".svn",
		".hg",
		"node_modules",
		"vendor",
		".idea",
		".vscode",
		"__pycache__",
		".DS_Store",
		"*.tmp",
		"*.log",
	}

	// 合并用户指定的排除模式
	allExcludes := append(defaultExcludes, s.excludePatterns...)

	for _, pattern := range allExcludes {
		// 简单的模式匹配
		if strings.Contains(path, pattern) {
			return true
		}

		// 通配符匹配
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}

	return false
}

// getLanguageFromExtension 根据文件扩展名获取语言名称
func (s *Scanner) getLanguageFromExtension(ext string) string {
	languageMap := map[string]string{
		".go":         "Go",
		".js":         "JavaScript",
		".jsx":        "JavaScript",
		".ts":         "TypeScript",
		".tsx":        "TypeScript",
		".py":         "Python",
		".java":       "Java",
		".c":          "C",
		".cpp":        "C++",
		".cc":         "C++",
		".cxx":        "C++",
		".h":          "C/C++",
		".hpp":        "C++",
		".cs":         "C#",
		".php":        "PHP",
		".rb":         "Ruby",
		".rs":         "Rust",
		".swift":      "Swift",
		".kt":         "Kotlin",
		".scala":      "Scala",
		".clj":        "Clojure",
		".hs":         "Haskell",
		".ml":         "OCaml",
		".fs":         "F#",
		".lua":        "Lua",
		".r":          "R",
		".m":          "Objective-C",
		".mm":         "Objective-C++",
		".dart":       "Dart",
		".elm":        "Elm",
		".ex":         "Elixir",
		".exs":        "Elixir",
		".erl":        "Erlang",
		".hrl":        "Erlang",
		".sql":        "SQL",
		".sh":         "Shell",
		".bash":       "Bash",
		".zsh":        "Zsh",
		".fish":       "Fish",
		".ps1":        "PowerShell",
		".html":       "HTML",
		".css":        "CSS",
		".scss":       "SCSS",
		".sass":       "Sass",
		".less":       "Less",
		".xml":        "XML",
		".json":       "JSON",
		".yaml":       "YAML",
		".yml":        "YAML",
		".toml":       "TOML",
		".ini":        "INI",
		".cfg":        "Config",
		".conf":       "Config",
		".md":         "Markdown",
		".tex":        "LaTeX",
		".dockerfile": "Docker",
		".Dockerfile": "Docker",
	}

	if lang, exists := languageMap[ext]; exists {
		return lang
	}

	return ""
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
