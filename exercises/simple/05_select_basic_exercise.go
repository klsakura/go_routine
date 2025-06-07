package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// TODO: 在这里实现你的代码

func main() {
	ch1 := make(chan string, 10)
	ch2 := make(chan string, 10)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("请输入内容：")
			if scanner.Scan() {
				text := scanner.Text()
				ch1 <- text
			} else {
				// 读取失败（例如 Ctrl+D），退出
				close(ch1)
				break
			}
		}
	}()

	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "1秒后的消息"
		// close(ch2)
	}()

	for {
		select {
		case msg := <-ch1:
			fmt.Printf(msg)
		case msg := <-ch2:
			fmt.Printf("收到通道2的数据了%s\n", msg)
		case <-time.After(5 * time.Second):
			fmt.Printf("暂时没有消息进入\n")
		}

	}

}
