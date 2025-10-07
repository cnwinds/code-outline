package main

import (
	"fmt"
	"time"
)

// Greeter 是一个简单的问候器
type Greeter struct {
	Name string
}

// Config 配置结构体
type Config struct {
	DatabaseURL string
	APIVersion  string
	Environment string
}

// Options 选项结构体
type Options struct {
	EnableLogging  bool
	MaxConnections int
	CacheSize      int
}

// Handler 处理器接口
type Handler interface {
	Process(data interface{}) error
}

// ComplexStruct 复杂结构体
type ComplexStruct struct {
	Config    Config
	Options   Options
	Handlers  []Handler
	Timeout   int
	Retries   int
	Debug     bool
	CreatedAt time.Time
}

// NewGreeter 创建一个新的问候器实例
func NewGreeter(name string) *Greeter {
	return &Greeter{Name: name}
}

// SayHello 向指定的人打招呼
func (g *Greeter) SayHello(person string) string {
	return fmt.Sprintf("Hello %s, I'm %s!", person, g.Name)
}

// ProcessUserData 处理用户数据，包含多个参数
func ProcessUserData(
	userID int,
	userName string,
	email string,
	age int,
	isActive bool,
	preferences map[string]interface{},
	callback func(string) error,
) (string, error) {
	// 验证用户数据
	if userID <= 0 {
		return "", fmt.Errorf("invalid user ID: %d", userID)
	}

	if userName == "" {
		return "", fmt.Errorf("user name cannot be empty")
	}

	// 处理用户数据
	result := fmt.Sprintf("Processing user: %s (ID: %d, Email: %s, Age: %d, Active: %t)",
		userName, userID, email, age, isActive)

	// 执行回调函数
	if callback != nil {
		if err := callback(result); err != nil {
			return "", fmt.Errorf("callback error: %w", err)
		}
	}

	return result, nil
}

// CreateComplexStruct 创建一个复杂的结构体，包含多行参数
func CreateComplexStruct(
	config Config,
	options Options,
	handlers []Handler,
	timeout int,
	retries int,
	debug bool,
) (*ComplexStruct, error) {
	// 验证参数
	if timeout < 0 {
		return nil, fmt.Errorf("timeout must be non-negative")
	}

	if retries < 0 {
		return nil, fmt.Errorf("retries must be non-negative")
	}

	// 创建复杂结构体
	cs := &ComplexStruct{
		Config:    config,
		Options:   options,
		Handlers:  handlers,
		Timeout:   timeout,
		Retries:   retries,
		Debug:     debug,
		CreatedAt: time.Now(),
	}

	return cs, nil
}

// main 程序入口点
func main() {
	greeter := NewGreeter("CodeCartographer")
	message := greeter.SayHello("World")
	fmt.Println(message)

	// 演示多行参数函数的使用
	preferences := map[string]interface{}{
		"theme": "dark",
		"lang":  "zh-CN",
	}

	// 调用ProcessUserData函数
	result, err := ProcessUserData(
		123,
		"张三",
		"zhangsan@example.com",
		25,
		true,
		preferences,
		func(msg string) error {
			fmt.Println("回调函数执行:", msg)
			return nil
		},
	)

	if err != nil {
		fmt.Printf("处理用户数据时出错: %v\n", err)
	} else {
		fmt.Println("处理结果:", result)
	}

	// 演示CreateComplexStruct函数
	config := Config{
		DatabaseURL: "postgres://localhost:5432/mydb",
		APIVersion:  "v1",
		Environment: "development",
	}

	options := Options{
		EnableLogging:  true,
		MaxConnections: 100,
		CacheSize:      1024,
	}

	cs, err := CreateComplexStruct(
		config,
		options,
		nil,  // handlers
		30,   // timeout
		3,    // retries
		true, // debug
	)

	if err != nil {
		fmt.Printf("创建复杂结构体时出错: %v\n", err)
	} else {
		fmt.Printf("创建成功: %+v\n", cs)
	}
}
