/*
Golang并发编程学习Demo - 简单级别
文件：01_basic_goroutine.go
主题：基础Goroutine使用

本示例演示：
1. 如何创建和启动goroutine
2. goroutine的并发执行特性
3. 主goroutine与子goroutine的关系
4. 使用time.Sleep等待goroutine完成（不推荐的方式）

学习要点：
- go关键字启动goroutine
- 多个goroutine并发执行
- 主程序需要等待goroutine完成
- goroutine的执行顺序是不确定的

运行方式：go run simple/01_basic_goroutine.go
*/

package main

import (
	"fmt"
	"time"
)

// printNumbers 打印数字序列
// 演示一个简单的goroutine任务
func printNumbers() {
	for i := 1; i <= 5; i++ {
		fmt.Printf("数字: %d\n", i)
		// 模拟一些工作时间，让并发效果更明显
		time.Sleep(100 * time.Millisecond)
	}
}

// printLetters 打印字母序列
// 演示另一个并发执行的goroutine任务
func printLetters() {
	for i := 'A'; i <= 'E'; i++ {
		fmt.Printf("字母: %c\n", i)
		// 不同的睡眠时间，展示goroutine独立执行
		time.Sleep(150 * time.Millisecond)
	}
}

func main() {
	fmt.Println("=== 基础Goroutine演示 ===")
	fmt.Println("观察数字和字母的交替输出，体现并发执行特性")

	// 使用go关键字启动第一个goroutine
	// 这会立即返回，不会阻塞主程序
	go printNumbers()

	// 启动第二个goroutine
	// 现在有三个goroutine在运行：main + printNumbers + printLetters
	go printLetters()

	// 主goroutine需要等待子goroutine完成
	// 这里使用sleep是一种粗糙的方式，实际项目中应该使用WaitGroup
	fmt.Println("主程序等待goroutine完成...")
	time.Sleep(1 * time.Second)

	fmt.Println("程序结束")
	fmt.Println("注意：如果主程序提前结束，所有goroutine都会被强制终止")
}
