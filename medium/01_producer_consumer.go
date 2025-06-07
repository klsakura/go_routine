package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 生产者消费者模式演示
type Product struct {
	ID    int
	Name  string
	Price float64
}

func producer(id int, products chan<- Product, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 5; i++ {
		product := Product{
			ID:    id*100 + i,
			Name:  fmt.Sprintf("产品-%d-%d", id, i),
			Price: rand.Float64() * 100,
		}

		fmt.Printf("生产者 %d 生产: %+v\n", id, product)
		products <- product

		// 随机生产间隔
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}

	fmt.Printf("生产者 %d 完成生产\n", id)
}

func consumer(id int, products <-chan Product, wg *sync.WaitGroup) {
	defer wg.Done()

	for product := range products {
		fmt.Printf("消费者 %d 消费: ID=%d, Name=%s, Price=%.2f\n",
			id, product.ID, product.Name, product.Price)

		// 模拟消费时间
		time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)
	}

	fmt.Printf("消费者 %d 结束消费\n", id)
}

func main() {
	fmt.Println("=== 生产者消费者模式演示 ===")

	rand.Seed(time.Now().UnixNano())

	const bufferSize = 10
	const numProducers = 3
	const numConsumers = 2

	products := make(chan Product, bufferSize)

	var producerWg, consumerWg sync.WaitGroup

	// 启动生产者
	for i := 1; i <= numProducers; i++ {
		producerWg.Add(1)
		go producer(i, products, &producerWg)
	}

	// 启动消费者
	for i := 1; i <= numConsumers; i++ {
		consumerWg.Add(1)
		go consumer(i, products, &consumerWg)
	}

	// 等待所有生产者完成，然后关闭channel
	go func() {
		producerWg.Wait()
		close(products)
		fmt.Println("所有生产者完成，产品channel已关闭")
	}()

	// 等待所有消费者完成
	consumerWg.Wait()
	fmt.Println("所有消费者完成！")
}
