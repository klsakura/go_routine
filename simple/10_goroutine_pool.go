package main

import (
	"fmt"
	"sync"
	"time"
)

// 简单的goroutine池演示
type Job struct {
	ID   int
	Data string
}

type Result struct {
	Job Job
	Sum int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("工作者 %d 处理任务 %d\n", id, job.ID)

		// 模拟工作：计算字符串长度
		time.Sleep(500 * time.Millisecond)
		sum := len(job.Data)

		results <- Result{job, sum}
		fmt.Printf("工作者 %d 完成任务 %d\n", id, job.ID)
	}
}

func main() {
	fmt.Println("=== 简单Goroutine池演示 ===")

	const numWorkers = 3
	const numJobs = 5

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	// 启动工作者
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 发送任务
	jobData := []string{"hello", "world", "golang", "concurrency", "programming"}
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Data: jobData[j-1]}
	}
	close(jobs)

	// 等待所有工作者完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	fmt.Println("\n结果:")
	for result := range results {
		fmt.Printf("任务 %d (%s) -> 长度: %d\n",
			result.Job.ID, result.Job.Data, result.Sum)
	}

	fmt.Println("所有任务完成！")
}
