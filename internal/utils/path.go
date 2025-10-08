package utils

import (
	"path/filepath"
	"strings"
)

// NormalizePath 标准化路径，统一处理各种斜杠输入
// 这个函数确保所有路径都使用正斜杠格式，并处理各种边界情况
func NormalizePath(path string) string {
	if path == "" {
		return ""
	}

	// 统一使用正斜杠
	path = strings.ReplaceAll(path, "\\", "/")

	// 处理多个连续斜杠
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}

	// 移除末尾的斜杠（除非是根目录）
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	return path
}

// GetRelativePath 获取相对于项目根目录的相对路径
// 这个函数确保返回的路径使用正斜杠格式
func GetRelativePath(projectPath, filePath string) string {
	// 先将输入路径转换为系统原生格式，以便 filepath.Rel 正确工作
	projectPath = filepath.Clean(projectPath)
	filePath = filepath.Clean(filePath)

	relPath, err := filepath.Rel(projectPath, filePath)
	if err != nil {
		// 如果无法获取相对路径，返回标准化后的原始路径
		return NormalizePath(filePath)
	}

	// 标准化相对路径为正斜杠格式
	return NormalizePath(relPath)
}

// ResolveTargetPath 解析目标路径，支持相对路径和绝对路径
// 返回绝对路径，用于文件系统操作
func ResolveTargetPath(projectPath, targetPath string) string {
	// 如果是绝对路径，直接返回
	if filepath.IsAbs(targetPath) {
		return targetPath
	}

	// 如果是相对路径，相对于项目路径
	// 使用 filepath.Join 但然后标准化结果
	absPath := filepath.Join(projectPath, targetPath)
	return absPath
}
