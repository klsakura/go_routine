/*
Golang并发编程练习 - 简单级别
练习文件：02_waitgroup_basic_exercise.go
练习主题：WaitGroup基础使用

练习目标：
1. 掌握sync.WaitGroup的基本用法
2. 理解Add()、Done()、Wait()的作用
3. 学会正确同步多个goroutine
4. 避免goroutine泄露

练习任务：
- 任务1：使用WaitGroup管理多个工作者goroutine
- 任务2：实现并发下载模拟器
- 任务3：处理可变数量的goroutine
- 任务4：在goroutine中使用defer确保Done()被调用

运行方式：go run exercises/simple/02_waitgroup_basic_exercise.go
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// simpleWorker 简单工作者函数
// TODO: 实现这个函数，接收工作者ID和WaitGroup指针
// 提示：使用defer wg.Done()，模拟工作时间，打印开始和完成信息
func simpleWorker(id int, wg *sync.WaitGroup) {
	// 在这里实现您的代码

}

// downloadFile 模拟文件下载
// TODO: 实现文件下载模拟，接收文件名、大小(MB)和WaitGroup
// 提示：根据文件大小计算下载时间，显示下载进度
func downloadFile(filename string, sizeMB int, wg *sync.WaitGroup) {
	// 在这里实现您的代码
	// 建议：每MB需要100毫秒下载时间，显示下载进度

}

// processTask 处理任务函数
// TODO: 实现任务处理函数，接收任务ID、处理时间和WaitGroup
// 模拟可能失败的任务（20%失败率）
func processTask(taskID string, processingTime time.Duration, wg *sync.WaitGroup) {
	// 在这里实现您的代码
	// 提示：使用rand.Float32() < 0.2 模拟失败

}

func main() {
	fmt.Println("=== WaitGroup基础练习 ===")
	rand.Seed(time.Now().UnixNano())

	// 任务1：基础WaitGroup使用
	fmt.Println("\n任务1：管理多个工作者")
	// TODO: 创建WaitGroup，启动5个simpleWorker
	// 每个工作者有不同的ID(1-5)

	fmt.Println("任务1完成\n")

	// 任务2：模拟并发下载
	fmt.Println("任务2：并发文件下载")
	// TODO: 创建WaitGroup，模拟下载以下文件：
	files := []struct {
		name string
		size int // MB
	}{
		{"video.mp4", 100},
		{"music.mp3", 5},
		{"document.pdf", 2},
		{"image.jpg", 1},
		{"archive.zip", 50},
	}

	// TODO: 为每个文件启动下载goroutine

	fmt.Println("任务2完成\n")

	// 任务3：动态数量的goroutine
	fmt.Println("任务3：处理可变数量的任务")
	// TODO: 创建5-15个随机数量的任务
	taskCount := rand.Intn(11) + 5 // 5到15个任务
	fmt.Printf("需要处理 %d 个任务\n", taskCount)

	// TODO: 为每个任务启动goroutine，使用processTask函数
	// 任务ID格式："task-001", "task-002"等
	// 处理时间随机100-500毫秒

	fmt.Println("任务3完成\n")

	// 任务4：错误处理练习
	fmt.Println("任务4：WaitGroup错误处理")
	// TODO: 故意制造一个常见错误然后修复
	// 比如：忘记调用Done()、Add()和Done()数量不匹配等
	// 请在注释中写出可能的错误和解决方案

	fmt.Println("所有练习完成！")

	// 反思问题：
	fmt.Println("\n思考题：")
	fmt.Println("1. 为什么要使用defer wg.Done()？")
	fmt.Println("2. 如果Add()和Done()的数量不匹配会怎样？")
	fmt.Println("3. WaitGroup相比time.Sleep有什么优势？")
	fmt.Println("4. 在什么情况下goroutine可能会泄露？")
}
