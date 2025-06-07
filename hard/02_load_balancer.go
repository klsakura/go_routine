package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// 负载均衡器演示
type Server struct {
	ID      int
	Address string
	Weight  int
	Active  int64 // 当前活跃连接数
	Total   int64 // 总处理请求数
	Failed  int64 // 失败请求数
	Healthy bool
	mu      sync.RWMutex
}

func NewServer(id int, address string, weight int) *Server {
	return &Server{
		ID:      id,
		Address: address,
		Weight:  weight,
		Healthy: true,
	}
}

func (s *Server) ProcessRequest(requestID string) error {
	atomic.AddInt64(&s.Active, 1)
	defer atomic.AddInt64(&s.Active, -1)

	// 模拟请求处理时间
	processingTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond

	fmt.Printf("服务器 %d (%s) 开始处理请求 %s\n", s.ID, s.Address, requestID)

	time.Sleep(processingTime)

	// 模拟5%的失败率
	if rand.Float32() < 0.05 {
		atomic.AddInt64(&s.Failed, 1)
		fmt.Printf("服务器 %d 处理请求 %s 失败\n", s.ID, requestID)
		return fmt.Errorf("server %d failed to process request", s.ID)
	}

	atomic.AddInt64(&s.Total, 1)
	fmt.Printf("服务器 %d (%s) 完成请求 %s (用时: %v)\n", s.ID, s.Address, requestID, processingTime)

	return nil
}

func (s *Server) GetStats() (int64, int64, int64) {
	return atomic.LoadInt64(&s.Active), atomic.LoadInt64(&s.Total), atomic.LoadInt64(&s.Failed)
}

func (s *Server) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Healthy
}

func (s *Server) SetHealthy(healthy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Healthy = healthy
	if !healthy {
		fmt.Printf("服务器 %d 标记为不健康\n", s.ID)
	} else {
		fmt.Printf("服务器 %d 恢复健康\n", s.ID)
	}
}

// 负载均衡策略接口
type LoadBalanceStrategy interface {
	Select(servers []*Server) *Server
	GetName() string
}

// 轮询策略
type RoundRobinStrategy struct {
	current int64
}

func (rr *RoundRobinStrategy) Select(servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	// 只选择健康的服务器
	healthyServers := make([]*Server, 0)
	for _, server := range servers {
		if server.IsHealthy() {
			healthyServers = append(healthyServers, server)
		}
	}

	if len(healthyServers) == 0 {
		return nil
	}

	index := atomic.AddInt64(&rr.current, 1) % int64(len(healthyServers))
	return healthyServers[index]
}

func (rr *RoundRobinStrategy) GetName() string {
	return "RoundRobin"
}

// 最少连接策略
type LeastConnectionsStrategy struct{}

func (lc *LeastConnectionsStrategy) Select(servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	var selected *Server
	minConnections := int64(-1)

	for _, server := range servers {
		if !server.IsHealthy() {
			continue
		}

		active, _, _ := server.GetStats()
		if minConnections == -1 || active < minConnections {
			minConnections = active
			selected = server
		}
	}

	return selected
}

func (lc *LeastConnectionsStrategy) GetName() string {
	return "LeastConnections"
}

// 加权轮询策略
type WeightedRoundRobinStrategy struct {
	weights map[int]int
	current map[int]int
	mu      sync.Mutex
}

func NewWeightedRoundRobinStrategy() *WeightedRoundRobinStrategy {
	return &WeightedRoundRobinStrategy{
		weights: make(map[int]int),
		current: make(map[int]int),
	}
}

func (wrr *WeightedRoundRobinStrategy) Select(servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	// 初始化权重
	for _, server := range servers {
		if _, exists := wrr.weights[server.ID]; !exists {
			wrr.weights[server.ID] = server.Weight
			wrr.current[server.ID] = 0
		}
	}

	var selected *Server
	maxWeight := -1
	totalWeight := 0

	for _, server := range servers {
		if !server.IsHealthy() {
			continue
		}

		wrr.current[server.ID] += wrr.weights[server.ID]
		totalWeight += wrr.weights[server.ID]

		if wrr.current[server.ID] > maxWeight {
			maxWeight = wrr.current[server.ID]
			selected = server
		}
	}

	if selected != nil {
		wrr.current[selected.ID] -= totalWeight
	}

	return selected
}

func (wrr *WeightedRoundRobinStrategy) GetName() string {
	return "WeightedRoundRobin"
}

// 负载均衡器
type LoadBalancer struct {
	servers  []*Server
	strategy LoadBalanceStrategy
	stats    struct {
		totalRequests  int64
		failedRequests int64
	}
	mu sync.RWMutex
}

func NewLoadBalancer(strategy LoadBalanceStrategy) *LoadBalancer {
	return &LoadBalancer{
		servers:  make([]*Server, 0),
		strategy: strategy,
	}
}

func (lb *LoadBalancer) AddServer(server *Server) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.servers = append(lb.servers, server)
	fmt.Printf("添加服务器: %d (%s) 权重=%d\n", server.ID, server.Address, server.Weight)
}

func (lb *LoadBalancer) RemoveServer(serverID int) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i, server := range lb.servers {
		if server.ID == serverID {
			lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)
			fmt.Printf("移除服务器: %d\n", serverID)
			break
		}
	}
}

func (lb *LoadBalancer) ProcessRequest(requestID string) error {
	atomic.AddInt64(&lb.stats.totalRequests, 1)

	lb.mu.RLock()
	servers := make([]*Server, len(lb.servers))
	copy(servers, lb.servers)
	lb.mu.RUnlock()

	server := lb.strategy.Select(servers)
	if server == nil {
		atomic.AddInt64(&lb.stats.failedRequests, 1)
		return fmt.Errorf("no healthy server available")
	}

	err := server.ProcessRequest(requestID)
	if err != nil {
		atomic.AddInt64(&lb.stats.failedRequests, 1)
		return err
	}

	return nil
}

func (lb *LoadBalancer) PrintStats() {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	fmt.Printf("\n=== 负载均衡器统计 (策略: %s) ===\n", lb.strategy.GetName())
	fmt.Printf("总请求数: %d\n", atomic.LoadInt64(&lb.stats.totalRequests))
	fmt.Printf("失败请求数: %d\n", atomic.LoadInt64(&lb.stats.failedRequests))

	fmt.Println("\n服务器统计:")
	for _, server := range lb.servers {
		active, total, failed := server.GetStats()
		health := "健康"
		if !server.IsHealthy() {
			health = "不健康"
		}
		fmt.Printf("服务器 %d: 活跃=%d, 总计=%d, 失败=%d, 状态=%s\n",
			server.ID, active, total, failed, health)
	}
}

// 健康检查器
func (lb *LoadBalancer) StartHealthCheck() {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			lb.mu.RLock()
			servers := make([]*Server, len(lb.servers))
			copy(servers, lb.servers)
			lb.mu.RUnlock()

			for _, server := range servers {
				// 模拟健康检查：10%概率变为不健康，20%概率恢复
				if server.IsHealthy() {
					if rand.Float32() < 0.1 {
						server.SetHealthy(false)
					}
				} else {
					if rand.Float32() < 0.2 {
						server.SetHealthy(true)
					}
				}
			}
		}
	}()
}

func main() {
	fmt.Println("=== 负载均衡器演示 ===")

	rand.Seed(time.Now().UnixNano())

	// 创建不同策略的负载均衡器
	strategies := []LoadBalanceStrategy{
		&RoundRobinStrategy{},
		&LeastConnectionsStrategy{},
		NewWeightedRoundRobinStrategy(),
	}

	for _, strategy := range strategies {
		fmt.Printf("\n--- 测试策略: %s ---\n", strategy.GetName())

		lb := NewLoadBalancer(strategy)

		// 添加服务器
		lb.AddServer(NewServer(1, "192.168.1.1:8080", 3))
		lb.AddServer(NewServer(2, "192.168.1.2:8080", 2))
		lb.AddServer(NewServer(3, "192.168.1.3:8080", 1))
		lb.AddServer(NewServer(4, "192.168.1.4:8080", 4))

		// 启动健康检查
		lb.StartHealthCheck()

		// 模拟并发请求
		var wg sync.WaitGroup
		numRequests := 20

		for i := 1; i <= numRequests; i++ {
			wg.Add(1)
			go func(reqID int) {
				defer wg.Done()
				err := lb.ProcessRequest(fmt.Sprintf("req-%d", reqID))
				if err != nil {
					fmt.Printf("请求 req-%d 失败: %v\n", reqID, err)
				}
			}(i)

			// 错开请求时间
			time.Sleep(50 * time.Millisecond)
		}

		wg.Wait()
		time.Sleep(1 * time.Second) // 等待请求完成

		lb.PrintStats()
		time.Sleep(2 * time.Second) // 间隔时间
	}

	fmt.Println("\n负载均衡演示完成！")
}
