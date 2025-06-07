package main

import (
	"fmt"
	"sync"
)

// 基础互斥锁演示
var (
	counter int
	mutex   sync.Mutex
)

func increment(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		mutex.Lock()   // 加锁
		counter++      // 临界区
		mutex.Unlock() // 解锁
	}
}

func main() {
	fmt.Println("=== 基础互斥锁演示 ===")

	var wg sync.WaitGroup

	// 启动3个goroutine同时增加计数器
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go increment(&wg)
	}

	wg.Wait()

	fmt.Printf("最终计数器值: %d (期望值: 3000)\n", counter)

	// 演示不使用锁的情况
	counter = 0
	fmt.Println("\n不使用锁的情况:")

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter++ // 没有锁保护，可能出现竞态条件
			}
		}()
	}

	wg.Wait()
	fmt.Printf("不安全的计数器值: %d (可能不等于3000)\n", counter)
}
