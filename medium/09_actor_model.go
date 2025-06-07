package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Actor模型演示
type Message interface {
	GetType() string
}

// 不同类型的消息
type StartMessage struct{}

func (m StartMessage) GetType() string { return "START" }

type StopMessage struct{}

func (m StopMessage) GetType() string { return "STOP" }

type DataMessage struct {
	Data   interface{}
	Sender string
}

func (m DataMessage) GetType() string { return "DATA" }

type QueryMessage struct {
	Query    string
	Response chan interface{}
}

func (m QueryMessage) GetType() string { return "QUERY" }

// Actor接口
type Actor interface {
	Start()
	Stop()
	Send(msg Message)
	GetAddress() string
}

// 基础Actor实现
type BaseActor struct {
	address  string
	mailbox  chan Message
	done     chan bool
	running  bool
	handlers map[string]func(Message)
	mu       sync.RWMutex
}

func NewBaseActor(address string, mailboxSize int) *BaseActor {
	actor := &BaseActor{
		address:  address,
		mailbox:  make(chan Message, mailboxSize),
		done:     make(chan bool),
		handlers: make(map[string]func(Message)),
	}

	// 注册默认处理器
	actor.RegisterHandler("START", actor.handleStart)
	actor.RegisterHandler("STOP", actor.handleStop)

	return actor
}

func (a *BaseActor) RegisterHandler(msgType string, handler func(Message)) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.handlers[msgType] = handler
}

func (a *BaseActor) Start() {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return
	}
	a.running = true
	a.mu.Unlock()

	fmt.Printf("Actor %s 启动\n", a.address)

	go a.messageLoop()
}

func (a *BaseActor) Stop() {
	a.mu.RLock()
	if !a.running {
		a.mu.RUnlock()
		return
	}
	a.mu.RUnlock()

	a.Send(StopMessage{})
}

func (a *BaseActor) Send(msg Message) {
	select {
	case a.mailbox <- msg:
	default:
		fmt.Printf("Actor %s 邮箱已满，消息被丢弃: %s\n", a.address, msg.GetType())
	}
}

func (a *BaseActor) GetAddress() string {
	return a.address
}

func (a *BaseActor) messageLoop() {
	for {
		select {
		case msg := <-a.mailbox:
			a.processMessage(msg)
		case <-a.done:
			return
		}
	}
}

func (a *BaseActor) processMessage(msg Message) {
	a.mu.RLock()
	handler, exists := a.handlers[msg.GetType()]
	a.mu.RUnlock()

	if exists {
		handler(msg)
	} else {
		fmt.Printf("Actor %s 收到未知消息类型: %s\n", a.address, msg.GetType())
	}
}

func (a *BaseActor) handleStart(msg Message) {
	fmt.Printf("Actor %s 收到启动消息\n", a.address)
}

func (a *BaseActor) handleStop(msg Message) {
	fmt.Printf("Actor %s 收到停止消息，准备关闭\n", a.address)
	a.mu.Lock()
	a.running = false
	a.mu.Unlock()
	close(a.done)
}

// 计算器Actor
type CalculatorActor struct {
	*BaseActor
	result float64
}

func NewCalculatorActor(address string) *CalculatorActor {
	calc := &CalculatorActor{
		BaseActor: NewBaseActor(address, 100),
		result:    0,
	}

	// 注册计算器特定的处理器
	calc.RegisterHandler("ADD", calc.handleAdd)
	calc.RegisterHandler("MULTIPLY", calc.handleMultiply)
	calc.RegisterHandler("QUERY", calc.handleQuery)

	return calc
}

type AddMessage struct {
	Value float64
}

func (m AddMessage) GetType() string { return "ADD" }

type MultiplyMessage struct {
	Value float64
}

func (m MultiplyMessage) GetType() string { return "MULTIPLY" }

func (c *CalculatorActor) handleAdd(msg Message) {
	addMsg := msg.(AddMessage)
	c.result += addMsg.Value
	fmt.Printf("Calculator %s: 加法操作 +%.2f, 结果=%.2f\n",
		c.address, addMsg.Value, c.result)
}

func (c *CalculatorActor) handleMultiply(msg Message) {
	mulMsg := msg.(MultiplyMessage)
	c.result *= mulMsg.Value
	fmt.Printf("Calculator %s: 乘法操作 *%.2f, 结果=%.2f\n",
		c.address, mulMsg.Value, c.result)
}

func (c *CalculatorActor) handleQuery(msg Message) {
	queryMsg := msg.(QueryMessage)
	if queryMsg.Query == "result" {
		queryMsg.Response <- c.result
		fmt.Printf("Calculator %s: 查询结果=%.2f\n", c.address, c.result)
	}
}

// 日志Actor
type LoggerActor struct {
	*BaseActor
	logs []string
}

func NewLoggerActor(address string) *LoggerActor {
	logger := &LoggerActor{
		BaseActor: NewBaseActor(address, 1000),
		logs:      make([]string, 0),
	}

	logger.RegisterHandler("LOG", logger.handleLog)
	logger.RegisterHandler("QUERY", logger.handleQuery)

	return logger
}

type LogMessage struct {
	Level   string
	Content string
}

func (m LogMessage) GetType() string { return "LOG" }

func (l *LoggerActor) handleLog(msg Message) {
	logMsg := msg.(LogMessage)
	logEntry := fmt.Sprintf("[%s] %s: %s",
		time.Now().Format("15:04:05"), logMsg.Level, logMsg.Content)
	l.logs = append(l.logs, logEntry)
	fmt.Printf("Logger %s: %s\n", l.address, logEntry)
}

func (l *LoggerActor) handleQuery(msg Message) {
	queryMsg := msg.(QueryMessage)
	if queryMsg.Query == "count" {
		queryMsg.Response <- len(l.logs)
	} else if queryMsg.Query == "logs" {
		queryMsg.Response <- l.logs
	}
}

// Actor系统
type ActorSystem struct {
	actors map[string]Actor
	mu     sync.RWMutex
}

func NewActorSystem() *ActorSystem {
	return &ActorSystem{
		actors: make(map[string]Actor),
	}
}

func (as *ActorSystem) RegisterActor(actor Actor) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.actors[actor.GetAddress()] = actor
	fmt.Printf("注册Actor: %s\n", actor.GetAddress())
}

func (as *ActorSystem) GetActor(address string) Actor {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.actors[address]
}

func (as *ActorSystem) StartAll() {
	as.mu.RLock()
	defer as.mu.RUnlock()

	for _, actor := range as.actors {
		actor.Start()
	}
}

func (as *ActorSystem) StopAll() {
	as.mu.RLock()
	defer as.mu.RUnlock()

	for _, actor := range as.actors {
		actor.Stop()
	}
}

func main() {
	fmt.Println("=== Actor模型演示 ===")

	rand.Seed(time.Now().UnixNano())

	// 创建Actor系统
	system := NewActorSystem()

	// 创建Actors
	calc1 := NewCalculatorActor("calculator-1")
	calc2 := NewCalculatorActor("calculator-2")
	logger := NewLoggerActor("logger")

	// 注册Actors
	system.RegisterActor(calc1)
	system.RegisterActor(calc2)
	system.RegisterActor(logger)

	// 启动所有Actors
	system.StartAll()

	time.Sleep(500 * time.Millisecond)

	// 向日志Actor发送日志消息
	logger.Send(LogMessage{Level: "INFO", Content: "系统启动"})

	// 向计算器Actor发送计算消息
	calc1.Send(AddMessage{Value: 10})
	calc1.Send(AddMessage{Value: 5})
	calc1.Send(MultiplyMessage{Value: 2})

	calc2.Send(AddMessage{Value: 20})
	calc2.Send(MultiplyMessage{Value: 3})

	// 并发发送更多消息
	var wg sync.WaitGroup

	// 模拟多个客户端向计算器发送请求
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			calcAddr := "calculator-1"
			if clientID%2 == 0 {
				calcAddr = "calculator-2"
			}

			calc := system.GetActor(calcAddr)
			if calc != nil {
				calc.Send(AddMessage{Value: float64(clientID)})
				time.Sleep(100 * time.Millisecond)
				calc.Send(MultiplyMessage{Value: 1.5})
			}

			// 记录日志
			logger.Send(LogMessage{
				Level:   "INFO",
				Content: fmt.Sprintf("客户端 %d 完成计算", clientID),
			})
		}(i)
	}

	wg.Wait()
	time.Sleep(1 * time.Second)

	// 查询结果
	fmt.Println("\n=== 查询Actor状态 ===")

	// 查询计算器结果
	for _, calcAddr := range []string{"calculator-1", "calculator-2"} {
		calc := system.GetActor(calcAddr)
		if calc != nil {
			response := make(chan interface{})
			calc.Send(QueryMessage{Query: "result", Response: response})

			select {
			case result := <-response:
				fmt.Printf("%s 的最终结果: %.2f\n", calcAddr, result)
			case <-time.After(1 * time.Second):
				fmt.Printf("%s 查询超时\n", calcAddr)
			}
		}
	}

	// 查询日志数量
	response := make(chan interface{})
	logger.Send(QueryMessage{Query: "count", Response: response})

	select {
	case count := <-response:
		fmt.Printf("日志记录数量: %d\n", count)
	case <-time.After(1 * time.Second):
		fmt.Println("日志查询超时")
	}

	// 停止所有Actors
	fmt.Println("\n=== 停止Actor系统 ===")
	system.StopAll()

	time.Sleep(1 * time.Second)
	fmt.Println("Actor模型演示完成！")
}
