package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 流水线处理演示
type DataItem struct {
	ID    int
	Value int
	Stage string
}

// 管道阶段接口
type PipelineStage interface {
	Process(input <-chan DataItem) <-chan DataItem
	GetName() string
}

// 数据生成阶段
type DataGeneratorStage struct {
	name  string
	count int
}

func NewDataGeneratorStage(count int) *DataGeneratorStage {
	return &DataGeneratorStage{
		name:  "Generator",
		count: count,
	}
}

func (g *DataGeneratorStage) Process(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)

	go func() {
		defer close(output)

		for i := 1; i <= g.count; i++ {
			item := DataItem{
				ID:    i,
				Value: rand.Intn(100),
				Stage: g.name,
			}

			fmt.Printf("%s: 生成数据 ID=%d, Value=%d\n", g.name, item.ID, item.Value)
			output <- item

			time.Sleep(100 * time.Millisecond) // 模拟生成时间
		}

		fmt.Printf("%s: 完成数据生成\n", g.name)
	}()

	return output
}

func (g *DataGeneratorStage) GetName() string {
	return g.name
}

// 数据过滤阶段
type FilterStage struct {
	name      string
	predicate func(DataItem) bool
}

func NewFilterStage(name string, predicate func(DataItem) bool) *FilterStage {
	return &FilterStage{
		name:      name,
		predicate: predicate,
	}
}

func (f *FilterStage) Process(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)

	go func() {
		defer close(output)

		for item := range input {
			if f.predicate(item) {
				item.Stage = f.name
				fmt.Printf("%s: 通过过滤 ID=%d, Value=%d\n", f.name, item.ID, item.Value)
				output <- item
			} else {
				fmt.Printf("%s: 被过滤掉 ID=%d, Value=%d\n", f.name, item.ID, item.Value)
			}
		}

		fmt.Printf("%s: 完成过滤处理\n", f.name)
	}()

	return output
}

func (f *FilterStage) GetName() string {
	return f.name
}

// 数据转换阶段
type TransformStage struct {
	name        string
	transformer func(DataItem) DataItem
}

func NewTransformStage(name string, transformer func(DataItem) DataItem) *TransformStage {
	return &TransformStage{
		name:        name,
		transformer: transformer,
	}
}

func (t *TransformStage) Process(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)

	go func() {
		defer close(output)

		for item := range input {
			transformed := t.transformer(item)
			transformed.Stage = t.name

			fmt.Printf("%s: 转换数据 ID=%d, %d -> %d\n",
				t.name, item.ID, item.Value, transformed.Value)

			output <- transformed

			time.Sleep(50 * time.Millisecond) // 模拟转换时间
		}

		fmt.Printf("%s: 完成转换处理\n", t.name)
	}()

	return output
}

func (t *TransformStage) GetName() string {
	return t.name
}

// 数据聚合阶段
type AggregateStage struct {
	name string
	size int
}

func NewAggregateStage(name string, batchSize int) *AggregateStage {
	return &AggregateStage{
		name: name,
		size: batchSize,
	}
}

func (a *AggregateStage) Process(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)

	go func() {
		defer close(output)

		batch := make([]DataItem, 0, a.size)
		batchID := 1

		for item := range input {
			batch = append(batch, item)

			if len(batch) >= a.size {
				// 计算批次汇总
				sum := 0
				for _, b := range batch {
					sum += b.Value
				}

				aggregated := DataItem{
					ID:    batchID,
					Value: sum,
					Stage: a.name,
				}

				fmt.Printf("%s: 聚合批次 %d, 包含 %d 项, 总和=%d\n",
					a.name, batchID, len(batch), sum)

				output <- aggregated

				batch = batch[:0] // 清空批次
				batchID++
			}
		}

		// 处理剩余的项目
		if len(batch) > 0 {
			sum := 0
			for _, b := range batch {
				sum += b.Value
			}

			aggregated := DataItem{
				ID:    batchID,
				Value: sum,
				Stage: a.name,
			}

			fmt.Printf("%s: 聚合最后批次 %d, 包含 %d 项, 总和=%d\n",
				a.name, batchID, len(batch), sum)

			output <- aggregated
		}

		fmt.Printf("%s: 完成聚合处理\n", a.name)
	}()

	return output
}

func (a *AggregateStage) GetName() string {
	return a.name
}

// 并行处理阶段
type ParallelStage struct {
	name       string
	workerFunc func(DataItem) DataItem
	workers    int
}

func NewParallelStage(name string, workers int, workerFunc func(DataItem) DataItem) *ParallelStage {
	return &ParallelStage{
		name:       name,
		workerFunc: workerFunc,
		workers:    workers,
	}
}

func (p *ParallelStage) Process(input <-chan DataItem) <-chan DataItem {
	output := make(chan DataItem)

	var wg sync.WaitGroup

	// 启动多个工作者
	for i := 1; i <= p.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for item := range input {
				processed := p.workerFunc(item)
				processed.Stage = fmt.Sprintf("%s-Worker%d", p.name, workerID)

				fmt.Printf("%s 工作者%d: 处理 ID=%d, %d -> %d\n",
					p.name, workerID, item.ID, item.Value, processed.Value)

				output <- processed

				time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
			}
		}(i)
	}

	// 等待所有工作者完成后关闭输出
	go func() {
		wg.Wait()
		close(output)
		fmt.Printf("%s: 所有工作者完成\n", p.name)
	}()

	return output
}

func (p *ParallelStage) GetName() string {
	return p.name
}

// 管道构建器
type Pipeline struct {
	stages []PipelineStage
	name   string
}

func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		name:   name,
		stages: make([]PipelineStage, 0),
	}
}

func (p *Pipeline) AddStage(stage PipelineStage) *Pipeline {
	p.stages = append(p.stages, stage)
	return p
}

func (p *Pipeline) Execute() <-chan DataItem {
	if len(p.stages) == 0 {
		output := make(chan DataItem)
		close(output)
		return output
	}

	fmt.Printf("=== 启动管道: %s ===\n", p.name)

	// 从第一个阶段开始
	var current <-chan DataItem = p.stages[0].Process(nil)

	// 连接所有阶段
	for i := 1; i < len(p.stages); i++ {
		stage := p.stages[i]
		fmt.Printf("连接阶段: %s\n", stage.GetName())
		current = stage.Process(current)
	}

	return current
}

func main() {
	fmt.Println("=== 流水线处理演示 ===")

	rand.Seed(time.Now().UnixNano())

	// 创建管道
	pipeline := NewPipeline("数据处理管道")

	// 添加各个阶段
	pipeline.
		AddStage(NewDataGeneratorStage(10)).                      // 生成10个数据项
		AddStage(NewFilterStage("过滤器", func(item DataItem) bool { // 过滤偶数
			return item.Value%2 == 0
		})).
		AddStage(NewTransformStage("转换器", func(item DataItem) DataItem { // 平方转换
			item.Value = item.Value * item.Value
			return item
		})).
		AddStage(NewParallelStage("并行处理器", 3, func(item DataItem) DataItem { // 并行加10
			item.Value = item.Value + 10
			return item
		})).
		AddStage(NewAggregateStage("聚合器", 3)) // 每3个聚合一次

	// 执行管道
	result := pipeline.Execute()

	// 收集最终结果
	fmt.Println("\n=== 最终结果 ===")
	var finalResults []DataItem

	for item := range result {
		finalResults = append(finalResults, item)
		fmt.Printf("最终输出: 批次ID=%d, 总和=%d, 来源=%s\n",
			item.ID, item.Value, item.Stage)
	}

	fmt.Printf("\n管道处理完成！共产生 %d 个最终结果\n", len(finalResults))

	// 计算总体统计
	totalSum := 0
	for _, item := range finalResults {
		totalSum += item.Value
	}

	fmt.Printf("所有批次总和: %d\n", totalSum)
}
