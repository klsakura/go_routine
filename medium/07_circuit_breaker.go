/*
Golang并发编程学习Demo - 中等级别
文件：07_circuit_breaker.go
主题：熔断器模式

本示例演示：
1. 熔断器的三种状态：CLOSED、OPEN、HALF_OPEN
2. 失败率阈值触发熔断
3. 超时后的自动恢复机制
4. 半开状态的试探性调用
5. 保护不稳定服务的策略

核心概念：
- CLOSED：正常状态，请求正常通过
- OPEN：熔断状态，快速失败，不调用服务
- HALF_OPEN：半开状态，允许少量请求试探服务是否恢复

应用场景：
- 微服务架构中的服务保护
- 防止故障服务拖垮整个系统
- 提供快速失败机制

运行方式：go run medium/07_circuit_breaker.go
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// CircuitBreakerState 熔断器状态枚举
type CircuitBreakerState int

const (
	StateClosed   CircuitBreakerState = iota // 关闭状态：正常通过请求
	StateOpen                                // 开启状态：拒绝请求，快速失败
	StateHalfOpen                            // 半开状态：允许少量请求试探
)

// String 实现Stringer接口，便于打印状态
func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig 熔断器配置参数
type CircuitBreakerConfig struct {
	MaxFailures     int           // 最大失败次数（暂未使用）
	ResetTimeout    time.Duration // 从OPEN到HALF_OPEN的等待时间
	FailureRatio    float64       // 失败率阈值（0.0-1.0）
	MinRequestCount int           // 最小请求数，低于此数不触发熔断
}

// CircuitBreaker 熔断器核心结构
type CircuitBreaker struct {
	config       CircuitBreakerConfig // 配置参数
	state        CircuitBreakerState  // 当前状态
	failures     int64                // 失败计数（原子操作）
	requests     int64                // 请求计数（原子操作）
	lastFailTime time.Time            // 最后失败时间，用于计算重置时间
	mu           sync.RWMutex         // 读写锁，保护状态变更
}

// NewCircuitBreaker 创建新的熔断器实例
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed, // 初始状态为关闭
	}
}

// Call 执行被保护的函数调用
// fn: 需要被保护的函数，返回error表示成功或失败
func (cb *CircuitBreaker) Call(fn func() error) error {
	// 首先检查是否允许请求通过
	if !cb.allowRequest() {
		return fmt.Errorf("circuit breaker is OPEN")
	}

	// 增加请求计数（用defer确保一定执行）
	defer func() {
		atomic.AddInt64(&cb.requests, 1)
	}()

	// 执行实际的业务函数
	err := fn()

	// 根据执行结果更新熔断器状态
	if err != nil {
		cb.onFailure() // 处理失败情况
		return err
	}

	cb.onSuccess() // 处理成功情况
	return nil
}

// allowRequest 检查当前状态是否允许请求通过
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		// 关闭状态：允许所有请求
		return true
	case StateOpen:
		// 开启状态：检查是否达到重置时间
		if time.Since(cb.lastFailTime) > cb.config.ResetTimeout {
			// 达到重置时间，转为半开状态
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = StateHalfOpen
			cb.mu.Unlock()
			cb.mu.RLock()
			fmt.Println("熔断器状态: OPEN -> HALF_OPEN")
			return true
		}
		return false
	case StateHalfOpen:
		// 半开状态：允许请求，用于试探服务是否恢复
		return true
	default:
		return false
	}
}

// onSuccess 处理成功调用
func (cb *CircuitBreaker) onSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// 如果当前是半开状态，成功调用说明服务恢复，转为关闭状态
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		cb.failures = 0 // 重置失败计数
		fmt.Println("熔断器状态: HALF_OPEN -> CLOSED")
	}
}

// onFailure 处理失败调用
func (cb *CircuitBreaker) onFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// 增加失败计数并记录失败时间
	atomic.AddInt64(&cb.failures, 1)
	cb.lastFailTime = time.Now()

	// 如果当前是半开状态，失败说明服务仍有问题，立即转为开启状态
	if cb.state == StateHalfOpen {
		cb.state = StateOpen
		fmt.Println("熔断器状态: HALF_OPEN -> OPEN")
		return
	}

	// 如果当前是关闭状态，检查是否需要触发熔断
	if cb.state == StateClosed {
		requests := atomic.LoadInt64(&cb.requests)
		failures := atomic.LoadInt64(&cb.failures)

		// 只有在请求数达到最小值时才考虑熔断
		if requests >= int64(cb.config.MinRequestCount) {
			failureRatio := float64(failures) / float64(requests)
			if failureRatio >= cb.config.FailureRatio {
				cb.state = StateOpen
				fmt.Printf("熔断器状态: CLOSED -> OPEN (失败率: %.2f%%)\n", failureRatio*100)
			}
		}
	}
}

// GetState 获取当前状态（线程安全）
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetStats 获取统计信息（线程安全）
func (cb *CircuitBreaker) GetStats() (int64, int64, CircuitBreakerState) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return atomic.LoadInt64(&cb.requests), atomic.LoadInt64(&cb.failures), cb.state
}

// UnstableService 模拟不稳定的外部服务
type UnstableService struct {
	failureRate float64      // 失败率（0.0-1.0）
	mu          sync.RWMutex // 保护失败率的并发修改
}

// NewUnstableService 创建不稳定服务实例
func NewUnstableService(failureRate float64) *UnstableService {
	return &UnstableService{
		failureRate: failureRate,
	}
}

// Call 模拟服务调用
func (s *UnstableService) Call(requestID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 模拟网络延迟和处理时间
	time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)

	// 根据失败率随机决定成功或失败
	if rand.Float64() < s.failureRate {
		fmt.Printf("服务调用失败: %s\n", requestID)
		return fmt.Errorf("service call failed for request %s", requestID)
	}

	fmt.Printf("服务调用成功: %s\n", requestID)
	return nil
}

// SetFailureRate 动态设置失败率（用于演示）
func (s *UnstableService) SetFailureRate(rate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failureRate = rate
	fmt.Printf("服务失败率调整为: %.2f%%\n", rate*100)
}

func main() {
	fmt.Println("=== 熔断器模式演示 ===")
	fmt.Println("演示熔断器如何保护不稳定的服务调用")

	rand.Seed(time.Now().UnixNano())

	// 创建熔断器配置
	config := CircuitBreakerConfig{
		MaxFailures:     5,               // 最大失败次数
		ResetTimeout:    3 * time.Second, // 3秒后尝试恢复
		FailureRatio:    0.5,             // 50%失败率触发熔断
		MinRequestCount: 10,              // 至少10个请求后才考虑熔断
	}

	// 创建熔断器和不稳定服务
	circuitBreaker := NewCircuitBreaker(config)
	service := NewUnstableService(0.7) // 初始70%失败率

	fmt.Printf("熔断器配置: 失败率阈值=%.0f%%, 最小请求数=%d, 重置超时=%v\n",
		config.FailureRatio*100, config.MinRequestCount, config.ResetTimeout)

	// 模拟客户端并发请求
	var wg sync.WaitGroup
	requestCount := 50

	// 第一阶段：高失败率，触发熔断
	fmt.Println("\n=== 第一阶段：高失败率测试 ===")
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(reqID int) {
			defer wg.Done()

			// 通过熔断器调用服务
			err := circuitBreaker.Call(func() error {
				return service.Call(fmt.Sprintf("req-%d", reqID))
			})

			if err != nil {
				fmt.Printf("请求 req-%d 失败: %v\n", reqID, err)
			}

			// 定期打印统计信息
			requests, failures, state := circuitBreaker.GetStats()
			if reqID%5 == 0 {
				fmt.Printf("当前统计: 请求=%d, 失败=%d, 状态=%s\n",
					requests, failures, state)
			}
		}(i)

		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	time.Sleep(1 * time.Second)

	// 第二阶段：等待熔断器恢复
	fmt.Println("\n=== 第二阶段：等待熔断器恢复 ===")
	service.SetFailureRate(0.2) // 降低失败率到20%

	fmt.Println("等待熔断器重置...")
	time.Sleep(4 * time.Second) // 等待超过重置时间

	// 第三阶段：低失败率，熔断器恢复
	fmt.Println("\n=== 第三阶段：低失败率测试 ===")
	for i := 21; i <= requestCount; i++ {
		wg.Add(1)
		go func(reqID int) {
			defer wg.Done()

			err := circuitBreaker.Call(func() error {
				return service.Call(fmt.Sprintf("req-%d", reqID))
			})

			if err != nil {
				fmt.Printf("请求 req-%d 失败: %v\n", reqID, err)
			}

			requests, failures, state := circuitBreaker.GetStats()
			if reqID%10 == 0 {
				fmt.Printf("当前统计: 请求=%d, 失败=%d, 状态=%s\n",
					requests, failures, state)
			}
		}(i)

		time.Sleep(150 * time.Millisecond)
	}

	wg.Wait()

	// 最终统计
	fmt.Println("\n=== 最终统计 ===")
	requests, failures, state := circuitBreaker.GetStats()
	fmt.Printf("总请求数: %d\n", requests)
	fmt.Printf("总失败数: %d\n", failures)
	fmt.Printf("失败率: %.2f%%\n", float64(failures)/float64(requests)*100)
	fmt.Printf("最终状态: %s\n", state)

	fmt.Println("\n熔断器演示完成！")
	fmt.Println("观察要点：")
	fmt.Println("1. 失败率达到阈值时自动熔断")
	fmt.Println("2. 熔断期间快速失败，保护系统")
	fmt.Println("3. 超时后自动尝试恢复")
	fmt.Println("4. 半开状态的试探机制")
}
