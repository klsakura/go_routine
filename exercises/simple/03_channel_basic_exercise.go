/*
Golang并发编程练习 - 简单级别
练习文件：03_channel_basic_exercise.go
练习主题：基础Channel通信

练习目标：
1. 掌握channel的创建和基本操作
2. 理解无缓冲channel的同步特性
3. 学会使用单向channel
4. 掌握channel的关闭和range遍历

练习任务：
- 任务1：实现简单的发送者和接收者
- 任务2：创建数据处理管道
- 任务3：使用channel实现简单的通信协议
- 任务4：处理多个channel的数据

运行方式：go run exercises/simple/03_channel_basic_exercise.go
*/

package main

import (
	"fmt"
)

// numberSender 数字发送者
// TODO: 实现发送1到n的数字到channel，然后关闭channel
func numberSender(ch chan<- int, n int) {
	// 在这里实现您的代码

}

// numberReceiver 数字接收者
// TODO: 实现从channel接收所有数字并计算总和
func numberReceiver(ch <-chan int) int {
	// 在这里实现您的代码
	// 提示：使用range遍历channel
	return 0
}

// stringProcessor 字符串处理器
// TODO: 实现接收字符串，处理后发送到输出channel
// 处理规则：添加前缀"Processed: "，转换为大写
func stringProcessor(input <-chan string, output chan<- string) {
	// 在这里实现您的代码

}

// dataGenerator 数据生成器
// TODO: 生成指定数量的随机数据并发送到channel
func dataGenerator(ch chan<- string, count int) {
	// 在这里实现您的代码
	// 提示：生成格式为"data-001", "data-002"的字符串

}

// messageRouter 消息路由器
// TODO: 从输入channel接收消息，根据消息内容路由到不同的输出channel
func messageRouter(input <-chan string, evenCh, oddCh chan<- string) {
	// 在这里实现您的代码
	// 提示：如果消息包含偶数，发送到evenCh；包含奇数，发送到oddCh

}

func main() {
	fmt.Println("=== Channel基础练习 ===")

	// 任务1：基础发送和接收
	fmt.Println("\n任务1：数字发送和接收")
	// TODO: 创建一个int类型的channel
	// 启动numberSender goroutine发送1到10的数字
	// 在主goroutine中接收并打印总和

	fmt.Println("任务1完成\n")

	// 任务2：数据处理管道
	fmt.Println("任务2：字符串处理管道")
	// TODO: 创建输入和输出channel
	// 启动stringProcessor goroutine
	// 发送一些字符串："hello", "world", "golang"
	// 接收并打印处理后的结果

	fmt.Println("任务2完成\n")

	// 任务3：多生产者单消费者
	fmt.Println("任务3：多生产者单消费者")
	// TODO: 创建一个共享的channel
	// 启动3个dataGenerator goroutine，每个生成5个数据
	// 在主goroutine中接收所有数据并打印

	fmt.Println("任务3完成\n")

	// 任务4：消息路由
	fmt.Println("任务4：消息路由系统")
	// TODO: 创建输入channel和两个输出channel(evenCh, oddCh)
	// 启动messageRouter goroutine
	// 启动两个接收者goroutine分别处理偶数和奇数消息
	// 发送消息："msg-1", "msg-2", ..., "msg-10"

	fmt.Println("任务4完成\n")

	// 任务5：channel关闭检测
	fmt.Println("任务5：检测channel是否关闭")
	// TODO: 创建一个channel，发送一些数据然后关闭
	// 使用两种方法检测channel关闭：
	// 方法1：使用ok模式 (val, ok := <-ch)
	// 方法2：使用range

	fmt.Println("所有练习完成！")

	// 反思问题：
	fmt.Println("\n思考题：")
	fmt.Println("1. 无缓冲channel的发送操作什么时候会阻塞？")
	fmt.Println("2. 单向channel有什么作用？")
	fmt.Println("3. 忘记关闭channel会有什么后果？")
	fmt.Println("4. 如何避免从已关闭的channel发送数据？")
}
