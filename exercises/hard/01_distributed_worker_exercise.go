/*
Golang并发编程练习 - 困难级别
练习文件：01_distributed_worker_exercise.go
练习主题：分布式工作者系统

练习目标：
1. 实现一致性哈希算法
2. 设计分布式任务调度系统
3. 处理节点动态添加和移除
4. 实现负载均衡和故障转移

练习任务：
- 任务1：实现基础的一致性哈希环
- 任务2：添加虚拟节点提高平衡性
- 任务3：实现任务路由和分发
- 任务4：处理节点故障和恢复

运行方式：go run exercises/hard/01_distributed_worker_exercise.go
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task 任务结构
type Task struct {
	ID       string
	Data     interface{}
	Priority int
	Created  time.Time
}

// WorkerNode 工作节点接口
type WorkerNode interface {
	GetID() string
	ProcessTask(task Task) error
	IsHealthy() bool
	GetLoad() float64
}

// SimpleWorker 简单工作者实现
type SimpleWorker struct {
	id             string
	processedCount int64
	isHealthy      bool
	processingTime time.Duration
	mu             sync.RWMutex
}

// TODO: 实现SimpleWorker的所有方法
func NewSimpleWorker(id string, processingTime time.Duration) *SimpleWorker {
	// 在这里实现您的代码
	return nil
}

func (w *SimpleWorker) GetID() string {
	// 在这里实现您的代码
	return ""
}

func (w *SimpleWorker) ProcessTask(task Task) error {
	// 在这里实现您的代码
	// 提示：
	// 1. 检查健康状态
	// 2. 模拟处理时间
	// 3. 更新处理计数
	// 4. 打印处理信息
	return nil
}

func (w *SimpleWorker) IsHealthy() bool {
	// 在这里实现您的代码
	return false
}

func (w *SimpleWorker) GetLoad() float64 {
	// 在这里实现您的代码
	// 提示：基于处理时间计算负载
	return 0.0
}

func (w *SimpleWorker) SetHealthy(healthy bool) {
	// 在这里实现您的代码
}

// ConsistentHashRing 一致性哈希环
type ConsistentHashRing struct {
	// TODO: 定义需要的字段
	// 提示：需要存储虚拟节点、排序的哈希值、工作者等
}

// TODO: 实现一致性哈希环的所有方法
func NewConsistentHashRing(virtualNodes int) *ConsistentHashRing {
	// 在这里实现您的代码
	return nil
}

func (ring *ConsistentHashRing) hashFunction(key string) uint32 {
	// 在这里实现您的代码
	// 提示：使用crc32.ChecksumIEEE
	return 0
}

func (ring *ConsistentHashRing) AddNode(node WorkerNode) {
	// 在这里实现您的代码
	// 提示：
	// 1. 为每个节点创建多个虚拟节点
	// 2. 将虚拟节点添加到哈希环
	// 3. 保持哈希值的排序
}

func (ring *ConsistentHashRing) RemoveNode(nodeID string) {
	// 在这里实现您的代码
	// 提示：移除节点的所有虚拟节点
}

func (ring *ConsistentHashRing) GetNode(taskID string) (WorkerNode, error) {
	// 在这里实现您的代码
	// 提示：
	// 1. 计算任务的哈希值
	// 2. 在环上找到第一个大于等于该哈希值的节点
	// 3. 如果没找到，选择第一个节点（环形特性）
	return nil, nil
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	// TODO: 定义需要的字段
	// 提示：需要哈希环、任务队列、同步原语等
}

// TODO: 实现任务调度器的所有方法
func NewTaskScheduler(virtualNodes int) *TaskScheduler {
	// 在这里实现您的代码
	return nil
}

func (ts *TaskScheduler) AddWorker(worker WorkerNode) {
	// 在这里实现您的代码
}

func (ts *TaskScheduler) RemoveWorker(workerID string) {
	// 在这里实现您的代码
}

func (ts *TaskScheduler) SubmitTask(task Task) error {
	// 在这里实现您的代码
	// 提示：将任务添加到队列
	return nil
}

func (ts *TaskScheduler) Start() {
	// 在这里实现您的代码
	// 提示：启动后台任务处理协程
}

func (ts *TaskScheduler) Stop() {
	// 在这里实现您的代码
}

func (ts *TaskScheduler) processTask() {
	// 在这里实现您的代码
	// 提示：
	// 1. 从队列获取任务
	// 2. 根据任务ID路由到对应工作者
	// 3. 检查工作者健康状态
	// 4. 异步处理任务
}

func main() {
	fmt.Println("=== 分布式工作者系统练习 ===")
	rand.Seed(time.Now().UnixNano())

	// 任务1：测试一致性哈希环
	fmt.Println("\n任务1：一致性哈希环基础测试")
	// TODO:
	// 1. 创建一致性哈希环（3个虚拟节点）
	// 2. 添加4个工作者节点
	// 3. 测试任务路由：提交10个任务，观察分布情况

	fmt.Println("任务1完成\n")

	// 任务2：负载均衡测试
	fmt.Println("任务2：负载均衡测试")
	// TODO:
	// 1. 创建5个不同处理能力的工作者
	// 2. 提交50个任务
	// 3. 统计每个工作者处理的任务数
	// 4. 分析负载分布是否均匀

	fmt.Println("任务2完成\n")

	// 任务3：动态节点管理
	fmt.Println("任务3：动态节点管理")
	// TODO:
	// 1. 启动系统，提交一批任务
	// 2. 运行期间添加新的工作者节点
	// 3. 模拟某个节点故障（设置为不健康）
	// 4. 移除故障节点
	// 5. 恢复节点并重新添加

	fmt.Println("任务3完成\n")

	// 任务4：高并发压力测试
	fmt.Println("任务4：高并发压力测试")
	// TODO:
	// 1. 创建一个包含10个工作者的系统
	// 2. 并发提交1000个任务
	// 3. 测量处理时间和吞吐量
	// 4. 模拟部分节点故障情况下的性能

	fmt.Println("任务4完成\n")

	// 任务5：一致性验证
	fmt.Println("任务5：一致性验证（挑战任务）")
	// TODO: 高级挑战
	// 验证相同的任务ID总是路由到相同的节点
	// 在节点添加/删除前后测试路由一致性

	fmt.Println("所有练习完成！")

	// 反思问题：
	fmt.Println("\n思考题：")
	fmt.Println("1. 虚拟节点数量如何影响负载均衡？")
	fmt.Println("2. 节点故障时，如何保证任务不丢失？")
	fmt.Println("3. 一致性哈希相比简单哈希有什么优势？")
	fmt.Println("4. 如何设计一个支持权重的负载均衡算法？")
}
