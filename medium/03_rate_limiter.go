package main

import (
	"fmt"
	"sync"
	"time"
)

// 速率限制器演示
type RateLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	done   chan bool
}

func NewRateLimiter(rate int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, rate),
		ticker: time.NewTicker(time.Second / time.Duration(rate)),
		done:   make(chan bool),
	}

	// 初始填充tokens
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}

	// 启动token补充器
	go rl.refill()

	return rl
}

func (rl *RateLimiter) refill() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
				// token池已满，跳过
			}
		case <-rl.done:
			return
		}
	}
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

func (rl *RateLimiter) Close() {
	rl.ticker.Stop()
	close(rl.done)
}

func worker(id int, limiter *RateLimiter, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Printf("工作者 %d: 请求 %d 通过 (时间: %s)\n",
				id, i, time.Now().Format("15:04:05.000"))
		} else {
			fmt.Printf("工作者 %d: 请求 %d 被限制，等待中...\n", id, i)
			limiter.Wait()
			fmt.Printf("工作者 %d: 请求 %d 重试成功 (时间: %s)\n",
				id, i, time.Now().Format("15:04:05.000"))
		}

		// 模拟处理时间
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("工作者 %d 完成所有请求\n", id)
}

func main() {
	fmt.Println("=== 速率限制器演示 ===")
	fmt.Println("限制: 每秒最多5个请求")

	// 创建每秒5个请求的速率限制器
	limiter := NewRateLimiter(5)
	defer limiter.Close()

	var wg sync.WaitGroup

	// 启动3个工作者，每个尝试发送10个请求
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, limiter, &wg)
	}

	wg.Wait()
	fmt.Println("所有工作者完成！")
}
