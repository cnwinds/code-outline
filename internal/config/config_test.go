package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cnwinds/code-outline/internal/models"
)

func TestGetDefaultLanguagesConfig(t *testing.T) {
	// 测试获取默认语言配置
	config := GetDefaultLanguagesConfig()
	require.NotNil(t, config)
	assert.Contains(t, config, "go")
	assert.Contains(t, config, "javascript")
	assert.Contains(t, config, "python")
	assert.Contains(t, config, "java")
	assert.Contains(t, config, "csharp")
	assert.Contains(t, config, "cpp")
	assert.Contains(t, config, "c")
	assert.Contains(t, config, "rust")
	assert.Contains(t, config, "typescript")
}

func TestGetDefaultLanguagesConfigContent(t *testing.T) {
	// 测试默认配置的内容
	config := GetDefaultLanguagesConfig()

	// 验证Go配置
	goConfig, exists := config["go"]
	assert.True(t, exists)
	assert.Equal(t, []string{".go"}, goConfig.Extensions)

	// 验证JavaScript配置
	jsConfig, exists := config["javascript"]
	assert.True(t, exists)
	assert.Equal(t, []string{".js", ".jsx"}, jsConfig.Extensions)

	// 验证Python配置
	pyConfig, exists := config["python"]
	assert.True(t, exists)
	assert.Equal(t, []string{".py"}, pyConfig.Extensions)
}

func TestGetDefaultLanguagesConfigAllLanguages(t *testing.T) {
	// 测试所有支持的语言都被包含
	config := GetDefaultLanguagesConfig()

	expectedLanguages := []string{
		"go", "java", "csharp", "cpp", "c", "rust",
		"javascript", "typescript", "python",
	}

	for _, lang := range expectedLanguages {
		assert.Contains(t, config, lang, "语言 %s 应该被包含在默认配置中", lang)
	}
}

func TestGetLanguageByExtension(t *testing.T) {
	config := models.LanguagesConfig{
		"go": {
			Extensions: []string{".go"},
		},
		"javascript": {
			Extensions: []string{".js", ".jsx"},
		},
	}

	// 测试存在的扩展名
	langName, langConfig, found := GetLanguageByExtension(config, ".go")
	assert.True(t, found)
	assert.Equal(t, "go", langName)
	assert.Equal(t, []string{".go"}, langConfig.Extensions)

	// 测试不存在的扩展名
	_, _, found = GetLanguageByExtension(config, ".py")
	assert.False(t, found)

	// 测试JavaScript扩展名
	langName, langConfig, found = GetLanguageByExtension(config, ".js")
	assert.True(t, found)
	assert.Equal(t, "javascript", langName)
	assert.Equal(t, []string{".js", ".jsx"}, langConfig.Extensions)
}

func TestGetDefaultLanguagesConfigExtensions(t *testing.T) {
	// 测试各种语言的扩展名配置
	config := GetDefaultLanguagesConfig()

	// 测试Go扩展名
	goConfig := config["go"]
	assert.Equal(t, []string{".go"}, goConfig.Extensions)

	// 测试JavaScript扩展名
	jsConfig := config["javascript"]
	assert.Equal(t, []string{".js", ".jsx"}, jsConfig.Extensions)

	// 测试C++扩展名
	cppConfig := config["cpp"]
	assert.Equal(t, []string{".cpp", ".hpp", ".cc", ".cxx"}, cppConfig.Extensions)
}

func TestGetDefaultLanguagesConfigConsistency(t *testing.T) {
	// 测试配置的一致性
	config1 := GetDefaultLanguagesConfig()
	config2 := GetDefaultLanguagesConfig()

	// 验证多次调用返回相同的配置
	assert.Equal(t, config1, config2)

	// 验证配置不为空
	assert.NotEmpty(t, config1)
	assert.NotEmpty(t, config2)
}

func TestConfigStruct(t *testing.T) {
	// 测试Config结构体
	config := Config{
		Languages: models.LanguagesConfig{
			"go": {
				Extensions: []string{".go"},
			},
		},
		Output:      "output.json",
		Exclude:     []string{"node_modules"},
		ProjectPath: "/path/to/project",
	}

	assert.NotNil(t, config.Languages)
	assert.Equal(t, "output.json", config.Output)
	assert.Equal(t, []string{"node_modules"}, config.Exclude)
	assert.Equal(t, "/path/to/project", config.ProjectPath)
}
