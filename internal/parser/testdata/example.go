package main

import "fmt"

// Greeter 是一个简单的问候器
type Greeter struct {
	Name string
}

// NewGreeter 创建一个新的问候器实例
func NewGreeter(name string) *Greeter {
	return &Greeter{Name: name}
}

// SayHello 向指定的人打招呼
func (g *Greeter) SayHello(person string) string {
	return fmt.Sprintf("Hello %s, I'm %s!", person, g.Name)
}

// main 程序入口点
func main() {
	greeter := NewGreeter("CodeCartographer")
	message := greeter.SayHello("World")
	fmt.Println(message)
}
