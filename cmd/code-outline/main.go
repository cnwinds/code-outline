package main

import (
	"fmt"
	"os"

	"github.com/cnwinds/code-outline/internal/cmd"
)

// Version 版本号
var Version = "v0.1.0"

// main 程序入口点，启动 CodeCartographer 命令行工具
func main() {
	if err := cmd.Execute(Version); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
