package main

import (
	"fmt"
	"time"
)

// 带超时的select演示
func main() {
	fmt.Println("=== 带超时的Select演示 ===")

	ch := make(chan string)

	// 启动一个goroutine，3秒后发送消息
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "延迟消息"
	}()

	fmt.Println("等待消息（超时时间：2秒）...")

	select {
	case msg := <-ch:
		fmt.Printf("接收到消息: %s\n", msg)
	case <-time.After(2 * time.Second):
		fmt.Println("超时！没有接收到消息")
	}

	fmt.Println("程序结束")
}
