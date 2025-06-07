/*
Golang并发编程练习 - 简单级别
练习文件：04_buffered_channel_exercise.go
练习主题：缓冲Channel使用

练习目标：
1. 理解缓冲channel与无缓冲channel的区别
2. 掌握缓冲区容量和长度的概念
3. 学会在生产者-消费者场景中应用缓冲channel
4. 理解缓冲channel的异步特性

练习任务：
- 任务1：实现基础的缓冲channel操作
- 任务2：创建生产者-消费者模型
- 任务3：实验不同缓冲区大小的影响
- 任务4：处理缓冲区满的情况

运行方式：go run exercises/simple/04_buffered_channel_exercise.go
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Message 消息结构
type Message struct {
	ID      int
	Content string
	Time    time.Time
}

// TODO: 实现缓冲区基础操作演示
func demonstrateBufferBasics() {
	fmt.Println("=== 缓冲区基础操作 ===")

	// TODO: 创建一个容量为3的string类型缓冲channel

	// TODO: 向channel发送数据，观察缓冲区长度变化
	// 发送3条消息："msg1", "msg2", "msg3"
	// 每次发送后打印：len(ch), cap(ch)

	// TODO: 尝试发送第4条消息（会阻塞），需要在goroutine中处理

	// TODO: 接收所有消息并打印

	fmt.Println("缓冲区基础操作完成\n")
}

// Producer 生产者函数
// TODO: 实现生产者，向channel发送指定数量的Message
func Producer(ch chan<- Message, count int, producerID string) {
	// TODO: 在这里实现您的代码
	// 提示：
	// 1. 循环创建count个消息
	// 2. 每个消息有唯一ID，包含生产者信息
	// 3. 模拟生产时间（50-200ms随机延迟）
	// 4. 发送完成后打印统计信息
}

// Consumer 消费者函数
// TODO: 实现消费者，从channel接收并处理Message
func Consumer(ch <-chan Message, consumerID string) {
	// TODO: 在这里实现您的代码
	// 提示：
	// 1. 使用range遍历channel
	// 2. 模拟处理时间（100-300ms随机延迟）
	// 3. 打印接收到的消息信息
	// 4. 统计处理的消息数量
}

// TODO: 实现测试不同缓冲区大小对性能的影响
func testBufferSizeImpact() {
	fmt.Println("=== 缓冲区大小影响测试 ===")

	bufferSizes := []int{1, 5, 10, 20}

	for _, size := range bufferSizes {
		fmt.Printf("\n--- 测试缓冲区大小: %d ---\n", size)

		// TODO: 为每个缓冲区大小进行测试
		// 1. 创建指定大小的缓冲channel
		// 2. 启动1个生产者和1个消费者
		// 3. 测量完成时间
		// 4. 对比不同缓冲区大小的性能差异
	}

	fmt.Println("缓冲区大小影响测试完成\n")
}

// TODO: 实现多生产者多消费者场景
func multiProducerConsumer() {
	fmt.Println("=== 多生产者多消费者测试 ===")

	// TODO: 在这里实现您的代码
	// 提示：
	// 1. 创建一个适当大小的缓冲channel
	// 2. 启动3个生产者，每个生产5条消息
	// 3. 启动2个消费者
	// 4. 使用sync.WaitGroup确保同步
	// 5. 观察消息的分布情况

	fmt.Println("多生产者多消费者测试完成\n")
}

// TODO: 实现缓冲区监控功能
func bufferMonitor(ch chan Message, interval time.Duration, duration time.Duration) {
	// TODO: 在这里实现您的代码
	// 提示：
	// 1. 定期检查channel的长度和容量
	// 2. 计算缓冲区使用率
	// 3. 打印监控信息
	// 4. 运行指定的时间后停止
}

func main() {
	fmt.Println("=== 缓冲Channel练习 ===")
	rand.Seed(time.Now().UnixNano())

	// 任务1：基础缓冲区操作
	demonstrateBufferBasics()

	// 任务2：生产者-消费者模型
	fmt.Println("任务2：生产者-消费者模型")
	// TODO: 实现以下场景：
	// 1. 创建容量为5的Message channel
	// 2. 启动1个生产者（生产10条消息）
	// 3. 启动1个消费者
	// 4. 使用适当的同步机制

	fmt.Println("任务2完成\n")

	// 任务3：缓冲区大小影响
	testBufferSizeImpact()

	// 任务4：多生产者多消费者
	multiProducerConsumer()

	// 任务5：缓冲区监控（挑战任务）
	fmt.Println("任务5：缓冲区监控")
	// TODO: 创建一个channel，启动生产者和消费者
	// 同时运行监控功能，观察缓冲区使用情况

	fmt.Println("任务5完成\n")

	fmt.Println("所有练习完成！")

	// 反思问题：
	fmt.Println("\n思考题：")
	fmt.Println("1. 缓冲channel和无缓冲channel的主要区别是什么？")
	fmt.Println("2. 如何选择合适的缓冲区大小？")
	fmt.Println("3. 缓冲区满时会发生什么？")
	fmt.Println("4. 在什么场景下缓冲channel比无缓冲channel更适合？")
}
