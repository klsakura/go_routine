package main

import (
	"fmt"
	"sync"
	"time"
)

// 发布订阅模式演示
type Message struct {
	Topic   string
	Content string
	Time    time.Time
}

type Subscriber struct {
	ID       int
	Messages chan Message
	Topics   map[string]bool
	mu       sync.RWMutex
}

func NewSubscriber(id int) *Subscriber {
	return &Subscriber{
		ID:       id,
		Messages: make(chan Message, 10),
		Topics:   make(map[string]bool),
	}
}

func (s *Subscriber) Subscribe(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Topics[topic] = true
	fmt.Printf("订阅者 %d 订阅了主题: %s\n", s.ID, topic)
}

func (s *Subscriber) IsSubscribed(topic string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Topics[topic]
}

func (s *Subscriber) Listen() {
	for msg := range s.Messages {
		fmt.Printf("订阅者 %d 收到消息 [%s]: %s (时间: %s)\n",
			s.ID, msg.Topic, msg.Content, msg.Time.Format("15:04:05"))
		time.Sleep(100 * time.Millisecond) // 模拟处理时间
	}
	fmt.Printf("订阅者 %d 停止监听\n", s.ID)
}

type Publisher struct {
	subscribers []*Subscriber
	mu          sync.RWMutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make([]*Subscriber, 0),
	}
}

func (p *Publisher) AddSubscriber(sub *Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers = append(p.subscribers, sub)
	fmt.Printf("添加订阅者 %d\n", sub.ID)
}

func (p *Publisher) Publish(topic, content string) {
	msg := Message{
		Topic:   topic,
		Content: content,
		Time:    time.Now(),
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	fmt.Printf("发布消息到主题 [%s]: %s\n", topic, content)

	// 发送给所有订阅了该主题的订阅者
	for _, sub := range p.subscribers {
		if sub.IsSubscribed(topic) {
			select {
			case sub.Messages <- msg:
			default:
				fmt.Printf("订阅者 %d 的消息队列已满，跳过消息\n", sub.ID)
			}
		}
	}
}

func (p *Publisher) Close() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, sub := range p.subscribers {
		close(sub.Messages)
	}
}

func main() {
	fmt.Println("=== 发布订阅模式演示 ===")

	publisher := NewPublisher()

	// 创建订阅者
	sub1 := NewSubscriber(1)
	sub2 := NewSubscriber(2)
	sub3 := NewSubscriber(3)

	// 添加订阅者到发布者
	publisher.AddSubscriber(sub1)
	publisher.AddSubscriber(sub2)
	publisher.AddSubscriber(sub3)

	// 订阅者订阅不同主题
	sub1.Subscribe("tech")
	sub1.Subscribe("news")
	sub2.Subscribe("tech")
	sub2.Subscribe("sports")
	sub3.Subscribe("news")
	sub3.Subscribe("sports")

	// 启动订阅者监听
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		sub1.Listen()
	}()

	go func() {
		defer wg.Done()
		sub2.Listen()
	}()

	go func() {
		defer wg.Done()
		sub3.Listen()
	}()

	// 发布消息
	time.Sleep(500 * time.Millisecond)

	publisher.Publish("tech", "Go 1.21 发布了新特性")
	time.Sleep(200 * time.Millisecond)

	publisher.Publish("news", "今日重要新闻")
	time.Sleep(200 * time.Millisecond)

	publisher.Publish("sports", "世界杯决赛结果")
	time.Sleep(200 * time.Millisecond)

	publisher.Publish("tech", "新的并发编程模式")
	time.Sleep(200 * time.Millisecond)

	// 关闭发布者
	publisher.Close()

	// 等待所有订阅者完成
	wg.Wait()
	fmt.Println("发布订阅演示完成！")
}
