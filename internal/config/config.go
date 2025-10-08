package config

import (
	"github.com/cnwinds/code-outline/internal/models"
)

// Config 表示应用程序配置
type Config struct {
	Languages   models.LanguagesConfig
	Output      string
	Exclude     []string
	ProjectPath string
}

// GetDefaultLanguagesConfig 获取默认的语言配置
func GetDefaultLanguagesConfig() models.LanguagesConfig {
	return models.LanguagesConfig{
		"go": {
			Extensions: []string{".go"},
		},
		"java": {
			Extensions: []string{".java"},
		},
		"csharp": {
			Extensions: []string{".cs"},
		},
		"cpp": {
			Extensions: []string{".cpp", ".hpp", ".cc", ".cxx"},
		},
		"c": {
			Extensions: []string{".c", ".h"},
		},
		"rust": {
			Extensions: []string{".rs"},
		},
		"javascript": {
			Extensions: []string{".js", ".jsx"},
		},
		"typescript": {
			Extensions: []string{".ts", ".tsx"},
		},
		"python": {
			Extensions: []string{".py"},
		},
	}
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
