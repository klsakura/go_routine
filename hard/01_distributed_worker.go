/*
Golang并发编程学习Demo - 困难级别
文件：01_distributed_worker.go
主题：分布式工作者系统

本示例演示：
1. 分布式任务分发和处理
2. 一致性哈希算法实现
3. 工作者节点的动态管理
4. 任务路由和负载均衡
5. 节点故障处理和恢复

核心技术：
- 一致性哈希：解决分布式系统中的数据分布问题
- 虚拟节点：提高哈希环的平衡性
- 任务路由：根据任务ID路由到对应工作者
- 并发处理：多个工作者并发处理任务

应用场景：
- 分布式计算系统
- 微服务任务调度
- 缓存系统的数据分片
- 负载均衡系统

学习要点：
- 理解一致性哈希的工作原理
- 掌握分布式系统的基本概念
- 学习节点管理和故障处理
- 了解虚拟节点的作用

运行方式：go run hard/01_distributed_worker.go
*/

package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
	"time"
)

// Task 表示一个需要处理的任务
type Task struct {
	ID      string      // 任务唯一标识
	Payload interface{} // 任务数据
	Created time.Time   // 创建时间
}

// Worker 工作者接口定义
type Worker interface {
	GetID() string               // 获取工作者ID
	ProcessTask(task Task) error // 处理任务
	IsHealthy() bool             // 健康检查
	GetProcessedCount() int64    // 获取已处理任务数
	GetLoad() float64            // 获取当前负载
}

// DistributedWorker 分布式工作者具体实现
type DistributedWorker struct {
	id             string        // 工作者唯一标识
	processedCount int64         // 已处理任务计数
	isHealthy      bool          // 健康状态
	processingTime time.Duration // 模拟处理时间
	mu             sync.RWMutex  // 保护并发访问
}

// NewDistributedWorker 创建新的分布式工作者
func NewDistributedWorker(id string, processingTime time.Duration) *DistributedWorker {
	return &DistributedWorker{
		id:             id,
		isHealthy:      true,
		processingTime: processingTime,
	}
}

// GetID 实现Worker接口 - 获取工作者ID
func (w *DistributedWorker) GetID() string {
	return w.id
}

// ProcessTask 实现Worker接口 - 处理任务
func (w *DistributedWorker) ProcessTask(task Task) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 检查健康状态
	if !w.isHealthy {
		return fmt.Errorf("worker %s is not healthy", w.id)
	}

	fmt.Printf("工作者 %s 开始处理任务 %s\n", w.id, task.ID)

	// 模拟任务处理时间
	time.Sleep(w.processingTime)

	w.processedCount++
	fmt.Printf("工作者 %s 完成处理任务 %s (总计: %d)\n", w.id, task.ID, w.processedCount)

	return nil
}

// IsHealthy 实现Worker接口 - 健康检查
func (w *DistributedWorker) IsHealthy() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.isHealthy
}

// GetProcessedCount 实现Worker接口 - 获取已处理任务数
func (w *DistributedWorker) GetProcessedCount() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.processedCount
}

// GetLoad 实现Worker接口 - 获取当前负载（简化实现）
func (w *DistributedWorker) GetLoad() float64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	// 简化的负载计算：基于处理时间
	return float64(w.processingTime) / float64(time.Second)
}

// SetHealthy 设置健康状态（用于模拟故障）
func (w *DistributedWorker) SetHealthy(healthy bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isHealthy = healthy
	status := "健康"
	if !healthy {
		status = "故障"
	}
	fmt.Printf("工作者 %s 状态变更为: %s\n", w.id, status)
}

// ConsistentHash 一致性哈希环实现
type ConsistentHash struct {
	replicas   int               // 虚拟节点数量
	ring       map[uint32]string // 哈希环：哈希值 -> 节点ID
	sortedKeys []uint32          // 排序的哈希值列表
	workers    map[string]Worker // 工作者映射：节点ID -> Worker
	mu         sync.RWMutex      // 保护并发访问
}

// NewConsistentHash 创建一致性哈希环
func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		replicas: replicas,
		ring:     make(map[uint32]string),
		workers:  make(map[string]Worker),
	}
}

// hashFunction 哈希函数：将字符串映射为uint32
func (ch *ConsistentHash) hashFunction(data string) uint32 {
	return crc32.ChecksumIEEE([]byte(data))
}

// AddWorker 添加工作者到哈希环
func (ch *ConsistentHash) AddWorker(worker Worker) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	workerID := worker.GetID()
	ch.workers[workerID] = worker

	// 为每个工作者创建多个虚拟节点
	for i := 0; i < ch.replicas; i++ {
		// 创建虚拟节点ID
		virtualNode := fmt.Sprintf("%s#%d", workerID, i)
		hash := ch.hashFunction(virtualNode)

		// 添加到哈希环
		ch.ring[hash] = workerID
		ch.sortedKeys = append(ch.sortedKeys, hash)
	}

	// 对哈希值进行排序，维护环的有序性
	sort.Slice(ch.sortedKeys, func(i, j int) bool {
		return ch.sortedKeys[i] < ch.sortedKeys[j]
	})

	fmt.Printf("工作者 %s 已添加到哈希环 (虚拟节点数: %d)\n", workerID, ch.replicas)
}

// RemoveWorker 从哈希环移除工作者
func (ch *ConsistentHash) RemoveWorker(workerID string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// 移除所有虚拟节点
	for i := 0; i < ch.replicas; i++ {
		virtualNode := fmt.Sprintf("%s#%d", workerID, i)
		hash := ch.hashFunction(virtualNode)
		delete(ch.ring, hash)

		// 从排序列表中移除
		for j, key := range ch.sortedKeys {
			if key == hash {
				ch.sortedKeys = append(ch.sortedKeys[:j], ch.sortedKeys[j+1:]...)
				break
			}
		}
	}

	// 移除工作者
	delete(ch.workers, workerID)
	fmt.Printf("工作者 %s 已从哈希环移除\n", workerID)
}

// GetWorker 根据任务ID获取对应的工作者
func (ch *ConsistentHash) GetWorker(taskID string) (Worker, error) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if len(ch.ring) == 0 {
		return nil, fmt.Errorf("no workers available")
	}

	// 计算任务的哈希值
	hash := ch.hashFunction(taskID)

	// 在哈希环上找到第一个大于等于该哈希值的节点
	// 如果找不到，则选择第一个节点（环形特性）
	idx := sort.Search(len(ch.sortedKeys), func(i int) bool {
		return ch.sortedKeys[i] >= hash
	})

	// 如果超出范围，回到环的开始
	if idx == len(ch.sortedKeys) {
		idx = 0
	}

	// 获取对应的工作者ID和工作者实例
	workerID := ch.ring[ch.sortedKeys[idx]]
	worker := ch.workers[workerID]

	return worker, nil
}

// GetAllWorkers 获取所有工作者
func (ch *ConsistentHash) GetAllWorkers() []Worker {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	workers := make([]Worker, 0, len(ch.workers))
	for _, worker := range ch.workers {
		workers = append(workers, worker)
	}
	return workers
}

// WorkerManager 工作者管理器
type WorkerManager struct {
	hash     *ConsistentHash // 一致性哈希环
	taskChan chan Task       // 任务队列
	wg       sync.WaitGroup  // 等待组
	stopChan chan bool       // 停止信号
}

// NewWorkerManager 创建工作者管理器
func NewWorkerManager(replicas int) *WorkerManager {
	return &WorkerManager{
		hash:     NewConsistentHash(replicas),
		taskChan: make(chan Task, 100), // 缓冲队列
		stopChan: make(chan bool),
	}
}

// AddWorker 添加工作者
func (wm *WorkerManager) AddWorker(worker Worker) {
	wm.hash.AddWorker(worker)
}

// RemoveWorker 移除工作者
func (wm *WorkerManager) RemoveWorker(workerID string) {
	wm.hash.RemoveWorker(workerID)
}

// SubmitTask 提交任务
func (wm *WorkerManager) SubmitTask(task Task) error {
	select {
	case wm.taskChan <- task:
		fmt.Printf("任务 %s 已提交到队列\n", task.ID)
		return nil
	default:
		return fmt.Errorf("task queue is full")
	}
}

// Start 启动任务处理器
func (wm *WorkerManager) Start() {
	wm.wg.Add(1)
	go wm.taskProcessor()
	fmt.Println("工作者管理器已启动")
}

// Stop 停止任务处理器
func (wm *WorkerManager) Stop() {
	close(wm.stopChan)
	wm.wg.Wait()
	fmt.Println("工作者管理器已停止")
}

// taskProcessor 任务处理器：从队列取任务并路由到工作者
func (wm *WorkerManager) taskProcessor() {
	defer wm.wg.Done()

	for {
		select {
		case task := <-wm.taskChan:
			// 根据任务ID路由到对应工作者
			worker, err := wm.hash.GetWorker(task.ID)
			if err != nil {
				fmt.Printf("获取工作者失败: %v\n", err)
				continue
			}

			// 检查工作者健康状态
			if !worker.IsHealthy() {
				fmt.Printf("工作者 %s 不健康，任务 %s 处理失败\n", worker.GetID(), task.ID)
				continue
			}

			// 异步处理任务
			go func(w Worker, t Task) {
				err := w.ProcessTask(t)
				if err != nil {
					fmt.Printf("任务处理失败: %v\n", err)
				}
			}(worker, task)

		case <-wm.stopChan:
			return
		}
	}
}

// GetStats 获取统计信息
func (wm *WorkerManager) GetStats() map[string]interface{} {
	workers := wm.hash.GetAllWorkers()
	stats := make(map[string]interface{})

	totalProcessed := int64(0)
	healthyCount := 0

	for _, worker := range workers {
		workerStats := map[string]interface{}{
			"processed": worker.GetProcessedCount(),
			"healthy":   worker.IsHealthy(),
			"load":      worker.GetLoad(),
		}
		stats[worker.GetID()] = workerStats

		totalProcessed += worker.GetProcessedCount()
		if worker.IsHealthy() {
			healthyCount++
		}
	}

	stats["total"] = map[string]interface{}{
		"workers":         len(workers),
		"healthy_workers": healthyCount,
		"total_processed": totalProcessed,
	}

	return stats
}

func main() {
	fmt.Println("=== 分布式工作者系统演示 ===")
	fmt.Println("演示一致性哈希在分布式任务调度中的应用")

	// 创建工作者管理器（每个工作者3个虚拟节点）
	manager := NewWorkerManager(3)

	// 创建多个工作者，模拟不同的处理能力
	workers := []*DistributedWorker{
		NewDistributedWorker("worker-1", 100*time.Millisecond), // 快速工作者
		NewDistributedWorker("worker-2", 200*time.Millisecond), // 中等工作者
		NewDistributedWorker("worker-3", 150*time.Millisecond), // 中等工作者
		NewDistributedWorker("worker-4", 300*time.Millisecond), // 慢速工作者
	}

	// 添加工作者到管理器
	fmt.Println("\n--- 初始化工作者 ---")
	for _, worker := range workers {
		manager.AddWorker(worker)
	}

	// 启动管理器
	manager.Start()
	defer manager.Stop()

	// 提交一批任务
	fmt.Println("\n--- 提交任务 ---")
	for i := 1; i <= 20; i++ {
		task := Task{
			ID:      fmt.Sprintf("task-%02d", i),
			Payload: fmt.Sprintf("任务数据 %d", i),
			Created: time.Now(),
		}

		err := manager.SubmitTask(task)
		if err != nil {
			fmt.Printf("提交任务失败: %v\n", err)
		}

		time.Sleep(50 * time.Millisecond) // 控制提交速度
	}

	// 等待一段时间让任务处理
	time.Sleep(3 * time.Second)

	// 模拟工作者故障
	fmt.Println("\n--- 模拟工作者故障 ---")
	workers[1].SetHealthy(false) // worker-2故障

	// 继续提交任务
	fmt.Println("\n--- 故障后继续提交任务 ---")
	for i := 21; i <= 30; i++ {
		task := Task{
			ID:      fmt.Sprintf("task-%02d", i),
			Payload: fmt.Sprintf("任务数据 %d", i),
			Created: time.Now(),
		}

		manager.SubmitTask(task)
		time.Sleep(50 * time.Millisecond)
	}

	// 等待处理完成
	time.Sleep(2 * time.Second)

	// 移除故障工作者
	fmt.Println("\n--- 移除故障工作者 ---")
	manager.RemoveWorker("worker-2")

	// 恢复工作者
	fmt.Println("\n--- 恢复工作者 ---")
	workers[1].SetHealthy(true)
	manager.AddWorker(workers[1])

	// 最后一批任务
	fmt.Println("\n--- 最后一批任务 ---")
	for i := 31; i <= 40; i++ {
		task := Task{
			ID:      fmt.Sprintf("task-%02d", i),
			Payload: fmt.Sprintf("任务数据 %d", i),
			Created: time.Now(),
		}

		manager.SubmitTask(task)
		time.Sleep(50 * time.Millisecond)
	}

	// 等待所有任务完成
	time.Sleep(3 * time.Second)

	// 打印最终统计
	fmt.Println("\n=== 最终统计 ===")
	stats := manager.GetStats()
	for workerID, workerStats := range stats {
		if workerID != "total" {
			fmt.Printf("工作者 %s: %+v\n", workerID, workerStats)
		}
	}
	fmt.Printf("总计: %+v\n", stats["total"])

	fmt.Println("\n分布式工作者演示完成！")
	fmt.Println("观察要点：")
	fmt.Println("1. 一致性哈希确保任务均匀分布")
	fmt.Println("2. 虚拟节点提高负载均衡效果")
	fmt.Println("3. 工作者故障时的处理策略")
	fmt.Println("4. 动态添加/移除工作者的影响")
}
