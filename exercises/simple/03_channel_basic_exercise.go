package main

import (
	"fmt"
)

func sender(ch chan<- string) {
	messages := []string{"hello", "world", "hehehe"}
	for _, v := range messages {
		ch <- v
	}
	close(ch)
}

func main() {
	//定义一个字符串通道
	ch := make(chan string)

	//启动一个发送者
	go sender(ch)

	//主程序接受ch的数据
	for msg := range ch {
		fmt.Printf("接受到数据%s\n", msg)
	}

}
