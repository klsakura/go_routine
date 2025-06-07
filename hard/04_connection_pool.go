package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// 连接池演示
type Connection interface {
	Connect() error
	Close() error
	Execute(query string) (interface{}, error)
	IsAlive() bool
	GetID() string
	GetCreatedTime() time.Time
	GetLastUsed() time.Time
	SetLastUsed(time.Time)
}

// 模拟数据库连接
type DBConnection struct {
	ID          string
	connected   bool
	createdTime time.Time
	lastUsed    time.Time
	queries     int64
	mu          sync.RWMutex
}

func NewDBConnection(id string) *DBConnection {
	return &DBConnection{
		ID:          id,
		createdTime: time.Now(),
		lastUsed:    time.Now(),
	}
}

func (c *DBConnection) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("connection %s already connected", c.ID)
	}

	// 模拟连接时间
	time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)

	// 模拟连接失败（5%概率）
	if rand.Float32() < 0.05 {
		return fmt.Errorf("failed to connect %s", c.ID)
	}

	c.connected = true
	fmt.Printf("连接 %s 已建立\n", c.ID)
	return nil
}

func (c *DBConnection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return fmt.Errorf("connection %s not connected", c.ID)
	}

	c.connected = false
	fmt.Printf("连接 %s 已关闭\n", c.ID)
	return nil
}

func (c *DBConnection) Execute(query string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil, fmt.Errorf("connection %s not connected", c.ID)
	}

	// 模拟查询执行时间
	time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)

	// 模拟查询失败（10%概率）
	if rand.Float32() < 0.1 {
		return nil, fmt.Errorf("query failed on connection %s", c.ID)
	}

	atomic.AddInt64(&c.queries, 1)
	c.lastUsed = time.Now()

	result := fmt.Sprintf("Result from %s: %s", c.ID, query)
	return result, nil
}

func (c *DBConnection) IsAlive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 模拟连接检查失败（3%概率）
	if rand.Float32() < 0.03 {
		c.connected = false
		return false
	}

	return c.connected
}

func (c *DBConnection) GetID() string {
	return c.ID
}

func (c *DBConnection) GetCreatedTime() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.createdTime
}

func (c *DBConnection) GetLastUsed() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastUsed
}

func (c *DBConnection) SetLastUsed(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastUsed = t
}

func (c *DBConnection) GetQueries() int64 {
	return atomic.LoadInt64(&c.queries)
}

// 连接池配置
type PoolConfig struct {
	MinConnections    int           // 最小连接数
	MaxConnections    int           // 最大连接数
	MaxIdleTime       time.Duration // 最大空闲时间
	ConnectionTimeout time.Duration // 连接超时时间
	HealthCheckPeriod time.Duration // 健康检查周期
}

// 连接池实现
type ConnectionPool struct {
	config      PoolConfig
	connections chan Connection
	active      map[string]Connection
	factory     func(id string) Connection
	mu          sync.RWMutex
	closed      bool
	wg          sync.WaitGroup
	stopCh      chan bool
	stats       struct {
		created     int64
		borrowed    int64
		returned    int64
		failed      int64
		evicted     int64
		healthCheck int64
	}
}

func NewConnectionPool(config PoolConfig, factory func(id string) Connection) (*ConnectionPool, error) {
	if config.MinConnections < 0 || config.MaxConnections < config.MinConnections {
		return nil, fmt.Errorf("invalid pool configuration")
	}

	pool := &ConnectionPool{
		config:      config,
		connections: make(chan Connection, config.MaxConnections),
		active:      make(map[string]Connection),
		factory:     factory,
		stopCh:      make(chan bool),
	}

	// 创建最小连接数
	for i := 0; i < config.MinConnections; i++ {
		conn := pool.createConnection()
		if conn != nil {
			pool.connections <- conn
		}
	}

	// 启动健康检查
	pool.wg.Add(1)
	go pool.healthChecker()

	// 启动空闲连接清理
	pool.wg.Add(1)
	go pool.idleConnectionEvector()

	return pool, nil
}

func (p *ConnectionPool) createConnection() Connection {
	id := fmt.Sprintf("conn-%d", atomic.AddInt64(&p.stats.created, 1))
	conn := p.factory(id)

	err := conn.Connect()
	if err != nil {
		atomic.AddInt64(&p.stats.failed, 1)
		fmt.Printf("创建连接失败: %v\n", err)
		return nil
	}

	return conn
}

func (p *ConnectionPool) BorrowConnection() (Connection, error) {
	if p.closed {
		return nil, fmt.Errorf("connection pool is closed")
	}

	atomic.AddInt64(&p.stats.borrowed, 1)

	select {
	case conn := <-p.connections:
		// 检查连接是否仍然有效
		if conn.IsAlive() {
			p.mu.Lock()
			p.active[conn.GetID()] = conn
			p.mu.Unlock()

			conn.SetLastUsed(time.Now())
			return conn, nil
		} else {
			// 连接已失效，创建新连接
			atomic.AddInt64(&p.stats.evicted, 1)
			conn.Close()

			newConn := p.createConnection()
			if newConn != nil {
				p.mu.Lock()
				p.active[newConn.GetID()] = newConn
				p.mu.Unlock()
				return newConn, nil
			}
		}

	case <-time.After(p.config.ConnectionTimeout):
		return nil, fmt.Errorf("connection timeout")
	}

	// 如果池中没有可用连接，尝试创建新连接
	p.mu.RLock()
	activeCount := len(p.active)
	p.mu.RUnlock()

	if activeCount < p.config.MaxConnections {
		conn := p.createConnection()
		if conn != nil {
			p.mu.Lock()
			p.active[conn.GetID()] = conn
			p.mu.Unlock()
			return conn, nil
		}
	}

	return nil, fmt.Errorf("no available connections")
}

func (p *ConnectionPool) ReturnConnection(conn Connection) error {
	if p.closed {
		conn.Close()
		return fmt.Errorf("connection pool is closed")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.active[conn.GetID()]; !exists {
		return fmt.Errorf("connection not from this pool")
	}

	delete(p.active, conn.GetID())
	atomic.AddInt64(&p.stats.returned, 1)

	// 检查连接是否仍然有效
	if conn.IsAlive() {
		select {
		case p.connections <- conn:
			return nil
		default:
			// 池已满，关闭连接
			conn.Close()
			return nil
		}
	} else {
		// 连接已失效，关闭它
		atomic.AddInt64(&p.stats.evicted, 1)
		conn.Close()
		return nil
	}
}

func (p *ConnectionPool) healthChecker() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.config.HealthCheckPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.performHealthCheck()
		case <-p.stopCh:
			return
		}
	}
}

func (p *ConnectionPool) performHealthCheck() {
	atomic.AddInt64(&p.stats.healthCheck, 1)

	// 检查池中的空闲连接
	poolSize := len(p.connections)
	for i := 0; i < poolSize; i++ {
		select {
		case conn := <-p.connections:
			if conn.IsAlive() {
				// 连接健康，放回池中
				p.connections <- conn
			} else {
				// 连接不健康，关闭并可能创建新连接
				atomic.AddInt64(&p.stats.evicted, 1)
				conn.Close()

				// 如果连接数低于最小值，创建新连接
				if len(p.connections) < p.config.MinConnections {
					newConn := p.createConnection()
					if newConn != nil {
						p.connections <- newConn
					}
				}
			}
		default:
			break
		}
	}
}

func (p *ConnectionPool) idleConnectionEvector() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.config.MaxIdleTime / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.evictIdleConnections()
		case <-p.stopCh:
			return
		}
	}
}

func (p *ConnectionPool) evictIdleConnections() {
	now := time.Now()
	poolSize := len(p.connections)

	for i := 0; i < poolSize; i++ {
		select {
		case conn := <-p.connections:
			if now.Sub(conn.GetLastUsed()) > p.config.MaxIdleTime {
				// 连接空闲时间过长，关闭它
				atomic.AddInt64(&p.stats.evicted, 1)
				conn.Close()

				// 确保不低于最小连接数
				if len(p.connections) < p.config.MinConnections {
					newConn := p.createConnection()
					if newConn != nil {
						p.connections <- newConn
					}
				}
			} else {
				// 连接仍在有效期内，放回池中
				p.connections <- conn
			}
		default:
			break
		}
	}
}

func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return fmt.Errorf("pool already closed")
	}
	p.closed = true
	p.mu.Unlock()

	// 停止后台任务
	close(p.stopCh)
	p.wg.Wait()

	// 关闭所有连接
	close(p.connections)
	for conn := range p.connections {
		conn.Close()
	}

	// 关闭活跃连接
	p.mu.RLock()
	for _, conn := range p.active {
		conn.Close()
	}
	p.mu.RUnlock()

	fmt.Println("连接池已关闭")
	return nil
}

func (p *ConnectionPool) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"idle_connections":   len(p.connections),
		"active_connections": len(p.active),
		"created":            atomic.LoadInt64(&p.stats.created),
		"borrowed":           atomic.LoadInt64(&p.stats.borrowed),
		"returned":           atomic.LoadInt64(&p.stats.returned),
		"failed":             atomic.LoadInt64(&p.stats.failed),
		"evicted":            atomic.LoadInt64(&p.stats.evicted),
		"health_checks":      atomic.LoadInt64(&p.stats.healthCheck),
	}
}

// 数据库客户端
type DBClient struct {
	pool *ConnectionPool
}

func NewDBClient(pool *ConnectionPool) *DBClient {
	return &DBClient{pool: pool}
}

func (c *DBClient) Query(query string) (interface{}, error) {
	conn, err := c.pool.BorrowConnection()
	if err != nil {
		return nil, err
	}
	defer c.pool.ReturnConnection(conn)

	return conn.Execute(query)
}

func main() {
	fmt.Println("=== 连接池演示 ===")

	rand.Seed(time.Now().UnixNano())

	// 创建连接池配置
	config := PoolConfig{
		MinConnections:    3,
		MaxConnections:    10,
		MaxIdleTime:       5 * time.Second,
		ConnectionTimeout: 3 * time.Second,
		HealthCheckPeriod: 2 * time.Second,
	}

	// 创建连接池
	pool, err := NewConnectionPool(config, func(id string) Connection {
		return NewDBConnection(id)
	})
	if err != nil {
		fmt.Printf("创建连接池失败: %v\n", err)
		return
	}
	defer pool.Close()

	// 创建数据库客户端
	client := NewDBClient(pool)

	fmt.Printf("连接池创建成功，配置: 最小=%d, 最大=%d, 空闲超时=%v\n",
		config.MinConnections, config.MaxConnections, config.MaxIdleTime)

	// 并发测试
	var wg sync.WaitGroup
	numClients := 20
	queriesPerClient := 5

	fmt.Printf("\n启动 %d 个并发客户端，每个执行 %d 次查询\n", numClients, queriesPerClient)

	for i := 1; i <= numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			for j := 1; j <= queriesPerClient; j++ {
				query := fmt.Sprintf("SELECT * FROM table WHERE client=%d AND seq=%d", clientID, j)

				result, err := client.Query(query)
				if err != nil {
					fmt.Printf("客户端 %d 查询 %d 失败: %v\n", clientID, j, err)
				} else {
					fmt.Printf("客户端 %d 查询 %d 成功: %v\n", clientID, j, result)
				}

				// 随机等待时间
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			}

			fmt.Printf("客户端 %d 完成所有查询\n", clientID)
		}(i)
	}

	// 在测试期间定期打印统计信息
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			stats := pool.GetStats()
			fmt.Printf("\n--- 连接池统计 (第%d次) ---\n", i+1)
			for key, value := range stats {
				fmt.Printf("%s: %v\n", key, value)
			}
		}
	}()

	wg.Wait()

	// 最终统计
	fmt.Println("\n=== 最终统计 ===")
	stats := pool.GetStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	// 等待一段时间观察空闲连接清理
	fmt.Println("\n等待空闲连接清理...")
	time.Sleep(6 * time.Second)

	finalStats := pool.GetStats()
	fmt.Println("\n=== 清理后统计 ===")
	for key, value := range finalStats {
		fmt.Printf("%s: %v\n", key, value)
	}

	fmt.Println("\n连接池演示完成！")
}
