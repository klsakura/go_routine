// 多个消费者
package main

import (
	"fmt"
	"time"
)

func consumer(ch <-chan string, i int) {
	for data := range ch {
		fmt.Printf("消费者%d 消费了 %s\n", i, data)
	}
}

func main() {

	ch := make(chan string, 8)
	for i := 0; i < 3; i++ {
		go consumer(ch, i)
	}

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- fmt.Sprintf("数据%d", i)
			time.Sleep(80 * time.Millisecond) // 控制发送速度
		}
	}()

	time.Sleep(4 * time.Second)

}
