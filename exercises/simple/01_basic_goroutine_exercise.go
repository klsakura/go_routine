/*
Golang并发编程练习 - 简单级别
练习文件：01_basic_goroutine_exercise.go
练习主题：基础Goroutine使用

练习目标：
1. 创建和启动多个goroutine
2. 理解goroutine的并发执行特性
3. 掌握goroutine与主程序的同步
4. 观察并发执行的不确定性

练习任务：
- 任务1：创建两个函数，一个打印奇数1-9，一个打印偶数2-10
- 任务2：使用goroutine并发执行这两个函数
- 任务3：确保主程序等待所有goroutine完成
- 任务4：尝试不同的等待时间，观察效果

运行方式：go run exercises/simple/01_basic_goroutine_exercise.go
*/

package main

import (
	"fmt"
)

// printOddNumbers 打印奇数的函数
// TODO: 实现这个函数，打印1, 3, 5, 7, 9
// 提示：使用for循环，每次打印后睡眠100毫秒
func printOddNumbers() {
	// 在这里实现您的代码

}

// printEvenNumbers 打印偶数的函数
// TODO: 实现这个函数，打印2, 4, 6, 8, 10
// 提示：使用for循环，每次打印后睡眠150毫秒
func printEvenNumbers() {
	// 在这里实现您的代码

}

// 练习扩展函数
// printCountdown 倒计时函数
// TODO: 实现从n倒数到1的函数
func printCountdown(n int, name string) {
	// 在这里实现您的代码
	// 提示：从n开始倒数到1，每秒打印一个数字

}

func main() {
	fmt.Println("=== Goroutine基础练习 ===")

	// 任务1：基础并发执行
	fmt.Println("\n任务1：奇数偶数并发打印")
	// TODO: 使用goroutine启动printOddNumbers和printEvenNumbers

	// TODO: 等待goroutine完成（使用time.Sleep）

	fmt.Println("任务1完成\n")

	// 任务2：多个goroutine
	fmt.Println("任务2：多个倒计时")
	// TODO: 启动3个倒计时goroutine，分别从5、3、7开始倒数
	// 给每个倒计时一个名字，比如"Timer-A"、"Timer-B"、"Timer-C"

	// TODO: 等待所有倒计时完成

	fmt.Println("任务2完成\n")

	// 任务3：观察执行顺序
	fmt.Println("任务3：观察执行顺序")
	// TODO: 启动5个goroutine，每个都打印自己的ID（1-5）
	// 每个goroutine打印3次，观察输出顺序的随机性

	// TODO: 等待完成

	fmt.Println("所有练习完成！")

	// 反思问题：
	fmt.Println("\n思考题：")
	fmt.Println("1. 为什么每次运行程序，输出的顺序可能不同？")
	fmt.Println("2. 如果主程序提前结束会发生什么？")
	fmt.Println("3. time.Sleep是等待goroutine的好方法吗？为什么？")
}
