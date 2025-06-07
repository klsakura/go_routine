# Golang并发编程练习集

这是配套的练习文件集合，帮助您通过动手实践来掌握Golang并发编程技术。

## 练习结构

```
exercises/
├── simple/          # 简单级别练习 (10个练习)
├── medium/          # 中等级别练习 (10个练习)
├── hard/            # 困难级别练习 (4个练习)
├── solutions/       # 参考答案 (可选)
├── run_exercises.sh # 练习运行脚本
└── README.md        # 说明文档
```

## 如何使用练习文件

### 1. 选择练习级别
- **Simple**: 适合初学者，掌握基础概念
- **Medium**: 适合有一定基础，学习实际应用模式  
- **Hard**: 适合深入学习，挑战企业级复杂系统

### 2. 练习方法
每个练习文件都包含：
- 📋 **练习目标**: 明确学习目标
- 📝 **任务说明**: 详细的任务要求
- 💡 **代码提示**: 实现思路和方法提示
- 🔧 **TODO标记**: 明确标识需要填写代码的位置
- 🤔 **思考题**: 加深理解的反思问题

### 3. 推荐学习流程
1. **先看Demo**: 运行对应的demo文件，理解概念
2. **再做练习**: 在练习文件中实现相同或类似功能
3. **对比学习**: 将自己的实现与demo对比
4. **思考提升**: 回答思考题，加深理解

## 简单级别练习 (Simple)

| 文件 | 主题 | 核心概念 |
|------|------|----------|
| `01_basic_goroutine_exercise.go` | 基础Goroutine | go关键字、并发执行 |
| `02_waitgroup_basic_exercise.go` | WaitGroup同步 | Add/Done/Wait、同步机制 |
| `03_channel_basic_exercise.go` | 基础Channel | 发送/接收、单向channel |
| `04_buffered_channel_exercise.go` | 缓冲Channel | 缓冲区、异步通信 |
| `05_select_basic_exercise.go` | Select语句 | 多路复用、非阻塞 |
| `06_timeout_select_exercise.go` | 超时控制 | time.After、超时机制 |
| `07_mutex_basic_exercise.go` | 互斥锁 | sync.Mutex、资源保护 |
| `08_once_basic_exercise.go` | 单次执行 | sync.Once、初始化 |
| `09_channel_pipeline_exercise.go` | Channel管道 | 数据流、管道模式 |
| `10_goroutine_pool_exercise.go` | Goroutine池 | 工作者池、资源管理 |

## 中等级别练习 (Medium)

| 文件 | 主题 | 核心概念 |
|------|------|----------|
| `01_producer_consumer_exercise.go` | 生产者消费者 | 经典并发模式 |
| `02_worker_pool_advanced_exercise.go` | 高级工作池 | 任务调度、动态扩缩容 |
| `03_rate_limiter_exercise.go` | 速率限制器 | 流量控制、令牌桶 |
| `04_publish_subscribe_exercise.go` | 发布订阅 | 事件驱动、消息分发 |
| `05_context_cancellation_exercise.go` | Context取消 | 优雅退出、信号传播 |
| `06_fan_in_fan_out_exercise.go` | 扇入扇出 | 数据聚合、工作分发 |
| `07_circuit_breaker_exercise.go` | 熔断器 | 服务保护、故障处理 |
| `08_semaphore_exercise.go` | 信号量 | 资源控制、并发限制 |
| `09_actor_model_exercise.go` | Actor模型 | 消息传递、状态隔离 |
| `10_pipeline_processing_exercise.go` | 流水线处理 | 多阶段、数据流 |

## 困难级别练习 (Hard)

| 文件 | 主题 | 核心概念 |
|------|------|----------|
| `01_distributed_worker_exercise.go` | 分布式工作者 | 一致性哈希、负载均衡 |
| `02_load_balancer_exercise.go` | 负载均衡器 | 多种均衡策略、健康检查 |
| `03_message_queue_exercise.go` | 消息队列 | 重试机制、死信队列 |
| `04_connection_pool_exercise.go` | 连接池 | 资源管理、生命周期 |

## 练习技巧

### 🚀 开始练习前
1. **环境准备**: 确保Go环境正常
2. **概念复习**: 快速回顾相关理论知识
3. **Demo运行**: 先运行对应的demo了解预期效果

### 💻 编码过程中
1. **逐步实现**: 按TODO顺序逐个完成
2. **测试验证**: 每完成一个功能就测试
3. **打印调试**: 多使用fmt.Printf观察执行过程
4. **错误处理**: 注意并发程序的错误处理

### 🔍 完成后检查
1. **功能测试**: 验证所有功能是否正常
2. **并发安全**: 检查是否存在竞态条件
3. **资源清理**: 确保没有goroutine泄露
4. **性能考虑**: 思考性能优化的可能性

## 常见错误和解决方案

### ❌ 常见错误
1. **忘记关闭channel**: 导致goroutine阻塞
2. **WaitGroup计数错误**: Add和Done不匹配
3. **竞态条件**: 未保护共享资源
4. **死锁**: 循环等待资源

### ✅ 解决方案
1. **使用defer close()**: 确保channel关闭
2. **仔细管理WaitGroup**: 每个Add对应一个Done
3. **合理使用锁**: 保护共享数据访问
4. **避免循环依赖**: 设计清晰的依赖关系

## 进阶学习建议

### 📚 理论学习
- Go内存模型
- 并发模式设计
- 分布式系统原理
- 性能优化技术

### 🛠️ 实践项目
- 实现一个简单的HTTP服务器
- 构建分布式缓存系统
- 开发消息队列中间件
- 设计微服务框架

### 🔧 工具使用
```bash
# 竞态检测
go run -race exercise_file.go

# 性能分析
go build -o exercise exercise_file.go
./exercise -cpuprofile=cpu.prof

# 内存分析
go tool pprof cpu.prof
```

## 参考资源

- [Go并发编程官方文档](https://golang.org/doc/effective_go.html#concurrency)
- [Go并发模式](https://blog.golang.org/pipelines)
- [Advanced Go Concurrency Patterns](https://blog.golang.org/advanced-go-concurrency-patterns)

## 学习交流

完成练习后，建议：
1. 与同学分享实现方案
2. 讨论不同的设计思路
3. 总结遇到的问题和解决方案
4. 思考在实际项目中的应用

---

🎯 **学习目标**: 通过练习掌握Golang并发编程的核心技术  
📈 **进度跟踪**: 建议记录每个练习的完成时间和心得  
🚀 **持续改进**: 定期回顾和优化自己的实现  

Happy Coding! 🎉 