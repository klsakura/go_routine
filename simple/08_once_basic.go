package main

import (
	"fmt"
	"sync"
	"time"
)

// sync.Once演示
var once sync.Once

func initialize() {
	fmt.Println("执行初始化操作...")
	time.Sleep(1 * time.Second)
	fmt.Println("初始化完成！")
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("工作者 %d 尝试初始化\n", id)

	// 确保initialize只执行一次
	once.Do(initialize)

	fmt.Printf("工作者 %d 继续执行自己的任务\n", id)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("工作者 %d 完成\n", id)
}

func main() {
	fmt.Println("=== sync.Once演示 ===")

	var wg sync.WaitGroup

	// 启动5个工作者，每个都尝试初始化
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("所有工作者完成！")
}
