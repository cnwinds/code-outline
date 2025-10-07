package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cnwinds/CodeCartographer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadLanguagesConfig(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_languages.json")

	// 写入测试配置

	configData := `{
		"go": {
			"extensions": [".go"],
			"queries": {
				"top_level_symbols": ["(function_declaration) @symbol"]
			}
		}
	}`

	err := os.WriteFile(configPath, []byte(configData), 0644)
	require.NoError(t, err)

	// 测试加载配置
	config, err := LoadLanguagesConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Contains(t, config, "go")
}

func TestLoadLanguagesConfigFileNotExist(t *testing.T) {
	// 测试不存在的配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.json")

	config, err := LoadLanguagesConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, config)

	// 验证默认配置被创建
	assert.Contains(t, config, "go")
	assert.Contains(t, config, "javascript")
	assert.Contains(t, config, "python")
}

func TestLoadLanguagesConfigInvalidJSON(t *testing.T) {
	// 创建无效的JSON文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	invalidJSON := `{ invalid json }`
	err := os.WriteFile(configPath, []byte(invalidJSON), 0644)
	require.NoError(t, err)

	// 测试加载无效配置
	_, err = LoadLanguagesConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "解析配置文件失败")
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

func TestCreateDefaultLanguagesConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "default_languages.json")

	config, err := createDefaultLanguagesConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, config)

	// 验证默认配置包含预期语言
	assert.Contains(t, config, "go")
	assert.Contains(t, config, "javascript")
	assert.Contains(t, config, "python")

	// 验证配置文件被创建
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// 验证配置文件内容
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "go")
	assert.Contains(t, string(content), "javascript")
}

func TestCreateDefaultLanguagesConfigDirectoryCreation(t *testing.T) {
	// 测试在深层目录中创建配置文件
	tmpDir := t.TempDir()
	deepDir := filepath.Join(tmpDir, "deep", "nested", "directory")
	configPath := filepath.Join(deepDir, "languages.json")

	config, err := createDefaultLanguagesConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, config)

	// 验证目录被创建
	_, err = os.Stat(deepDir)
	assert.NoError(t, err)

	// 验证配置文件被创建
	_, err = os.Stat(configPath)
	assert.NoError(t, err)
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
