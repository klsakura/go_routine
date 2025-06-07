package main

import (
	"fmt"
	"sync"
	"time"
)

// 信号量实现演示
type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}

// 资源池演示
type Resource struct {
	ID   int
	Name string
}

type ResourcePool struct {
	resources []Resource
	semaphore *Semaphore
	mu        sync.Mutex
	available []Resource
	inUse     map[int]Resource
}

func NewResourcePool(resources []Resource) *ResourcePool {
	pool := &ResourcePool{
		resources: resources,
		semaphore: NewSemaphore(len(resources)),
		available: make([]Resource, len(resources)),
		inUse:     make(map[int]Resource),
	}

	copy(pool.available, resources)
	return pool
}

func (rp *ResourcePool) AcquireResource() (*Resource, error) {
	// 获取信号量
	rp.semaphore.Acquire()

	rp.mu.Lock()
	defer rp.mu.Unlock()

	if len(rp.available) == 0 {
		rp.semaphore.Release()
		return nil, fmt.Errorf("no resources available")
	}

	// 取出一个资源
	resource := rp.available[0]
	rp.available = rp.available[1:]
	rp.inUse[resource.ID] = resource

	return &resource, nil
}

func (rp *ResourcePool) ReleaseResource(resource *Resource) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	if _, exists := rp.inUse[resource.ID]; exists {
		delete(rp.inUse, resource.ID)
		rp.available = append(rp.available, *resource)
		rp.semaphore.Release()
	}
}

func (rp *ResourcePool) GetStats() (int, int, int) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	return len(rp.available), len(rp.inUse), rp.semaphore.Available()
}

// 限制并发连接数的例子
type ConnectionManager struct {
	semaphore *Semaphore
	active    int64
	mu        sync.Mutex
}

func NewConnectionManager(maxConnections int) *ConnectionManager {
	return &ConnectionManager{
		semaphore: NewSemaphore(maxConnections),
	}
}

func (cm *ConnectionManager) HandleConnection(clientID int) error {
	fmt.Printf("客户端 %d 尝试连接...\n", clientID)

	// 尝试获取连接许可
	if !cm.semaphore.TryAcquire() {
		fmt.Printf("客户端 %d 连接被拒绝 (连接数已满)\n", clientID)
		return fmt.Errorf("connection limit reached")
	}

	cm.mu.Lock()
	cm.active++
	currentActive := cm.active
	cm.mu.Unlock()

	fmt.Printf("客户端 %d 连接成功 (当前活跃连接: %d)\n", clientID, currentActive)

	// 模拟连接处理时间
	time.Sleep(time.Duration(clientID%3+1) * time.Second)

	// 释放连接
	cm.mu.Lock()
	cm.active--
	currentActive = cm.active
	cm.mu.Unlock()

	cm.semaphore.Release()
	fmt.Printf("客户端 %d 断开连接 (当前活跃连接: %d)\n", clientID, currentActive)

	return nil
}

func main() {
	fmt.Println("=== 信号量实现演示 ===")

	// 示例1: 资源池管理
	fmt.Println("\n1. 资源池管理演示:")

	resources := []Resource{
		{ID: 1, Name: "数据库连接1"},
		{ID: 2, Name: "数据库连接2"},
		{ID: 3, Name: "数据库连接3"},
	}

	pool := NewResourcePool(resources)
	var wg sync.WaitGroup

	// 启动5个工作者竞争3个资源
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			fmt.Printf("工作者 %d 请求资源\n", workerID)
			resource, err := pool.AcquireResource()
			if err != nil {
				fmt.Printf("工作者 %d 获取资源失败: %v\n", workerID, err)
				return
			}

			fmt.Printf("工作者 %d 获得资源: %s\n", workerID, resource.Name)

			// 模拟使用资源
			time.Sleep(time.Duration(workerID) * 500 * time.Millisecond)

			pool.ReleaseResource(resource)
			fmt.Printf("工作者 %d 释放资源: %s\n", workerID, resource.Name)

			available, inUse, semAvail := pool.GetStats()
			fmt.Printf("资源状态 - 可用: %d, 使用中: %d, 信号量可用: %d\n",
				available, inUse, semAvail)
		}(i)
	}

	wg.Wait()

	// 示例2: 连接数限制
	fmt.Println("\n2. 连接数限制演示:")

	connManager := NewConnectionManager(3) // 最多3个并发连接

	// 启动8个客户端尝试连接
	for i := 1; i <= 8; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			connManager.HandleConnection(clientID)
		}(i)

		time.Sleep(200 * time.Millisecond) // 错开连接时间
	}

	wg.Wait()

	// 示例3: 批量处理限制
	fmt.Println("\n3. 批量处理限制演示:")

	processSemaphore := NewSemaphore(2) // 最多2个并发处理

	tasks := []string{"任务A", "任务B", "任务C", "任务D", "任务E"}

	for _, task := range tasks {
		wg.Add(1)
		go func(taskName string) {
			defer wg.Done()

			fmt.Printf("%s 等待处理许可...\n", taskName)
			processSemaphore.Acquire()

			fmt.Printf("%s 开始处理 (可用许可: %d)\n", taskName, processSemaphore.Available())

			// 模拟处理时间
			time.Sleep(2 * time.Second)

			fmt.Printf("%s 处理完成\n", taskName)
			processSemaphore.Release()
		}(task)

		time.Sleep(300 * time.Millisecond)
	}

	wg.Wait()

	fmt.Println("\n信号量演示完成！")
}
