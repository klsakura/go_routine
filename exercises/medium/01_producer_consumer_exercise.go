/*
Golang并发编程练习 - 中等级别
练习文件：01_producer_consumer_exercise.go
练习主题：生产者-消费者模式

练习目标：
1. 实现经典的生产者-消费者模式
2. 掌握缓冲channel在此模式中的应用
3. 处理多生产者多消费者场景
4. 学会优雅地关闭和同步

练习任务：
- 任务1：单生产者单消费者
- 任务2：多生产者单消费者
- 任务3：单生产者多消费者
- 任务4：多生产者多消费者与负载均衡

运行方式：go run exercises/medium/01_producer_consumer_exercise.go
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Product 产品结构
type Product struct {
	ID       int
	Name     string
	Price    float64
	Category string
}

// Producer 生产者结构
type Producer struct {
	ID   string
	Type string
}

// Consumer 消费者结构
type Consumer struct {
	ID             string
	ProcessingTime time.Duration
}

// TODO: 实现Producer的Produce方法
// 生产指定数量的产品到channel
func (p *Producer) Produce(products chan<- Product, count int, wg *sync.WaitGroup) {
	// 在这里实现您的代码
	// 提示：
	// 1. 使用defer wg.Done()
	// 2. 根据生产者类型生成不同类别的产品
	// 3. 每个产品有唯一ID
	// 4. 模拟生产时间

}

// TODO: 实现Consumer的Consume方法
// 从channel消费产品
func (c *Consumer) Consume(products <-chan Product, wg *sync.WaitGroup) {
	// 在这里实现您的代码
	// 提示：
	// 1. 使用defer wg.Done()
	// 2. 使用range遍历channel
	// 3. 根据消费者的ProcessingTime模拟处理时间
	// 4. 打印消费信息

}

func main() {
	fmt.Println("=== 生产者-消费者模式练习 ===")
	rand.Seed(time.Now().UnixNano())

	// 任务1：单生产者单消费者
	fmt.Println("\n任务1：单生产者单消费者")
	// TODO:
	// 1. 创建一个缓冲为5的Product channel
	// 2. 创建一个生产者（ID: "P1", Type: "Electronics"）
	// 3. 创建一个消费者（ID: "C1", ProcessingTime: 200ms）
	// 4. 生产者生产10个产品
	// 5. 使用WaitGroup同步

	fmt.Println("任务1完成\n")

	// 任务2：多生产者单消费者
	fmt.Println("任务2：多生产者单消费者")
	// TODO:
	// 1. 创建一个缓冲为10的Product channel
	// 2. 创建3个不同类型的生产者
	// 3. 创建1个消费者
	// 4. 每个生产者生产8个产品
	// 5. 注意：需要在所有生产者完成后关闭channel

	fmt.Println("任务2完成\n")

	// 任务3：单生产者多消费者
	fmt.Println("任务3：单生产者多消费者")
	// TODO:
	// 1. 创建一个缓冲为15的Product channel
	// 2. 创建1个生产者
	// 3. 创建4个消费者，每个有不同的处理时间
	// 4. 生产者生产20个产品
	// 5. 观察消费者之间的竞争

	fmt.Println("任务3完成\n")

	// 任务4：多生产者多消费者
	fmt.Println("任务4：多生产者多消费者系统")
	// TODO:
	// 1. 创建一个缓冲为20的Product channel
	// 2. 创建5个不同类型的生产者
	// 3. 创建3个不同速度的消费者
	// 4. 每个生产者生产随机数量(5-15)的产品
	// 5. 实现生产者统计和消费者统计

	fmt.Println("任务4完成\n")

	// 任务5：带优先级的生产消费
	fmt.Println("任务5：优先级处理（挑战任务）")
	// TODO: 高级挑战
	// 实现一个支持优先级的生产者消费者系统
	// 提示：可以使用多个channel或者自定义排序

	fmt.Println("所有练习完成！")

	// 反思问题：
	fmt.Println("\n思考题：")
	fmt.Println("1. 缓冲区大小如何影响系统性能？")
	fmt.Println("2. 如何确保所有生产者完成后正确关闭channel？")
	fmt.Println("3. 多消费者竞争时，如何保证负载均衡？")
	fmt.Println("4. 在什么情况下需要多个channel？")
}
