package parser

// ExtractorFactory 提取器工厂
type ExtractorFactory struct{}

// NewExtractorFactory 创建提取器工厂
func NewExtractorFactory() *ExtractorFactory {
	return &ExtractorFactory{}
}

// GetExtractor 根据语言获取对应的提取器
func (f *ExtractorFactory) GetExtractor(language string) LanguageExtractor {
	switch language {
	case "go":
		return NewGoExtractor()
	case "java":
		return NewJavaExtractor()
	case "csharp":
		return NewCSharpExtractor()
	case "cpp":
		return NewCppExtractor()
	case "c":
		return NewCExtractor()
	case "rust":
		return NewRustExtractor()
	case "javascript", "typescript":
		return NewJSExtractor()
	case "python":
		return NewPythonExtractor()
	default:
		// 默认返回Go提取器
		return NewGoExtractor()
	}
}
