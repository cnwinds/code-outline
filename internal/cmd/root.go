package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cnwinds/CodeCartographer/internal/config"
	"github.com/cnwinds/CodeCartographer/internal/models"
	"github.com/cnwinds/CodeCartographer/internal/parser"
	"github.com/cnwinds/CodeCartographer/internal/scanner"
	"github.com/cnwinds/CodeCartographer/internal/updater"
	"github.com/spf13/cobra"
)

var (
	projectPath string
	outputPath  string
	configPath  string
	excludeDirs string
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "contextgen",
	Short: "CodeCartographer - 通用型项目上下文生成器",
	Long: `CodeCartographer 是一个高性能、跨平台的命令行工具，
用于通过静态分析为任何复杂的代码仓库生成统一、简洁且信息丰富的 project_context.json 文件。

此文件将作为大语言模型（LLM）的"全局上下文记忆"，使其能够以前所未有的
准确性和深度来理解项目架构，从而提升代码生成、需求变更、重构和调试等任务的表现。`,
}

// generateCmd 生成命令
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成项目上下文文件",
	Long:  `扫描指定项目目录，解析代码文件，并生成 project_context.json 文件。`,
	RunE:  runGenerate,
}

// updateCmd 更新命令
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "增量更新项目上下文文件",
	Long:  `检测文件变更并增量更新现有的 project_context.json 文件，只重新解析已修改的文件。`,
	RunE:  runUpdate,
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(updateCmd)

	// 添加generate命令行参数
	generateCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "项目路径")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "project_context.json", "输出文件路径")
	generateCmd.Flags().StringVarP(&configPath, "config", "c", "", "语言配置文件路径")
	generateCmd.Flags().StringVarP(&excludeDirs, "exclude", "e", "", "要排除的目录或文件模式，用逗号分隔")

	// 添加update命令行参数
	updateCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "项目路径")
	updateCmd.Flags().StringVarP(&outputPath, "output", "o", "project_context.json", "输出文件路径")
	updateCmd.Flags().StringVarP(&configPath, "config", "c", "", "语言配置文件路径")
	updateCmd.Flags().StringVarP(&excludeDirs, "exclude", "e", "", "要排除的目录或文件模式，用逗号分隔")
}

// Execute 执行根命令
func Execute(version string) error {
	// 添加版本命令
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("CodeCartographer %s\n", version)
			fmt.Printf("Go版本: %s\n", runtime.Version())
			fmt.Printf("操作系统: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}

// runGenerate 执行生成命令
func runGenerate(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 开始生成项目上下文...")

	// 1. 加载语言配置
	fmt.Println("📋 加载语言配置...")
	languagesConfig, err := config.LoadLanguagesConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载语言配置失败: %w", err)
	}
	fmt.Printf("✅ 已加载 %d 种语言的配置\n", len(languagesConfig))

	// 2. 创建解析器
	fmt.Println("🔧 初始化解析器...")
	fmt.Println("🌳 使用 Tree-sitter 解析器")
	treeSitterParser, err := parser.NewTreeSitterParser(languagesConfig)
	if err != nil {
		return fmt.Errorf("tree-sitter 解析器初始化失败: %w", err)
	}
	codeParser := treeSitterParser

	// 3. 解析排除模式
	var excludePatterns []string
	if excludeDirs != "" {
		excludePatterns = strings.Split(excludeDirs, ",")
		for i, pattern := range excludePatterns {
			excludePatterns[i] = strings.TrimSpace(pattern)
		}
	}

	// 4. 创建扫描器并扫描项目
	fmt.Printf("🔍 扫描项目: %s\n", projectPath)
	fileScanner := scanner.NewScanner(codeParser, excludePatterns)
	files, techStack, err := fileScanner.ScanProject(projectPath)
	if err != nil {
		return fmt.Errorf("扫描项目失败: %w", err)
	}
	fmt.Printf("✅ 扫描完成，找到 %d 个文件\n", len(files))

	// 5. 构建项目上下文
	fmt.Println("📦 构建项目上下文...")
	projectName := filepath.Base(projectPath)
	if projectName == "." {
		if cwd, err := os.Getwd(); err == nil {
			projectName = filepath.Base(cwd)
		} else {
			projectName = "Unknown Project"
		}
	}

	context := models.ProjectContext{
		ProjectName: projectName,
		ProjectGoal: "TODO: 请在此描述项目目标和主要功能",
		TechStack:   techStack,
		LastUpdated: time.Now(),
		Architecture: models.Architecture{
			Overview:      "TODO: 请在此描述项目的整体架构",
			ModuleSummary: generateModuleSummary(files),
		},
		Files: files,
	}

	// 6. 生成JSON文件
	fmt.Printf("💾 生成输出文件: %s\n", outputPath)
	err = saveProjectContext(&context, outputPath)
	if err != nil {
		return fmt.Errorf("保存项目上下文失败: %w", err)
	}

	// 7. 显示统计信息
	printStatistics(&context)

	fmt.Println("🎉 项目上下文生成完成!")
	return nil
}

// generateModuleSummary 生成模块摘要
func generateModuleSummary(files map[string]models.FileInfo) map[string]string {
	moduleSummary := make(map[string]string)

	// 按目录分组文件
	dirGroups := make(map[string][]string)
	for filePath := range files {
		dir := filepath.Dir(filePath)
		if dir == "." {
			dir = "root"
		}
		dirGroups[dir] = append(dirGroups[dir], filePath)
	}

	// 为每个目录生成摘要
	for dir, fileList := range dirGroups {
		if len(fileList) == 1 {
			moduleSummary[dir] = fmt.Sprintf("包含 1 个文件: %s", filepath.Base(fileList[0]))
		} else {
			moduleSummary[dir] = fmt.Sprintf("包含 %d 个文件，主要用于 TODO: 请描述此模块的用途", len(fileList))
		}
	}

	return moduleSummary
}

// saveProjectContext 保存项目上下文到JSON文件
func saveProjectContext(context *models.ProjectContext, outputPath string) error {
	// 创建输出目录（如果不存在）
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(outputPath, data, 0600)
}

// runUpdate 执行更新命令
func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("🔄 开始增量更新项目上下文...")

	// 1. 加载语言配置
	languagesConfig, err := config.LoadLanguagesConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载语言配置失败: %w", err)
	}

	// 2. 创建解析器
	fmt.Println("🌳 使用 Tree-sitter 解析器")
	treeSitterParser, err := parser.NewTreeSitterParser(languagesConfig)
	if err != nil {
		return fmt.Errorf("tree-sitter 解析器初始化失败: %w", err)
	}
	fileParser := treeSitterParser

	// 3. 创建增量更新器
	incrementalUpdater := updater.NewIncrementalUpdater(fileParser)

	// 4. 解析排除模式
	var excludePatterns []string
	if excludeDirs != "" {
		excludePatterns = strings.Split(excludeDirs, ",")
		for i, pattern := range excludePatterns {
			excludePatterns[i] = strings.TrimSpace(pattern)
		}
	}

	// 5. 执行增量更新
	updatedContext, changes, err := incrementalUpdater.UpdateProject(outputPath, projectPath, excludePatterns)
	if err != nil {
		return fmt.Errorf("增量更新失败: %w", err)
	}

	// 6. 如果有变更，保存更新后的上下文
	if len(changes) > 0 {
		fmt.Printf("\n📝 应用了 %d 个文件变更\n", len(changes))

		if err := saveProjectContext(updatedContext, outputPath); err != nil {
			return fmt.Errorf("保存更新后的上下文失败: %w", err)
		}

		fmt.Printf("💾 更新文件: %s\n", outputPath)
	}

	// 7. 打印统计信息
	printUpdateStatistics(updatedContext, changes)

	return nil
}

// printStatistics 打印统计信息
func printStatistics(context *models.ProjectContext) {
	fmt.Println("\n📊 统计信息:")
	fmt.Printf("  项目名称: %s\n", context.ProjectName)
	fmt.Printf("  技术栈: %s\n", strings.Join(context.TechStack, ", "))
	fmt.Printf("  文件数量: %d\n", len(context.Files))
	fmt.Printf("  模块数量: %d\n", len(context.Architecture.ModuleSummary))

	// 统计符号数量
	totalSymbols := 0
	for _, fileInfo := range context.Files {
		totalSymbols += len(fileInfo.Symbols)
	}
	fmt.Printf("  符号数量: %d\n", totalSymbols)

	fmt.Printf("  最后更新: %s\n", context.LastUpdated.Format("2006-01-02 15:04:05"))
}

// printUpdateStatistics 打印更新统计信息
func printUpdateStatistics(context *models.ProjectContext, changes []updater.FileChange) {
	fmt.Printf("\n📊 更新统计:\n")

	addedCount := 0
	modifiedCount := 0
	deletedCount := 0

	for _, change := range changes {
		switch change.ChangeType {
		case updater.FileAdded:
			addedCount++
		case updater.FileModified:
			modifiedCount++
		case updater.FileDeleted:
			deletedCount++
		}
	}

	if len(changes) > 0 {
		fmt.Printf("  📁 文件变更: +%d ✏️%d 🗑️%d\n", addedCount, modifiedCount, deletedCount)
	}

	fmt.Printf("  📄 总文件数量: %d\n", len(context.Files))

	// 统计符号数量
	symbolCount := 0
	for _, fileInfo := range context.Files {
		symbolCount += len(fileInfo.Symbols)
	}
	fmt.Printf("  🔍 总符号数量: %d\n", symbolCount)
	fmt.Printf("  ⏰ 最后更新: %s\n", context.LastUpdated.Format("2006-01-02 15:04:05"))
}
