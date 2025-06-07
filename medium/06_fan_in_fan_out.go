package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 扇入扇出模式演示
type Task struct {
	ID   int
	Data int
}

type Result struct {
	TaskID int
	Result int
}

// 生成器：生成任务
func taskGenerator(tasks []Task) <-chan Task {
	out := make(chan Task)
	go func() {
		defer close(out)
		for _, task := range tasks {
			out <- task
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return out
}

// 扇出：将任务分发给多个工作者
func fanOut(input <-chan Task, numWorkers int) []<-chan Result {
	workers := make([]<-chan Result, numWorkers)

	for i := 0; i < numWorkers; i++ {
		worker := make(chan Result)
		workers[i] = worker

		go func(workerID int, input <-chan Task, output chan<- Result) {
			defer close(output)
			for task := range input {
				// 模拟处理时间
				processingTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
				time.Sleep(processingTime)

				// 计算结果（这里简单地平方）
				result := Result{
					TaskID: task.ID,
					Result: task.Data * task.Data,
				}

				fmt.Printf("工作者 %d 处理任务 %d: %d -> %d\n",
					workerID, task.ID, task.Data, result.Result)

				output <- result
			}
			fmt.Printf("工作者 %d 完成所有任务\n", workerID)
		}(i+1, input, worker)
	}

	return workers
}

// 扇入：将多个工作者的结果合并
func fanIn(inputs []<-chan Result) <-chan Result {
	output := make(chan Result)
	var wg sync.WaitGroup

	// 为每个输入channel启动一个goroutine进行合并
	for i, input := range inputs {
		wg.Add(1)
		go func(id int, ch <-chan Result) {
			defer wg.Done()
			for result := range ch {
				fmt.Printf("从工作者通道 %d 收集结果: 任务%d = %d\n",
					id+1, result.TaskID, result.Result)
				output <- result
			}
		}(i, input)
	}

	// 等待所有输入完成后关闭输出channel
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func main() {
	fmt.Println("=== 扇入扇出模式演示 ===")

	rand.Seed(time.Now().UnixNano())

	// 创建任务
	tasks := []Task{
		{ID: 1, Data: 2},
		{ID: 2, Data: 3},
		{ID: 3, Data: 4},
		{ID: 4, Data: 5},
		{ID: 5, Data: 6},
		{ID: 6, Data: 7},
		{ID: 7, Data: 8},
		{ID: 8, Data: 9},
	}

	fmt.Printf("总共 %d 个任务需要处理\n", len(tasks))

	// 创建任务流
	taskStream := taskGenerator(tasks)

	// 扇出：将任务分发给3个工作者
	const numWorkers = 3
	fmt.Printf("启动 %d 个工作者\n", numWorkers)
	workerChannels := fanOut(taskStream, numWorkers)

	// 扇入：合并工作者的结果
	results := fanIn(workerChannels)

	// 收集最终结果
	fmt.Println("\n最终结果:")
	var totalResults []Result
	for result := range results {
		totalResults = append(totalResults, result)
	}

	// 按任务ID排序显示结果
	fmt.Println("\n按任务ID排序的结果:")
	for i := 1; i <= len(tasks); i++ {
		for _, result := range totalResults {
			if result.TaskID == i {
				fmt.Printf("任务 %d: %d\n", result.TaskID, result.Result)
				break
			}
		}
	}

	fmt.Printf("\n处理完成！共处理 %d 个任务\n", len(totalResults))
}
