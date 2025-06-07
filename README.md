# Golang并发编程学习Demo

这是一个完整的Golang并发编程学习项目，包含从简单到困难三个级别的实用demo。

## 目录结构

```
.
├── simple/          # 简单级别 (10个demo)
├── medium/          # 中等级别 (10个demo)  
├── hard/            # 困难级别 (4个demo)
├── run_all.sh       # 运行脚本
└── README.md        # 说明文档
```

## 简单级别 (Simple)

基础并发概念和简单使用场景：

1. **01_basic_goroutine.go** - 基础goroutine使用和并发执行
2. **02_waitgroup_basic.go** - WaitGroup基础使用和同步
3. **03_channel_basic.go** - 基础channel通信和数据传递
4. **04_buffered_channel.go** - 缓冲channel使用和容量管理
5. **05_select_basic.go** - 基础select语句和多路复用
6. **06_timeout_select.go** - 带超时的select和时间控制
7. **07_mutex_basic.go** - 基础互斥锁使用和共享资源保护
8. **08_once_basic.go** - sync.Once使用和单次初始化
9. **09_channel_pipeline.go** - 简单的channel管道和数据流
10. **10_goroutine_pool.go** - 简单的goroutine池和工作者模式

## 中等级别 (Medium)

实际应用中的并发模式：

1. **01_producer_consumer.go** - 生产者消费者模式和缓冲处理
2. **02_worker_pool_advanced.go** - 高级工作池和任务调度
3. **03_rate_limiter.go** - 速率限制器和流量控制
4. **04_publish_subscribe.go** - 发布订阅模式和事件分发
5. **05_context_cancellation.go** - Context取消机制和优雅退出
6. **06_fan_in_fan_out.go** - 扇入扇出模式和工作分发
7. **07_circuit_breaker.go** - 熔断器模式和服务保护
8. **08_semaphore.go** - 信号量实现和资源控制
9. **09_actor_model.go** - Actor模型和消息传递
10. **10_pipeline_processing.go** - 流水线处理和多阶段数据处理

## 困难级别 (Hard)

高级并发编程技术和复杂系统：

1. **01_distributed_worker.go** - 分布式工作者和一致性哈希
2. **02_load_balancer.go** - 负载均衡器和多种均衡策略
3. **03_message_queue.go** - 消息队列系统和重试机制
4. **04_connection_pool.go** - 连接池管理和资源生命周期

**注意：** Hard级别目前包含4个高质量的企业级并发编程示例，每个都是完整的系统实现，涵盖了分布式系统、负载均衡、消息队列和连接池等核心技术。

## 如何使用

### 运行单个demo
```bash
# 进入对应目录
cd simple
go run 01_basic_goroutine.go

# 或者使用完整路径
go run simple/01_basic_goroutine.go
```

### 运行所有demo
```bash
# 创建运行脚本
chmod +x run_all.sh
./run_all.sh
```

### 注意事项
- 由于每个demo都是独立的main程序，在同一个目录下运行时会出现"main redeclared"的linter警告，这是正常现象
- 建议一次运行一个demo文件来学习，而不是同时编译整个目录
- 可以使用运行脚本来方便地选择和运行特定的demo

## 学习建议

1. **按顺序学习**: 从simple -> medium -> hard 逐步深入
2. **动手实践**: 不仅要看代码，更要运行和修改
3. **理解原理**: 理解每个demo背后的并发原理
4. **举一反三**: 尝试修改参数，观察不同的行为
5. **实际应用**: 思考如何在实际项目中应用这些模式

## 核心概念

### Goroutine
- 轻量级线程，由Go运行时管理
- 栈大小可动态调整（初始2KB）
- 通过`go`关键字启动

### Channel
- Goroutine间通信的管道
- 无缓冲channel：同步通信
- 缓冲channel：异步通信
- 单向channel：限制读写权限

### Select
- 多路复用channel操作
- 非阻塞通信
- 超时控制

### 同步原语
- **sync.Mutex**: 互斥锁
- **sync.RWMutex**: 读写锁
- **sync.WaitGroup**: 等待组
- **sync.Once**: 单次执行
- **sync.Atomic**: 原子操作

### Context
- 取消信号传播
- 超时控制
- 值传递

## 并发模式

### 基础模式
- **Pipeline**: 数据流水线处理
- **Fan-in/Fan-out**: 数据聚合和分发
- **Worker Pool**: 工作者池

### 高级模式
- **Producer-Consumer**: 生产者消费者
- **Publish-Subscribe**: 发布订阅
- **Circuit Breaker**: 熔断器
- **Rate Limiter**: 速率限制
- **Actor Model**: Actor模型
- **Semaphore**: 信号量控制

### 企业级模式
- **Message Queue**: 消息队列
- **Connection Pool**: 连接池
- **Load Balancer**: 负载均衡
- **Distributed Worker**: 分布式工作者

## 最佳实践

1. **避免共享内存**: 优先使用channel通信
2. **合理使用buffer**: 根据场景选择缓冲大小
3. **及时关闭channel**: 避免goroutine泄露
4. **使用context**: 实现优雅的取消机制
5. **错误处理**: 并发程序中的错误处理要更谨慎

## 常见陷阱

1. **Goroutine泄露**: 忘记关闭channel或context
2. **竞态条件**: 多个goroutine同时访问共享数据
3. **死锁**: 循环等待资源
4. **数据竞争**: 未同步的内存访问

## 性能考虑

1. **Goroutine数量**: 不是越多越好，要根据CPU核数调整
2. **Channel缓冲**: 合理的缓冲可以减少阻塞
3. **锁的粒度**: 细粒度锁可以提高并发性
4. **原子操作**: 简单操作优先使用atomic包

## 调试工具

```bash
# 竞态检测
go run -race main.go

# 性能分析
go build -o app main.go
./app -cpuprofile=cpu.prof -memprofile=mem.prof

# 查看goroutine
kill -QUIT <pid>  # 发送SIGQUIT信号
```

## 扩展学习

- [Go并发编程官方文档](https://golang.org/doc/effective_go.html#concurrency)
- [Go内存模型](https://golang.org/ref/mem)
- [Go竞态检测器](https://golang.org/doc/articles/race_detector.html)

Happy Coding! 🚀 