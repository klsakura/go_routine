package main

import (
	"fmt"
	"time"
)

// 基础select演示
func main() {
	fmt.Println("=== 基础Select演示 ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// 启动两个goroutine发送数据
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "来自channel1的消息"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "来自channel2的消息"
	}()

	// 使用select同时监听多个channel
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("接收到: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("接收到: %s\n", msg2)
		}
	}

	fmt.Println("程序结束")
}
