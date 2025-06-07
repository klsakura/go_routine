/*
Golang并发编程学习Demo - 简单级别
文件：02_waitgroup_basic.go
主题：WaitGroup基础使用

本示例演示：
1. sync.WaitGroup的基本用法
2. 如何正确等待多个goroutine完成
3. Add()、Done()、Wait()方法的作用
4. 解决goroutine同步问题的标准方案

学习要点：
- WaitGroup比time.Sleep更可靠
- Add()增加等待计数
- Done()减少等待计数
- Wait()阻塞直到计数为0
- defer确保Done()一定被调用

运行方式：go run simple/02_waitgroup_basic.go
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// worker 工作者函数，模拟一个需要时间的任务
// id: 工作者编号，用于标识不同的goroutine
// wg: WaitGroup指针，用于通知任务完成
func worker(id int, wg *sync.WaitGroup) {
	// defer确保函数结束时一定会调用Done()
	// 即使函数中途panic也会执行
	defer wg.Done() // 任务完成时调用Done()，计数器-1

	fmt.Printf("工作者 %d 开始工作\n", id)

	// 模拟不同的工作时间，让效果更明显
	// 工作者id越大，工作时间越长
	time.Sleep(time.Duration(id) * 100 * time.Millisecond)

	fmt.Printf("工作者 %d 完成工作\n", id)
}

func main() {
	fmt.Println("=== WaitGroup基础演示 ===")
	fmt.Println("演示如何使用WaitGroup等待多个goroutine完成")

	// 创建WaitGroup实例
	// WaitGroup内部维护一个计数器
	var wg sync.WaitGroup

	// 启动5个工作者goroutine
	fmt.Println("启动5个工作者...")
	for i := 1; i <= 5; i++ {
		// 每启动一个goroutine前，先Add(1)增加计数器
		wg.Add(1) // 计数器+1，告诉WaitGroup要等待一个goroutine

		// 启动goroutine，传递工作者ID和WaitGroup指针
		go worker(i, &wg)
	}

	// Wait()会阻塞当前goroutine（main goroutine）
	// 直到WaitGroup内部计数器变为0
	fmt.Println("主程序等待所有工作者完成...")
	wg.Wait()

	fmt.Println("所有工作者完成！")
	fmt.Println("程序正常结束，没有goroutine泄露")
}
