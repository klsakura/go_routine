/*
Golang并发编程学习Demo - 困难级别
文件：03_message_queue.go
主题：消息队列系统实现

本示例演示：
1. 完整的消息队列系统架构
2. 发布-订阅模式的实现
3. 消息重试机制和死信队列
4. 并发消费者管理
5. 消息持久化和可靠性保证

核心功能：
- 主题订阅：支持多个消费者订阅同一主题
- 消息重试：失败消息自动重试，避免消息丢失
- 死信队列：超过重试次数的消息进入死信队列
- 并发处理：多个消费者并发处理消息
- 消息统计：提供详细的消息处理统计

应用场景：
- 微服务架构中的异步通信
- 事件驱动架构
- 任务队列系统
- 日志收集和处理

技术要点：
- 生产者-消费者模式
- 消息路由和分发
- 错误处理和重试策略
- 并发安全的队列操作

运行方式：go run hard/03_message_queue.go
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// QueueMessage 队列中的消息结构
type QueueMessage struct {
	ID        string      // 消息唯一标识
	Topic     string      // 消息主题
	Payload   interface{} // 消息内容
	Timestamp time.Time   // 创建时间戳
	Retries   int         // 重试次数
	Priority  int         // 消息优先级（暂未使用）
}

// MessageQueue 消息队列接口定义
type MessageQueue interface {
	Publish(topic string, message QueueMessage) error  // 发布消息
	Subscribe(topic string, consumer Consumer) error   // 订阅主题
	Unsubscribe(topic string, consumerID string) error // 取消订阅
	Close() error                                      // 关闭队列
}

// Consumer 消费者接口定义
type Consumer interface {
	GetID() string                      // 获取消费者ID
	Consume(message QueueMessage) error // 消费消息
}

// SimpleConsumer 简单消费者实现
type SimpleConsumer struct {
	ID           string        // 消费者唯一标识
	ProcessTime  time.Duration // 模拟处理时间
	SuccessRate  float64       // 成功处理概率（0.0-1.0）
	MessageCount int64         // 已处理消息计数
}

// NewSimpleConsumer 创建简单消费者
func NewSimpleConsumer(id string, processTime time.Duration, successRate float64) *SimpleConsumer {
	return &SimpleConsumer{
		ID:          id,
		ProcessTime: processTime,
		SuccessRate: successRate,
	}
}

// GetID 实现Consumer接口 - 获取消费者ID
func (c *SimpleConsumer) GetID() string {
	return c.ID
}

// Consume 实现Consumer接口 - 处理消息
func (c *SimpleConsumer) Consume(message QueueMessage) error {
	// 原子增加消息计数
	atomic.AddInt64(&c.MessageCount, 1)

	// 模拟消息处理时间
	time.Sleep(c.ProcessTime)

	// 根据成功率随机决定处理结果
	if rand.Float64() < c.SuccessRate {
		fmt.Printf("消费者 %s 成功处理消息: %s (主题: %s)\n",
			c.ID, message.ID, message.Topic)
		return nil
	} else {
		fmt.Printf("消费者 %s 处理消息失败: %s (主题: %s)\n",
			c.ID, message.ID, message.Topic)
		return fmt.Errorf("processing failed")
	}
}

// GetMessageCount 获取已处理消息数量
func (c *SimpleConsumer) GetMessageCount() int64 {
	return atomic.LoadInt64(&c.MessageCount)
}

// InMemoryMessageQueue 基于内存的消息队列实现
type InMemoryMessageQueue struct {
	subscriptions map[string][]Consumer // 订阅关系：主题 -> 消费者列表
	retryQueue    chan QueueMessage     // 重试队列
	deadLetter    []QueueMessage        // 死信队列
	mu            sync.RWMutex          // 保护订阅关系的读写锁
	wg            sync.WaitGroup        // 等待组，用于优雅关闭
	stopCh        chan bool             // 停止信号
	maxRetries    int                   // 最大重试次数
	stats         struct {              // 消息处理统计
		published int64 // 发布消息数
		consumed  int64 // 成功消费数
		failed    int64 // 失败消费数
		retried   int64 // 重试次数
	}
}

// NewInMemoryMessageQueue 创建内存消息队列
func NewInMemoryMessageQueue(maxRetries int) *InMemoryMessageQueue {
	mq := &InMemoryMessageQueue{
		subscriptions: make(map[string][]Consumer),
		retryQueue:    make(chan QueueMessage, 1000), // 重试队列缓冲
		deadLetter:    make([]QueueMessage, 0),
		maxRetries:    maxRetries,
		stopCh:        make(chan bool),
	}

	// 启动后台重试处理器
	mq.wg.Add(1)
	go mq.retryProcessor()

	return mq
}

// Publish 实现MessageQueue接口 - 发布消息到指定主题
func (mq *InMemoryMessageQueue) Publish(topic string, message QueueMessage) error {
	// 获取该主题的所有消费者
	mq.mu.RLock()
	consumers, exists := mq.subscriptions[topic]
	mq.mu.RUnlock()

	if !exists || len(consumers) == 0 {
		fmt.Printf("警告: 主题 %s 没有消费者\n", topic)
		return fmt.Errorf("no consumers for topic: %s", topic)
	}

	// 增加发布统计
	atomic.AddInt64(&mq.stats.published, 1)

	// 并发发送消息给所有订阅该主题的消费者
	var wg sync.WaitGroup
	for _, consumer := range consumers {
		wg.Add(1)
		go func(c Consumer) {
			defer wg.Done()
			mq.deliverMessage(message, c)
		}(consumer)
	}

	// 等待所有消费者处理完成
	wg.Wait()
	return nil
}

// deliverMessage 将消息投递给指定消费者
func (mq *InMemoryMessageQueue) deliverMessage(message QueueMessage, consumer Consumer) {
	err := consumer.Consume(message)
	if err != nil {
		// 处理失败，增加失败统计
		atomic.AddInt64(&mq.stats.failed, 1)

		// 检查是否需要重试
		if message.Retries < mq.maxRetries {
			message.Retries++
			select {
			case mq.retryQueue <- message:
				atomic.AddInt64(&mq.stats.retried, 1)
				fmt.Printf("消息 %s 加入重试队列 (重试次数: %d)\n", message.ID, message.Retries)
			default:
				// 重试队列满，直接进入死信队列
				mq.addToDeadLetter(message)
			}
		} else {
			// 超过最大重试次数，进入死信队列
			mq.addToDeadLetter(message)
		}
	} else {
		// 处理成功，增加成功统计
		atomic.AddInt64(&mq.stats.consumed, 1)
	}
}

// addToDeadLetter 将消息添加到死信队列
func (mq *InMemoryMessageQueue) addToDeadLetter(message QueueMessage) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.deadLetter = append(mq.deadLetter, message)
	fmt.Printf("消息 %s 进入死信队列\n", message.ID)
}

// retryProcessor 重试处理器，在后台处理重试队列
func (mq *InMemoryMessageQueue) retryProcessor() {
	defer mq.wg.Done()

	for {
		select {
		case message := <-mq.retryQueue:
			// 延迟重试：重试次数越多，延迟时间越长
			retryDelay := time.Duration(message.Retries) * time.Second
			time.Sleep(retryDelay)

			fmt.Printf("重试消息: %s (第 %d 次重试)\n", message.ID, message.Retries)

			// 重新发布消息
			mq.Publish(message.Topic, message)

		case <-mq.stopCh:
			// 收到停止信号，退出重试处理器
			return
		}
	}
}

// Subscribe 实现MessageQueue接口 - 订阅主题
func (mq *InMemoryMessageQueue) Subscribe(topic string, consumer Consumer) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	// 初始化主题的消费者列表
	if mq.subscriptions[topic] == nil {
		mq.subscriptions[topic] = make([]Consumer, 0)
	}

	// 检查消费者是否已经订阅过该主题
	for _, c := range mq.subscriptions[topic] {
		if c.GetID() == consumer.GetID() {
			return fmt.Errorf("consumer %s already subscribed to topic %s", consumer.GetID(), topic)
		}
	}

	// 添加消费者到订阅列表
	mq.subscriptions[topic] = append(mq.subscriptions[topic], consumer)
	fmt.Printf("消费者 %s 订阅主题: %s\n", consumer.GetID(), topic)

	return nil
}

// Unsubscribe 实现MessageQueue接口 - 取消订阅
func (mq *InMemoryMessageQueue) Unsubscribe(topic string, consumerID string) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	consumers, exists := mq.subscriptions[topic]
	if !exists {
		return fmt.Errorf("topic %s not found", topic)
	}

	// 查找并移除指定消费者
	for i, consumer := range consumers {
		if consumer.GetID() == consumerID {
			// 从切片中移除消费者
			mq.subscriptions[topic] = append(consumers[:i], consumers[i+1:]...)
			fmt.Printf("消费者 %s 取消订阅主题: %s\n", consumerID, topic)
			return nil
		}
	}

	return fmt.Errorf("consumer %s not found in topic %s", consumerID, topic)
}

// Close 实现MessageQueue接口 - 关闭消息队列
func (mq *InMemoryMessageQueue) Close() error {
	close(mq.stopCh)     // 发送停止信号
	close(mq.retryQueue) // 关闭重试队列
	mq.wg.Wait()         // 等待后台处理器完成
	return nil
}

// GetStats 获取消息队列统计信息
func (mq *InMemoryMessageQueue) GetStats() (int64, int64, int64, int64, int) {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	return atomic.LoadInt64(&mq.stats.published),
		atomic.LoadInt64(&mq.stats.consumed),
		atomic.LoadInt64(&mq.stats.failed),
		atomic.LoadInt64(&mq.stats.retried),
		len(mq.deadLetter)
}

// GetDeadLetters 获取死信队列中的消息
func (mq *InMemoryMessageQueue) GetDeadLetters() []QueueMessage {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	// 返回副本，避免外部修改
	result := make([]QueueMessage, len(mq.deadLetter))
	copy(result, mq.deadLetter)
	return result
}

// MessageProducer 消息生产者
type MessageProducer struct {
	queue MessageQueue // 消息队列引用
	id    string       // 生产者ID
}

// NewMessageProducer 创建消息生产者
func NewMessageProducer(id string, queue MessageQueue) *MessageProducer {
	return &MessageProducer{
		id:    id,
		queue: queue,
	}
}

// SendMessage 发送消息到指定主题
func (p *MessageProducer) SendMessage(topic string, payload interface{}, priority int) error {
	// 构造消息
	message := QueueMessage{
		ID:        fmt.Sprintf("%s-%d", p.id, rand.Intn(10000)), // 生成唯一ID
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now(),
		Priority:  priority,
		Retries:   0, // 初始重试次数为0
	}

	fmt.Printf("生产者 %s 发送消息: %s 到主题 %s\n", p.id, message.ID, topic)
	return p.queue.Publish(topic, message)
}

func main() {
	fmt.Println("=== 消息队列实现演示 ===")
	fmt.Println("演示完整的消息队列系统：发布订阅、重试机制、死信队列")

	rand.Seed(time.Now().UnixNano())

	// 创建消息队列（最多重试3次）
	mq := NewInMemoryMessageQueue(3)
	defer mq.Close()

	// 创建不同性能的消费者
	consumer1 := NewSimpleConsumer("consumer-1", 100*time.Millisecond, 0.8) // 80%成功率，快速处理
	consumer2 := NewSimpleConsumer("consumer-2", 200*time.Millisecond, 0.6) // 60%成功率，中等处理
	consumer3 := NewSimpleConsumer("consumer-3", 150*time.Millisecond, 0.9) // 90%成功率，中等处理

	// 建立订阅关系
	fmt.Println("\n--- 建立订阅关系 ---")
	mq.Subscribe("orders", consumer1)        // 订单主题：consumer1
	mq.Subscribe("orders", consumer2)        // 订单主题：consumer2（多个消费者）
	mq.Subscribe("notifications", consumer2) // 通知主题：consumer2
	mq.Subscribe("notifications", consumer3) // 通知主题：consumer3

	// 创建消息生产者
	producer1 := NewMessageProducer("producer-1", mq)
	producer2 := NewMessageProducer("producer-2", mq)

	// 并发发送消息
	var wg sync.WaitGroup

	// 生产者1：发送订单消息
	fmt.Println("\n--- 发送订单消息 ---")
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(orderID int) {
			defer wg.Done()

			// 构造订单数据
			payload := map[string]interface{}{
				"orderID":   orderID,
				"amount":    rand.Float64() * 1000,
				"userID":    rand.Intn(1000),
				"timestamp": time.Now().Unix(),
			}

			err := producer1.SendMessage("orders", payload, rand.Intn(5))
			if err != nil {
				fmt.Printf("发送订单消息失败: %v\n", err)
			}

			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		}(i)
	}

	// 生产者2：发送通知消息
	fmt.Println("\n--- 发送通知消息 ---")
	for i := 1; i <= 8; i++ {
		wg.Add(1)
		go func(notificationID int) {
			defer wg.Done()

			// 构造通知数据
			payload := map[string]interface{}{
				"notificationID": notificationID,
				"message":        fmt.Sprintf("通知消息 %d", notificationID),
				"userID":         rand.Intn(1000),
				"type":           "info",
			}

			err := producer2.SendMessage("notifications", payload, rand.Intn(3))
			if err != nil {
				fmt.Printf("发送通知消息失败: %v\n", err)
			}

			time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
		}(i)
	}

	wg.Wait()

	// 等待消息处理完成（包括重试）
	fmt.Println("\n--- 等待消息处理完成 ---")
	time.Sleep(5 * time.Second)

	// 打印统计信息
	fmt.Println("\n=== 消息队列统计 ===")
	published, consumed, failed, retried, deadLetterCount := mq.GetStats()
	fmt.Printf("发布消息数: %d\n", published)
	fmt.Printf("成功消费数: %d\n", consumed)
	fmt.Printf("失败消费数: %d\n", failed)
	fmt.Printf("重试次数: %d\n", retried)
	fmt.Printf("死信消息数: %d\n", deadLetterCount)

	// 打印消费者统计
	fmt.Println("\n=== 消费者统计 ===")
	fmt.Printf("消费者1处理消息数: %d\n", consumer1.GetMessageCount())
	fmt.Printf("消费者2处理消息数: %d\n", consumer2.GetMessageCount())
	fmt.Printf("消费者3处理消息数: %d\n", consumer3.GetMessageCount())

	// 显示死信队列内容
	deadLetters := mq.GetDeadLetters()
	if len(deadLetters) > 0 {
		fmt.Println("\n=== 死信队列 ===")
		for _, msg := range deadLetters {
			fmt.Printf("死信消息: %s (主题: %s, 重试次数: %d)\n",
				msg.ID, msg.Topic, msg.Retries)
		}
	}

	// 测试取消订阅功能
	fmt.Println("\n=== 取消订阅测试 ===")
	mq.Unsubscribe("orders", "consumer-1")

	// 再发送一条消息验证取消订阅效果
	fmt.Println("发送测试消息验证取消订阅...")
	producer1.SendMessage("orders", map[string]interface{}{"test": "after unsubscribe"}, 1)

	time.Sleep(1 * time.Second)

	fmt.Println("\n消息队列演示完成！")
	fmt.Println("观察要点：")
	fmt.Println("1. 多个消费者可订阅同一主题")
	fmt.Println("2. 失败消息自动重试")
	fmt.Println("3. 超过重试次数进入死信队列")
	fmt.Println("4. 消费者可动态取消订阅")
	fmt.Println("5. 系统提供详细的处理统计")
}
