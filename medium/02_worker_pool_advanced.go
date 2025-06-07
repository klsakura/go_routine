package main

import (
	"fmt"
	"sync"
	"time"
)

// 高级工作池演示
type Task struct {
	ID       int
	Data     []int
	Priority int
}

type TaskResult struct {
	TaskID int
	Sum    int
	Worker int
}

type WorkerPool struct {
	tasks   chan Task
	results chan TaskResult
	workers int
	wg      sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		tasks:   make(chan Task, 100),
		results: make(chan TaskResult, 100),
		workers: numWorkers,
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for task := range wp.tasks {
		fmt.Printf("工作者 %d 开始处理任务 %d (优先级: %d)\n",
			id, task.ID, task.Priority)

		// 模拟处理时间，优先级高的任务处理更快
		processingTime := time.Duration(500-task.Priority*100) * time.Millisecond
		time.Sleep(processingTime)

		// 计算数组和
		sum := 0
		for _, v := range task.Data {
			sum += v
		}

		result := TaskResult{
			TaskID: task.ID,
			Sum:    sum,
			Worker: id,
		}

		wp.results <- result
		fmt.Printf("工作者 %d 完成任务 %d，结果: %d\n", id, task.ID, sum)
	}

	fmt.Printf("工作者 %d 退出\n", id)
}

func (wp *WorkerPool) Start() {
	// 启动工作者
	for i := 1; i <= wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Submit(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) Close() {
	close(wp.tasks)
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) Results() <-chan TaskResult {
	return wp.results
}

func main() {
	fmt.Println("=== 高级工作池演示 ===")

	pool := NewWorkerPool(3)
	pool.Start()

	// 提交不同优先级的任务
	tasks := []Task{
		{ID: 1, Data: []int{1, 2, 3, 4}, Priority: 1},
		{ID: 2, Data: []int{5, 6, 7, 8}, Priority: 3},
		{ID: 3, Data: []int{9, 10, 11}, Priority: 2},
		{ID: 4, Data: []int{12, 13, 14, 15}, Priority: 1},
		{ID: 5, Data: []int{16, 17}, Priority: 3},
	}

	for _, task := range tasks {
		pool.Submit(task)
	}

	pool.Close()

	// 启动结果收集器
	go func() {
		pool.Wait()
	}()

	// 收集结果
	fmt.Println("\n任务执行结果:")
	for result := range pool.Results() {
		fmt.Printf("任务 %d: 和=%d (由工作者 %d 完成)\n",
			result.TaskID, result.Sum, result.Worker)
	}

	fmt.Println("所有任务完成！")
}
