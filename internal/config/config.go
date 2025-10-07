package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cnwinds/CodeCartographer/internal/models"
)

// Config 表示应用程序配置
type Config struct {
	Languages   models.LanguagesConfig
	Output      string
	Exclude     []string
	ProjectPath string
}

// LoadLanguagesConfig 从指定路径加载语言配置文件
func LoadLanguagesConfig(configPath string) (models.LanguagesConfig, error) {
	// 如果没有指定配置路径，使用默认路径
	if configPath == "" {
		configPath = "languages.json"
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果配置文件不存在，创建默认配置
		return createDefaultLanguagesConfig(configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config models.LanguagesConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return config, nil
}

// createDefaultLanguagesConfig 创建默认的语言配置文件
func createDefaultLanguagesConfig(configPath string) (models.LanguagesConfig, error) {
	defaultConfig := models.LanguagesConfig{
		"go": {
			Extensions: []string{".go"},
			Queries: models.Queries{
				TopLevelSymbols: []string{
					"(function_declaration) @symbol",
					"(method_declaration) @symbol",
					"(type_declaration) @symbol",
					"(const_declaration) @symbol",
					"(var_declaration) @symbol",
				},
				ContainerBody:    "(block) @body | (struct_type) @body | (interface_type) @body",
				ContainerMethods: "(method_declaration) @method",
			},
		},
		"javascript": {
			Extensions: []string{".js", ".jsx"},
			Queries: models.Queries{
				TopLevelSymbols: []string{
					"(function_declaration) @symbol",
					"(arrow_function) @symbol",
					"(class_declaration) @symbol",
					"(const_declaration) @symbol",
					"(let_declaration) @symbol",
					"(var_declaration) @symbol",
				},
				ContainerBody:    "(class_body) @body | (object) @body",
				ContainerMethods: "(method_definition) @method",
			},
		},
		"python": {
			Extensions: []string{".py"},
			Queries: models.Queries{
				TopLevelSymbols: []string{
					"(function_definition) @symbol",
					"(class_definition) @symbol",
					"(assignment) @symbol",
				},
				ContainerBody:    "(block) @body",
				ContainerMethods: "(function_definition) @method",
			},
		},
	}

	// 创建配置文件目录
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 写入默认配置文件
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化默认配置失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return nil, fmt.Errorf("写入配置文件失败: %w", err)
	}

	fmt.Printf("已创建默认配置文件: %s\n", configPath)
	return defaultConfig, nil
}

// GetLanguageByExtension 根据文件扩展名获取语言配置
func GetLanguageByExtension(config models.LanguagesConfig, ext string) (string, models.LanguageConfig, bool) {
	for langName, langConfig := range config {
		for _, supportedExt := range langConfig.Extensions {
			if supportedExt == ext {
				return langName, langConfig, true
			}
		}
	}
	return "", models.LanguageConfig{}, false
}
