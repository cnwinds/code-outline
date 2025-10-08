package updater

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cnwinds/code-outline/internal/models"
	"github.com/cnwinds/code-outline/internal/scanner"
	"github.com/cnwinds/code-outline/internal/utils"
)

// IncrementalUpdater 增量更新器
type IncrementalUpdater struct {
	parser scanner.FileParser
}

// NewIncrementalUpdater 创建新的增量更新器
func NewIncrementalUpdater(p scanner.FileParser) *IncrementalUpdater {
	return &IncrementalUpdater{
		parser: p,
	}
}

// FileChangeType 文件变更类型
type FileChangeType int

const (
	FileAdded FileChangeType = iota
	FileModified
	FileDeleted
)

// FileChange 文件变更信息
type FileChange struct {
	Path       string
	ChangeType FileChangeType
	OldInfo    *models.FileInfo
	NewInfo    *models.FileInfo
}

// UpdateProject 增量更新项目上下文
func (u *IncrementalUpdater) UpdateProject(
	contextPath, projectPath string,
	excludePatterns []string,
	targetFiles []string,
	targetDirs []string,
) (*models.ProjectContext, []FileChange, error) {
	// 1. 加载现有的项目上下文
	existingContext, err := u.loadExistingContext(contextPath)
	if err != nil {
		return nil, nil, fmt.Errorf("加载现有上下文失败: %w", err)
	}

	// 2. 扫描项目文件，检测变更
	changes, err := u.detectFileChanges(existingContext, projectPath, excludePatterns, targetFiles, targetDirs)
	if err != nil {
		return nil, nil, fmt.Errorf("检测文件变更失败: %w", err)
	}

	// 3. 如果没有变更，直接返回
	if len(changes) == 0 {
		fmt.Println("✅ 没有检测到文件变更")
		return existingContext, changes, nil
	}

	// 4. 应用变更
	updatedContext := u.applyChanges(existingContext, changes)

	// 5. 更新时间戳
	updatedContext.LastUpdated = time.Now()

	return updatedContext, changes, nil
}

// loadExistingContext 加载现有的项目上下文
func (u *IncrementalUpdater) loadExistingContext(contextPath string) (*models.ProjectContext, error) {
	if _, err := os.Stat(contextPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("上下文文件不存在: %s", contextPath)
	}

	data, err := os.ReadFile(contextPath)
	if err != nil {
		return nil, fmt.Errorf("读取上下文文件失败: %w", err)
	}

	var context models.ProjectContext
	if err := json.Unmarshal(data, &context); err != nil {
		return nil, fmt.Errorf("解析上下文文件失败: %w", err)
	}

	// 标准化所有文件路径键，确保使用正斜杠格式
	normalizedFiles := make(map[string]models.FileInfo, len(context.Files))
	for filePath, fileInfo := range context.Files {
		normalizedPath := utils.NormalizePath(filePath)
		normalizedFiles[normalizedPath] = fileInfo
	}
	context.Files = normalizedFiles

	return &context, nil
}

// detectFileChanges 检测文件变更
func (u *IncrementalUpdater) detectFileChanges(
	context *models.ProjectContext,
	projectPath string,
	excludePatterns []string,
	targetFiles []string,
	targetDirs []string,
) ([]FileChange, error) {
	var changes []FileChange
	currentFiles := make(map[string]bool)

	// 如果指定了目标文件或目录，只处理这些文件
	if len(targetFiles) > 0 || len(targetDirs) > 0 {
		var err error
		changes, err = u.detectTargetChanges(context, projectPath, excludePatterns, targetFiles, targetDirs)
		if err != nil {
			return nil, err
		}
		return changes, nil
	}

	// 遍历项目文件
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查是否应该排除
		if u.shouldExclude(path, excludePatterns) {
			return nil
		}

		// 检查是否为支持的文件类型
		ext := filepath.Ext(path)
		if !u.isSupportedFile(ext) {
			return nil
		}

		// 转换为相对路径
		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath) // 统一使用斜杠

		currentFiles[relPath] = true

		// 检查文件是否存在于上下文中
		existingFile, exists := context.Files[relPath]
		if !exists {
			// 新文件
			newInfo, err := u.parser.ParseFile(path)
			if err != nil {
				return err
			}
			changes = append(changes, FileChange{
				Path:       relPath,
				ChangeType: FileAdded,
				NewInfo:    newInfo,
			})
		} else if u.isFileModified(path, &existingFile) {
			// 检查文件是否被修改
			newInfo, err := u.parser.ParseFile(path)
			if err != nil {
				return err
			}
			changes = append(changes, FileChange{
				Path:       relPath,
				ChangeType: FileModified,
				OldInfo:    &existingFile,
				NewInfo:    newInfo,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 检查已删除的文件
	for filePath := range context.Files {
		if !currentFiles[filePath] {
			existingFile := context.Files[filePath]
			changes = append(changes, FileChange{
				Path:       filePath,
				ChangeType: FileDeleted,
				OldInfo:    &existingFile,
			})
		}
	}

	return changes, nil
}

// detectTargetChanges 检测指定文件和目录的变更
func (u *IncrementalUpdater) detectTargetChanges(
	context *models.ProjectContext,
	projectPath string,
	excludePatterns []string,
	targetFiles []string,
	targetDirs []string,
) ([]FileChange, error) {
	var changes []FileChange

	// 处理指定的文件
	for _, targetFile := range targetFiles {
		// 标准化路径
		normalizedFile := utils.NormalizePath(targetFile)
		resolvedPath := utils.ResolveTargetPath(projectPath, normalizedFile)

		// 获取相对路径作为键
		relPath := utils.GetRelativePath(projectPath, resolvedPath)

		// 检查文件是否存在
		if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
			// 文件不存在，检查是否在上下文中（可能是被删除的文件）
			if existingFile, exists := context.Files[relPath]; exists {
				changes = append(changes, FileChange{
					Path:       relPath,
					ChangeType: FileDeleted,
					OldInfo:    &existingFile,
				})
			}
			continue
		}

		// 检查是否应该排除
		if u.shouldExclude(resolvedPath, excludePatterns) {
			continue
		}

		// 检查是否为支持的文件类型
		ext := filepath.Ext(resolvedPath)
		if !u.isSupportedFile(ext) {
			continue
		}

		// 检查文件是否存在于上下文中
		existingFile, exists := context.Files[relPath]
		if !exists {
			// 新文件
			newInfo, err := u.parser.ParseFile(resolvedPath)
			if err != nil {
				return nil, fmt.Errorf("解析文件失败 %s: %w", resolvedPath, err)
			}
			changes = append(changes, FileChange{
				Path:       relPath,
				ChangeType: FileAdded,
				NewInfo:    newInfo,
			})
		} else if u.isFileModified(resolvedPath, &existingFile) {
			// 文件被修改
			newInfo, err := u.parser.ParseFile(resolvedPath)
			if err != nil {
				return nil, fmt.Errorf("解析文件失败 %s: %w", resolvedPath, err)
			}
			changes = append(changes, FileChange{
				Path:       relPath,
				ChangeType: FileModified,
				OldInfo:    &existingFile,
				NewInfo:    newInfo,
			})
		}
	}

	// 处理指定的目录
	for _, targetDir := range targetDirs {
		// 标准化路径
		normalizedDir := utils.NormalizePath(targetDir)
		resolvedDirPath := utils.ResolveTargetPath(projectPath, normalizedDir)

		// 检查目录是否存在
		if _, err := os.Stat(resolvedDirPath); os.IsNotExist(err) {
			// 目录不存在，检查上下文中是否有该目录下的文件（可能是被删除的）
			for filePath := range context.Files {
				if strings.HasPrefix(filePath, normalizedDir) {
					existingFile := context.Files[filePath]
					changes = append(changes, FileChange{
						Path:       filePath,
						ChangeType: FileDeleted,
						OldInfo:    &existingFile,
					})
				}
			}
			continue
		}

		// 遍历目录下的文件
		err := filepath.Walk(resolvedDirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 跳过目录
			if info.IsDir() {
				return nil
			}

			// 检查是否应该排除
			if u.shouldExclude(path, excludePatterns) {
				return nil
			}

			// 检查是否为支持的文件类型
			ext := filepath.Ext(path)
			if !u.isSupportedFile(ext) {
				return nil
			}

			// 转换为相对路径
			relPath := utils.GetRelativePath(projectPath, path)

			// 检查文件是否存在于上下文中
			existingFile, exists := context.Files[relPath]
			if !exists {
				// 新文件
				newInfo, err := u.parser.ParseFile(path)
				if err != nil {
					return fmt.Errorf("解析文件失败 %s: %w", path, err)
				}
				changes = append(changes, FileChange{
					Path:       relPath,
					ChangeType: FileAdded,
					NewInfo:    newInfo,
				})
			} else if u.isFileModified(path, &existingFile) {
				// 文件被修改
				newInfo, err := u.parser.ParseFile(path)
				if err != nil {
					return fmt.Errorf("解析文件失败 %s: %w", path, err)
				}
				changes = append(changes, FileChange{
					Path:       relPath,
					ChangeType: FileModified,
					OldInfo:    &existingFile,
					NewInfo:    newInfo,
				})
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("遍历目录失败 %s: %w", resolvedDirPath, err)
		}
	}

	return changes, nil
}

// isFileModified 检查文件是否被修改
func (u *IncrementalUpdater) isFileModified(filePath string, existingInfo *models.FileInfo) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return true // 如果无法获取文件信息，假设已修改
	}

	// 比较修改时间和文件大小
	currentModTime := fileInfo.ModTime().Format(time.RFC3339)
	currentSize := fileInfo.Size()

	return currentModTime != existingInfo.LastModified || currentSize != existingInfo.FileSize
}

// applyChanges 应用文件变更
func (u *IncrementalUpdater) applyChanges(context *models.ProjectContext, changes []FileChange) *models.ProjectContext {
	// 创建上下文副本
	updatedContext := *context
	updatedFiles := make(map[string]models.FileInfo)

	// 复制现有文件
	for path, info := range context.Files {
		updatedFiles[path] = info
	}

	// 应用变更
	for _, change := range changes {
		switch change.ChangeType {
		case FileAdded:
			updatedFiles[change.Path] = *change.NewInfo
			fmt.Printf("➕ 添加文件: %s\n", change.Path)
		case FileModified:
			updatedFiles[change.Path] = *change.NewInfo
			fmt.Printf("✏️  修改文件: %s\n", change.Path)
		case FileDeleted:
			delete(updatedFiles, change.Path)
			fmt.Printf("🗑️  删除文件: %s\n", change.Path)
		}
	}

	updatedContext.Files = updatedFiles

	// 重新生成模块摘要
	updatedContext.ModuleSummary = u.generateModuleSummary(updatedFiles)

	return &updatedContext
}

// generateModuleSummary 生成模块摘要
func (u *IncrementalUpdater) generateModuleSummary(files map[string]models.FileInfo) map[string]string {
	moduleSummary := make(map[string]string)
	moduleFiles := make(map[string][]string)

	// 按模块分组文件
	for filePath := range files {
		dir := utils.NormalizePath(filepath.Dir(filePath))
		if dir == "." {
			dir = "root"
		}
		moduleFiles[dir] = append(moduleFiles[dir], filepath.Base(filePath))
	}

	// 生成摘要
	for module, fileList := range moduleFiles {
		if len(fileList) == 1 {
			moduleSummary[module] = fmt.Sprintf("包含 1 个文件: %s", fileList[0])
		} else {
			moduleSummary[module] = fmt.Sprintf("包含 %d 个文件: %s", len(fileList), strings.Join(fileList, ", "))
		}
	}

	return moduleSummary
}

// shouldExclude 检查路径是否应该被排除
func (u *IncrementalUpdater) shouldExclude(path string, excludePatterns []string) bool {
	path = filepath.ToSlash(path)
	for _, pattern := range excludePatterns {
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

// isSupportedFile 检查是否为支持的文件类型
func (u *IncrementalUpdater) isSupportedFile(ext string) bool {
	supportedExts := []string{".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".java", ".cs", ".rs", ".cpp", ".c", ".h"}
	for _, supportedExt := range supportedExts {
		if ext == supportedExt {
			return true
		}
	}
	return false
}
