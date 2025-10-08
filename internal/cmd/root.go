package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cnwinds/code-outline/internal/config"
	"github.com/cnwinds/code-outline/internal/models"
	"github.com/cnwinds/code-outline/internal/parser"
	"github.com/cnwinds/code-outline/internal/scanner"
	"github.com/cnwinds/code-outline/internal/updater"
	"github.com/cnwinds/code-outline/internal/utils"
)

var (
	projectPath string
	outputPath  string
	excludeDirs string
	updateFiles string
	updateDirs  string
	dataFiles   string
	dataDirs    string
	compact     bool
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "code-outline",
	Short: "code-outline - 通用型项目上下文生成器",
	Long: `code-outline 是一个高性能、跨平台的命令行工具，
用于通过静态分析为任何复杂的代码仓库生成统一、简洁且信息丰富的 code-outline.json 文件。

此文件将作为大语言模型（LLM）的"全局上下文记忆"，使其能够以前所未有的
准确性和深度来理解项目架构，从而提升代码生成、需求变更、重构和调试等任务的表现。`,
}

// generateCmd 生成命令
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成项目上下文文件",
	Long:  `扫描指定项目目录，解析代码文件，并生成 code-outline.json 文件。`,
	RunE:  runGenerate,
}

// updateCmd 更新命令
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "增量更新项目上下文文件",
	Long:  `检测文件变更并增量更新现有的 code-outline.json 文件，只重新解析已修改的文件。`,
	RunE:  runUpdate,
}

// queryCmd 查询命令
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "查询文件和方法的定义数据",
	Long:  `查询指定文件或目录中的所有文件和方法定义，返回JSON格式的数据。`,
	RunE:  runQuery,
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(queryCmd)

	// 添加generate命令行参数
	generateCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "项目路径")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "code-outline.json", "输出文件路径")
	generateCmd.Flags().StringVarP(&excludeDirs, "exclude", "e", "", "要排除的目录或文件模式，用逗号分隔")
	generateCmd.Flags().BoolVarP(&compact, "compact", "c", false, "生成紧凑的JSON输出（去掉所有空格）")

	// 添加update命令行参数
	updateCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "项目路径")
	updateCmd.Flags().StringVarP(&outputPath, "output", "o", "code-outline.json", "输出文件路径")
	updateCmd.Flags().StringVarP(&excludeDirs, "exclude", "e", "", "要排除的目录或文件模式，用逗号分隔")
	updateCmd.Flags().StringVarP(&updateFiles, "files", "f", "", "指定要更新的文件，用逗号分隔（如：file1.go,file2.js）")
	updateCmd.Flags().StringVarP(&updateDirs, "dirs", "d", "", "指定要更新的目录，用逗号分隔（如：src/,internal/）")
	updateCmd.Flags().BoolVarP(&compact, "compact", "c", false, "生成紧凑的JSON输出（去掉所有空格）")

	// 添加query命令行参数
	queryCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "项目路径")
	queryCmd.Flags().StringVarP(&outputPath, "output", "o", "", "输出文件路径（如果不指定则输出到标准输出）")
	queryCmd.Flags().BoolVarP(&compact, "compact", "c", false, "生成紧凑的JSON输出（去掉所有空格）")
	queryCmd.Flags().StringVarP(&excludeDirs, "exclude", "e", "", "要排除的目录或文件模式，用逗号分隔")
	queryCmd.Flags().StringVarP(&dataFiles, "files", "f", "", "指定要查询的文件，用逗号分隔（如：file1.go,file2.js）")
	queryCmd.Flags().StringVarP(&dataDirs, "dirs", "d", "", "指定要查询的目录，用逗号分隔（如：src/,internal/）")
}

// Execute 执行根命令
func Execute(version string) error {
	// 添加版本命令
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("code-outline %s\n", version)
			fmt.Printf("Go版本: %s\n", runtime.Version())
			fmt.Printf("操作系统: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}

// runGenerate 执行生成命令
func runGenerate(cmd *cobra.Command, args []string) error {
	// 设置默认输出路径（如果未指定）
	if outputPath == "" {
		outputPath = "code-outline.json"
	}

	fmt.Println("🚀 开始生成项目上下文...")

	// 1. 加载语言配置
	fmt.Println("📋 加载语言配置...")
	languagesConfig := config.GetDefaultLanguagesConfig()
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
		if cwd, getCwdErr := os.Getwd(); getCwdErr == nil {
			projectName = filepath.Base(cwd)
		} else {
			projectName = "Unknown Project"
		}
	}

	// 将文件路径转换为相对路径
	relativeFiles := make(map[string]models.FileInfo)
	for filePath, fileInfo := range files {
		relPath, err := filepath.Rel(projectPath, filePath)
		if err != nil {
			// 如果无法计算相对路径，使用原始路径
			relPath = filePath
		}
		relativeFiles[relPath] = fileInfo
	}

	// 获取项目根目录的绝对路径
	absProjectPath, err := filepath.Abs(projectPath)
	if err != nil {
		absProjectPath = projectPath // 如果获取绝对路径失败，使用原始路径
	}

	context := models.ProjectContext{
		ProjectName:   projectName,
		ProjectRoot:   absProjectPath,
		ProjectGoal:   "TODO: 请在此描述项目目标和主要功能",
		TechStack:     techStack,
		LastUpdated:   time.Now(),
		ModuleSummary: generateModuleSummary(relativeFiles),
		Files:         relativeFiles,
	}

	// 6. 生成JSON文件
	// 如果输出路径是相对路径，则相对于项目路径
	resolvedOutputPath := resolveOutputPath(outputPath, projectPath)
	fmt.Printf("💾 生成输出文件: %s\n", resolvedOutputPath)
	err = saveProjectContext(&context, resolvedOutputPath)
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
		dir := utils.NormalizePath(filepath.Dir(filePath))
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
	data, err := json.Marshal(context)
	if err != nil {
		return err
	}

	// 根据compact选项决定格式化方式
	if compact {
		// 紧凑模式：直接使用原始JSON（无缩进）
		// 不需要额外格式化
	} else {
		// 标准模式：格式化JSON，但保持数组在一行
		data, err = formatJSONCompact(data)
		if err != nil {
			return err
		}
	}

	// 写入文件
	return os.WriteFile(outputPath, data, 0600)
}

// formatJSONCompact 格式化JSON，保持range数组在一行，过滤空的purpose字段
func formatJSONCompact(data []byte) ([]byte, error) {
	// 解析JSON数据
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	// 过滤空的purpose字段
	jsonData = filterEmptyPurposeFields(jsonData)

	// 使用MarshalIndent格式化JSON
	formatted, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, err
	}

	// 将range数组格式化为单行
	rangePattern := regexp.MustCompile(`"range": \[\s*\n\s*(\d+),\s*\n\s*(\d+)\s*\n\s*\]`)
	formatted = rangePattern.ReplaceAll(formatted, []byte(`"range": [$1, $2]`))

	return formatted, nil
}

// filterEmptyPurposeFields 递归过滤空的purpose字段
func filterEmptyPurposeFields(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			// 跳过空的purpose字段
			if key == "purpose" {
				if str, ok := value.(string); ok && str == "" {
					continue
				}
			}
			// 递归处理嵌套结构
			result[key] = filterEmptyPurposeFields(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, 0, len(v))
		for _, item := range v {
			result = append(result, filterEmptyPurposeFields(item))
		}
		return result
	default:
		return v
	}
}

// runUpdate 执行更新命令
func runUpdate(cmd *cobra.Command, args []string) error {
	// 设置默认输出路径（如果未指定）
	if outputPath == "" {
		outputPath = "code-outline.json"
	}

	fmt.Println("🔄 开始增量更新项目上下文...")

	// 1. 加载语言配置
	languagesConfig := config.GetDefaultLanguagesConfig()

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

	// 5. 解析更新文件和目录
	var targetFiles []string
	var targetDirs []string

	if updateFiles != "" {
		rawFiles := strings.Split(updateFiles, ",")
		for _, file := range rawFiles {
			file = strings.TrimSpace(file)
			if file != "" {
				targetFiles = append(targetFiles, file)
			}
		}
	}

	if updateDirs != "" {
		rawDirs := strings.Split(updateDirs, ",")
		for _, dir := range rawDirs {
			dir = strings.TrimSpace(dir)
			if dir != "" {
				targetDirs = append(targetDirs, dir)
			}
		}
	}

	// 6. 执行增量更新
	updatedContext, changes, err := incrementalUpdater.UpdateProject(outputPath, projectPath, excludePatterns, targetFiles, targetDirs)
	if err != nil {
		return fmt.Errorf("增量更新失败: %w", err)
	}

	// 6. 如果有变更，保存更新后的上下文
	if len(changes) > 0 {
		fmt.Printf("\n📝 应用了 %d 个文件变更\n", len(changes))

		// 解析输出路径
		resolvedOutputPath := resolveOutputPath(outputPath, projectPath)
		if err := saveProjectContext(updatedContext, resolvedOutputPath); err != nil {
			return fmt.Errorf("保存更新后的上下文失败: %w", err)
		}

		fmt.Printf("💾 更新文件: %s\n", resolvedOutputPath)
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
	fmt.Printf("  模块数量: %d\n", len(context.ModuleSummary))

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

// runQuery 执行查询命令
func runQuery(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 开始查询文件和方法的定义数据...")

	// 1. 检查是否存在 code-outline.json 文件
	contextFile := filepath.Join(projectPath, "code-outline.json")
	if _, err := os.Stat(contextFile); os.IsNotExist(err) {
		return fmt.Errorf("未找到 code-outline.json 文件，请先运行 generate 命令生成项目上下文")
	}

	// 2. 加载项目上下文文件
	fmt.Println("📂 加载项目上下文文件...")
	context, err := loadProjectContext(contextFile)
	if err != nil {
		return fmt.Errorf("加载项目上下文失败: %w", err)
	}
	fmt.Printf("✅ 已加载项目上下文: %s\n", context.ProjectName)

	// 3. 解析目标文件和目录
	var targetFiles []string
	var targetDirs []string

	if dataFiles != "" {
		rawFiles := strings.Split(dataFiles, ",")
		for _, file := range rawFiles {
			file = strings.TrimSpace(file)
			if file != "" {
				targetFiles = append(targetFiles, file)
			}
		}
	}

	if dataDirs != "" {
		rawDirs := strings.Split(dataDirs, ",")
		for _, dir := range rawDirs {
			dir = strings.TrimSpace(dir)
			if dir != "" {
				targetDirs = append(targetDirs, dir)
			}
		}
	}

	// 4. 从上下文中提取数据
	fmt.Println("🔍 从项目上下文中提取数据...")
	dataResult, err := extractDataFromContext(context, targetFiles, targetDirs)
	if err != nil {
		return fmt.Errorf("提取数据失败: %w", err)
	}

	// 5. 输出结果
	if outputPath != "" {
		fmt.Printf("💾 保存数据到文件: %s\n", outputPath)
		err = saveDataToFile(dataResult, outputPath)
		if err != nil {
			return fmt.Errorf("保存数据失败: %w", err)
		}
		fmt.Println("✅ 数据已保存到文件")
	} else {
		// 输出到标准输出
		jsonData, err := json.Marshal(dataResult)
		if err != nil {
			return fmt.Errorf("序列化数据失败: %w", err)
		}

		// 根据compact选项决定格式化方式
		if compact {
			// 紧凑模式：直接输出原始JSON
			fmt.Println(string(jsonData))
		} else {
			// 标准模式：格式化JSON，但保持range数组在一行
			jsonData, err = formatJSONCompact(jsonData)
			if err != nil {
				return fmt.Errorf("格式化JSON失败: %w", err)
			}
			fmt.Println(string(jsonData))
		}
	}

	// 6. 显示统计信息
	printDataStatistics(dataResult)

	fmt.Println("🎉 查询完成!")
	return nil
}

// DataResult 数据获取结果
type DataResult struct {
	Files map[string]models.FileInfo `json:"files"`
	Stats DataStats                  `json:"stats"`
}

// DataStats 数据统计信息
type DataStats struct {
	TotalFiles   int      `json:"totalFiles"`
	TotalSymbols int      `json:"totalSymbols"`
	Languages    []string `json:"languages"`
}

// saveDataToFile 保存数据到文件
func saveDataToFile(data *DataResult, outputPath string) error {
	// 创建输出目录（如果不存在）
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// 序列化为JSON
	var dataBytes []byte
	var err error

	if compact {
		// 紧凑模式：无缩进
		dataBytes, err = json.Marshal(data)
	} else {
		// 标准模式：有缩进
		dataBytes, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(outputPath, dataBytes, 0600)
}

// printDataStatistics 打印数据统计信息
func printDataStatistics(data *DataResult) {
	fmt.Printf("\n📊 数据统计:\n")
	fmt.Printf("  📁 文件数量: %d\n", data.Stats.TotalFiles)
	fmt.Printf("  🔍 符号数量: %d\n", data.Stats.TotalSymbols)
	fmt.Printf("  🛠️  技术栈: %s\n", strings.Join(data.Stats.Languages, ", "))
}

// loadProjectContext 加载项目上下文文件
func loadProjectContext(filePath string) (*models.ProjectContext, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var context models.ProjectContext
	err = json.Unmarshal(data, &context)
	if err != nil {
		return nil, err
	}

	return &context, nil
}

// extractDataFromContext 从项目上下文中提取指定文件或目录的数据
func extractDataFromContext(context *models.ProjectContext, targetFiles, targetDirs []string) (*DataResult, error) {
	result := &DataResult{
		Files: make(map[string]models.FileInfo),
		Stats: DataStats{
			TotalFiles:   0,
			TotalSymbols: 0,
			Languages:    []string{},
		},
	}

	// 用于统计语言的map
	languageCount := make(map[string]int)

	// 如果没有指定目标，返回所有文件
	if len(targetFiles) == 0 && len(targetDirs) == 0 {
		for filePath, fileInfo := range context.Files {
			result.Files[filePath] = fileInfo
			result.Stats.TotalFiles++
			result.Stats.TotalSymbols += len(fileInfo.Symbols)

			// 统计语言
			ext := filepath.Ext(filePath)
			if ext != "" {
				ext = ext[1:] // 移除点号
				languageCount[ext]++
			}
		}
	} else {
		// 处理指定的文件
		for _, targetFile := range targetFiles {
			// 标准化文件路径
			normalizedFile := utils.NormalizePath(targetFile)

			// 查找匹配的文件
			for filePath, fileInfo := range context.Files {
				if strings.Contains(filePath, normalizedFile) || filepath.Base(filePath) == filepath.Base(normalizedFile) {
					result.Files[filePath] = fileInfo
					result.Stats.TotalFiles++
					result.Stats.TotalSymbols += len(fileInfo.Symbols)

					// 统计语言
					ext := filepath.Ext(filePath)
					if ext != "" {
						ext = ext[1:] // 移除点号
						languageCount[ext]++
					}
				}
			}
		}

		// 处理指定的目录
		for _, targetDir := range targetDirs {
			// 标准化目录路径
			normalizedDir := utils.NormalizePath(targetDir)

			// 查找匹配目录下的文件
			for filePath, fileInfo := range context.Files {
				fileDir := filepath.Dir(filePath)
				if strings.HasPrefix(fileDir, normalizedDir) || strings.Contains(filePath, normalizedDir) {
					// 避免重复添加
					if _, exists := result.Files[filePath]; !exists {
						result.Files[filePath] = fileInfo
						result.Stats.TotalFiles++
						result.Stats.TotalSymbols += len(fileInfo.Symbols)

						// 统计语言
						ext := filepath.Ext(filePath)
						if ext != "" {
							ext = ext[1:] // 移除点号
							languageCount[ext]++
						}
					}
				}
			}
		}
	}

	// 将语言统计转换为切片
	for lang := range languageCount {
		result.Stats.Languages = append(result.Stats.Languages, lang)
	}

	return result, nil
}

// resolveOutputPath 解析输出路径，如果输出路径是相对路径，则相对于项目路径
func resolveOutputPath(outputPath, projectPath string) string {
	// 如果是绝对路径，直接返回
	if filepath.IsAbs(outputPath) {
		return outputPath
	}

	// 如果是相对路径，相对于项目路径
	return filepath.Join(projectPath, outputPath)
}
