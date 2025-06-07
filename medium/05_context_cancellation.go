package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Context取消演示
func longRunningTask(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("任务 %d 被取消: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("任务 %d 执行步骤 %d\n", id, i)
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Printf("任务 %d 正常完成\n", id)
}

func taskWithTimeout(ctx context.Context, id int, duration time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	fmt.Printf("任务 %d 开始，超时时间: %v\n", id, duration)

	select {
	case <-time.After(2 * time.Second):
		fmt.Printf("任务 %d 完成工作\n", id)
	case <-ctx.Done():
		fmt.Printf("任务 %d 超时: %v\n", id, ctx.Err())
	}
}

func main() {
	fmt.Println("=== Context取消演示 ===")

	// 示例1: 手动取消
	fmt.Println("\n1. 手动取消演示:")
	ctx1, cancel1 := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	// 启动3个长时间运行的任务
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go longRunningTask(ctx1, i, &wg)
	}

	// 2秒后取消所有任务
	time.Sleep(2 * time.Second)
	fmt.Println("取消所有任务...")
	cancel1()

	wg.Wait()

	// 示例2: 超时取消
	fmt.Println("\n2. 超时取消演示:")

	// 启动几个不同超时时间的任务
	go taskWithTimeout(context.Background(), 1, 1*time.Second) // 会超时
	go taskWithTimeout(context.Background(), 2, 3*time.Second) // 会完成

	time.Sleep(4 * time.Second)

	// 示例3: 带deadline的取消
	fmt.Println("\n3. Deadline取消演示:")
	deadline := time.Now().Add(1500 * time.Millisecond)
	ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
	defer cancel3()

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go longRunningTask(ctx3, 99, &wg2)

	wg2.Wait()

	fmt.Println("Context演示完成！")
}
