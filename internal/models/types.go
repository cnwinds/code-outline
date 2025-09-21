package models

import "time"

// Symbol 表示代码中的一个符号（如函数、结构体、常量等）
type Symbol struct {
	Prototype string   `json:"prototype"`           // 符号的完整声明行
	Purpose   string   `json:"purpose"`             // 从注释中提取的说明
	Range     []int    `json:"range"`               // [start_line, end_line]
	Body      string   `json:"body,omitempty"`      // 用于类/结构体/接口等容器类型的内部内容
	Methods   []Symbol `json:"methods,omitempty"`   // 用于类/结构体的方法
}

// FileInfo 表示一个文件的信息
type FileInfo struct {
	Purpose string   `json:"purpose"` // 文件的用途描述
	Symbols []Symbol `json:"symbols"` // 文件中的符号列表
}

// Architecture 表示项目架构信息
type Architecture struct {
	Overview      string            `json:"overview"`      // 架构概述
	ModuleSummary map[string]string `json:"moduleSummary"` // 模块摘要
}

// ProjectContext 表示整个项目的上下文信息
type ProjectContext struct {
	ProjectName  string                `json:"projectName"`  // 项目名称
	ProjectGoal  string                `json:"projectGoal"`  // 项目目标
	TechStack    []string              `json:"techStack"`    // 技术栈
	LastUpdated  time.Time             `json:"lastUpdated"`  // 最后更新时间
	Architecture Architecture          `json:"architecture"` // 架构信息
	Files        map[string]FileInfo   `json:"files"`        // 文件信息映射
}

// LanguageConfig 表示单个语言的配置
type LanguageConfig struct {
	Extensions  []string `json:"extensions"`  // 文件扩展名列表
	GrammarPath string   `json:"grammar_path"` // Tree-sitter语法库路径
	Queries     Queries  `json:"queries"`     // 查询规则
}

// Queries 表示Tree-sitter查询规则
type Queries struct {
	TopLevelSymbols  []string `json:"top_level_symbols"`  // 顶级符号查询
	ContainerBody    string   `json:"container_body"`     // 容器主体查询
	ContainerMethods string   `json:"container_methods"`  // 容器方法查询
}

// LanguagesConfig 表示所有语言的配置
type LanguagesConfig map[string]LanguageConfig
