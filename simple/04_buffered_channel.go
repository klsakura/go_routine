/*
Golang并发编程学习Demo - 简单级别
文件：04_buffered_channel.go
主题：缓冲Channel使用

本示例演示：
1. 缓冲channel与无缓冲channel的区别
2. 缓冲channel的容量和长度概念
3. 异步发送和接收操作
4. 缓冲满时的阻塞行为
5. 缓冲channel的实际应用场景

学习要点：
- 缓冲channel允许异步通信
- 发送操作在缓冲区满之前不会阻塞
- 接收操作在缓冲区空之前不会阻塞
- len()获取当前缓冲区中的元素数量
- cap()获取缓冲区的总容量

运行方式：go run simple/04_buffered_channel.go
*/

package main

import (
	"fmt"
	"time"
)

// demonstrateBasicBufferedChannel 演示基础缓冲channel操作
func demonstrateBasicBufferedChannel() {
	fmt.Println("=== 基础缓冲Channel演示 ===")

	// 创建一个容量为3的缓冲channel
	bufferedCh := make(chan string, 3)

	fmt.Printf("初始状态 - 长度: %d, 容量: %d\n", len(bufferedCh), cap(bufferedCh))

	// 向缓冲channel发送数据（不会阻塞，因为还有空间）
	fmt.Println("发送数据到缓冲channel...")
	bufferedCh <- "第一条消息"
	fmt.Printf("发送后 - 长度: %d, 容量: %d\n", len(bufferedCh), cap(bufferedCh))

	bufferedCh <- "第二条消息"
	fmt.Printf("发送后 - 长度: %d, 容量: %d\n", len(bufferedCh), cap(bufferedCh))

	bufferedCh <- "第三条消息"
	fmt.Printf("发送后 - 长度: %d, 容量: %d\n", len(bufferedCh), cap(bufferedCh))

	// 现在缓冲区已满，再发送会阻塞（在goroutine中演示）
	go func() {
		fmt.Println("尝试发送第四条消息（会阻塞直到有空间）...")
		bufferedCh <- "第四条消息"
		fmt.Println("第四条消息发送成功！")
	}()

	// 接收消息，为第四条消息腾出空间
	time.Sleep(500 * time.Millisecond) // 让goroutine有时间尝试发送

	fmt.Println("\n开始接收消息...")
	for i := 0; i < 4; i++ {
		msg := <-bufferedCh
		fmt.Printf("接收到: %s (剩余: %d)\n", msg, len(bufferedCh))
		time.Sleep(200 * time.Millisecond)
	}

	close(bufferedCh)
}

// producer 生产者函数，向channel发送数据
func producer(ch chan<- int, start, count int) {
	defer close(ch) // 确保发送完成后关闭channel

	fmt.Printf("生产者开始工作：从%d开始，生产%d个数字\n", start, count)

	for i := 0; i < count; i++ {
		value := start + i

		// 检查发送前的缓冲区状态
		fmt.Printf("准备发送: %d (当前缓冲区长度: %d)\n", value, len(ch))

		ch <- value

		// 模拟生产时间
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("生产者完成工作")
}

// consumer 消费者函数，从channel接收数据
func consumer(ch <-chan int, consumerID string) {
	fmt.Printf("消费者 %s 开始工作\n", consumerID)

	messageCount := 0
	for value := range ch { // range会自动处理channel关闭
		messageCount++
		fmt.Printf("消费者 %s 接收到: %d (第%d条消息)\n", consumerID, value, messageCount)

		// 模拟处理时间
		time.Sleep(150 * time.Millisecond)
	}

	fmt.Printf("消费者 %s 完成工作，共处理 %d 条消息\n", consumerID, messageCount)
}

// demonstrateProducerConsumer 演示生产者-消费者模式
func demonstrateProducerConsumer() {
	fmt.Println("\n=== 生产者-消费者模式演示 ===")

	// 创建一个容量为5的缓冲channel
	// 这样生产者可以在消费者准备好之前发送一些数据
	productChannel := make(chan int, 5)

	// 启动生产者goroutine
	go producer(productChannel, 100, 10)

	// 启动消费者goroutine
	go consumer(productChannel, "A")

	// 等待足够的时间让生产者和消费者完成工作
	time.Sleep(3 * time.Second)
}

// demonstrateMultipleConsumers 演示多个消费者竞争同一个缓冲channel
func demonstrateMultipleConsumers() {
	fmt.Println("\n=== 多消费者竞争演示 ===")

	// 创建一个较大的缓冲channel
	jobChannel := make(chan int, 8)

	// 启动多个消费者
	for i := 1; i <= 3; i++ {
		consumerID := fmt.Sprintf("Consumer-%d", i)
		go consumer(jobChannel, consumerID)
	}

	// 生产者发送任务
	go func() {
		defer close(jobChannel)

		fmt.Println("发送10个任务到job channel...")
		for i := 1; i <= 10; i++ {
			fmt.Printf("发送任务: %d\n", i)
			jobChannel <- i
			time.Sleep(80 * time.Millisecond) // 控制发送速度
		}
		fmt.Println("所有任务发送完成")
	}()

	// 等待所有消费者处理完成
	time.Sleep(4 * time.Second)
}

// demonstrateChannelCapacityEffects 演示不同缓冲区大小的影响
func demonstrateChannelCapacityEffects() {
	fmt.Println("\n=== 缓冲区大小影响演示 ===")

	// 测试不同大小的缓冲区
	capacities := []int{1, 3, 10}

	for _, cap := range capacities {
		fmt.Printf("\n--- 测试容量为 %d 的缓冲channel ---\n", cap)

		ch := make(chan string, cap)

		// 记录开始时间
		start := time.Now()

		// 启动接收者（延迟启动以观察缓冲效果）
		go func() {
			time.Sleep(500 * time.Millisecond) // 延迟接收
			for i := 0; i < 5; i++ {
				msg := <-ch
				fmt.Printf("接收: %s\n", msg)
				time.Sleep(100 * time.Millisecond)
			}
		}()

		// 发送数据
		for i := 1; i <= 5; i++ {
			msg := fmt.Sprintf("消息-%d", i)
			fmt.Printf("发送: %s (时间: %v)\n", msg, time.Since(start))
			ch <- msg
		}

		// 等待处理完成
		time.Sleep(1 * time.Second)
		close(ch)
	}
}

func main() {
	fmt.Println("=== 缓冲Channel详细演示 ===")
	fmt.Println("观察缓冲channel如何改变goroutine间的通信行为\n")

	// 1. 基础缓冲channel操作
	demonstrateBasicBufferedChannel()

	// 2. 生产者-消费者模式
	demonstrateProducerConsumer()

	// 3. 多消费者竞争
	demonstrateMultipleConsumers()

	// 4. 缓冲区大小对性能的影响
	demonstrateChannelCapacityEffects()

	fmt.Println("\n=== 总结 ===")
	fmt.Println("缓冲channel的优势：")
	fmt.Println("1. 减少goroutine阻塞，提高并发性能")
	fmt.Println("2. 解耦生产者和消费者的处理速度")
	fmt.Println("3. 提供临时存储，平滑处理峰值")
	fmt.Println("4. 在生产者-消费者模式中特别有用")

	fmt.Println("\n注意事项：")
	fmt.Println("1. 缓冲区太小可能导致频繁阻塞")
	fmt.Println("2. 缓冲区太大可能浪费内存")
	fmt.Println("3. 需要根据实际场景选择合适的缓冲区大小")
	fmt.Println("4. 记住关闭channel以避免goroutine泄露")
}
