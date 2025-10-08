package models

import "time"

// Symbol 表示代码中的一个符号（如函数、结构体、常量等）
type Symbol struct {
	Prototype string   `json:"prototype"`         // 符号的完整声明行
	Purpose   string   `json:"purpose"`           // 从注释中提取的说明
	Range     []int    `json:"range"`             // [start_line, end_line]
	Body      string   `json:"body,omitempty"`    // 用于类/结构体/接口等容器类型的内部内容
	Methods   []Symbol `json:"methods,omitempty"` // 用于类/结构体的方法
}

// FileInfo 表示一个文件的信息
type FileInfo struct {
	Purpose      string   `json:"purpose"`      // 文件的用途描述
	Symbols      []Symbol `json:"symbols"`      // 文件中的符号列表
	LastModified string   `json:"lastModified"` // 文件最后修改时间
	FileSize     int64    `json:"fileSize"`     // 文件大小
}

// ProjectContext 表示整个项目的上下文信息
type ProjectContext struct {
	ProjectName   string              `json:"projectName"`   // 项目名称
	ProjectRoot   string              `json:"projectRoot"`   // 项目根目录
	ProjectGoal   string              `json:"projectGoal"`   // 项目目标
	TechStack     []string            `json:"techStack"`     // 技术栈
	LastUpdated   time.Time           `json:"lastUpdated"`   // 最后更新时间
	ModuleSummary map[string]string   `json:"moduleSummary"` // 模块摘要
	Files         map[string]FileInfo `json:"files"`         // 文件信息映射（相对路径）
}

// LanguageConfig 表示单个语言的配置
type LanguageConfig struct {
	Extensions []string `json:"extensions"` // 文件扩展名列表
}

// LanguagesConfig 表示所有语言的配置
type LanguagesConfig map[string]LanguageConfig
