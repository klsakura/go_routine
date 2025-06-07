/*
Golang并发编程学习Demo - 简单级别
文件：03_channel_basic.go
主题：基础Channel通信

本示例演示：
1. 无缓冲channel的创建和使用
2. goroutine间的数据传递
3. channel的发送和接收操作
4. 单向channel的概念
5. 使用range遍历channel
6. 关闭channel的重要性

学习要点：
- channel是goroutine间通信的管道
- <- 操作符用于发送和接收
- 无缓冲channel是同步的
- 关闭channel通知接收者没有更多数据
- range可以自动检测channel关闭

运行方式：go run simple/03_channel_basic.go
*/

package main

import (
	"fmt"
	"time"
)

// sender 发送者函数，向channel发送数据
// ch: 只写channel，只能向其发送数据
// chan<- string 是单向channel类型，提高代码安全性
func sender(ch chan<- string) {
	messages := []string{"Hello", "World", "Go", "Channel"}

	fmt.Println("开始发送消息...")
	for i, msg := range messages {
		fmt.Printf("发送第%d条消息: %s\n", i+1, msg)

		// 使用 <- 操作符向channel发送数据
		// 对于无缓冲channel，这会阻塞直到有接收者
		ch <- msg

		// 模拟一些处理时间
		time.Sleep(500 * time.Millisecond)
	}

	// 关闭channel非常重要！
	// 这会通知接收者没有更多数据会被发送
	close(ch)
	fmt.Println("发送完成，channel已关闭")
}

func main() {
	fmt.Println("=== 基础Channel演示 ===")
	fmt.Println("演示goroutine间通过channel进行通信")

	// 创建一个string类型的无缓冲channel
	// 无缓冲channel必须有发送者和接收者同时准备好才能完成传输
	ch := make(chan string)

	// 启动发送者goroutine
	// 必须在goroutine中运行，否则会死锁
	// 因为无缓冲channel的发送操作会阻塞直到有接收者
	fmt.Println("启动发送者goroutine...")
	go sender(ch)

	// 在主goroutine中接收消息
	// 使用range遍历channel，会自动处理channel关闭
	fmt.Println("开始接收消息...")
	messageCount := 0
	for msg := range ch {
		messageCount++
		fmt.Printf("接收第%d条消息: %s\n", messageCount, msg)
	}

	fmt.Printf("通信完成！共接收到%d条消息\n", messageCount)
	fmt.Println("注意：range循环在channel关闭时自动退出")
}
