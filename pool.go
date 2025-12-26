package mredis

import (
	"context"
	"fmt"
	"mlog"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// Pool Redis连接池管理器
type Pool struct {
	client *goredis.Client
	config *RedisConfig
	mu     sync.RWMutex
	closed bool
}

// NewPool 创建新的Redis连接池实例
// 返回:
//   - *Pool: 连接池对象，初始状态为关闭，需要调用LoadConfig后才能使用
//
// 注意:
//   - 新创建的连接池处于关闭状态
//   - 必须先调用LoadConfig加载配置才能正常使用
func NewPool() *Pool {
	mlog.Info("NewPool redis")
	return &Pool{
		closed: true, // 初始状态为关闭，需要LoadConfig后才能使用
	}
}

// LoadConfig 加载Redis配置并初始化连接池
// 参数:
//   - config: Redis配置对象
//
// 返回:
//   - error: 配置加载失败时返回错误
//
// 功能:
//   - 验证配置有效性（Addr为必填项）
//   - 关闭旧连接（如果存在）
//   - 设置默认值（PoolSize=10, MinIdleConns=2等）
//   - 创建新的Redis客户端
//   - 执行PING测试连接
//
// 注意:
//   - 此方法会关闭现有连接，请谨慎调用
//   - 配置中未设置的参数将使用默认值
func (p *Pool) LoadConfig(config *RedisConfig) error {
	if config == nil {
		return ErrInvalidConfig
	}

	// 验证必要配置
	if config.Addr == "" {
		return fmt.Errorf("%w: addr is required", ErrInvalidConfig)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// 关闭旧连接
	if p.client != nil {
		if err := p.client.Close(); err != nil {
			mlog.Warn("close old redis client error: %v", err)
		}
	}

	// 设置默认值
	if config.PoolSize <= 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns <= 0 {
		config.MinIdleConns = 2
	}
	if config.DialTimeout <= 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout <= 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout <= 0 {
		config.WriteTimeout = 3 * time.Second
	}

	// 创建新连接
	opts := &goredis.Options{
		Addr:            config.Addr,
		Password:        config.Pass,
		DB:              config.Index,
		PoolSize:        config.PoolSize,
		MinIdleConns:    config.MinIdleConns,
		MaxIdleConns:    config.MaxIdleConns,
		DialTimeout:     config.DialTimeout,
		ReadTimeout:     config.ReadTimeout,
		WriteTimeout:    config.WriteTimeout,
		PoolTimeout:     config.PoolTimeout,
		ConnMaxIdleTime: config.IdleTimeout,
		MaxRetries:      config.MaxRetries,
		MinRetryBackoff: config.MinRetryBackoff,
		MaxRetryBackoff: config.MaxRetryBackoff,
	}

	p.client = goredis.NewClient(opts)
	p.config = config
	p.closed = false

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := p.client.Ping(ctx).Err(); err != nil {
		p.closed = true
		return fmt.Errorf("redis connection test failed: %w", err)
	}

	mlog.Info("Redis pool loaded successfully, addr: %s, db: %d", config.Addr, config.Index)
	return nil
}

// GetClient 获取Redis客户端实例
// 返回:
//   - *goredis.Client: Redis客户端对象，如果连接池已关闭或未初始化则返回nil
//
// 注意:
//   - 返回的客户端由go-redis内置连接池管理
//   - 无需手动释放连接
//   - 使用前应检查返回值是否为nil
func (p *Pool) GetClient() *goredis.Client {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		mlog.Error("connection pool is closed")
		return nil
	}

	if p.client == nil {
		mlog.Error("redis client is not initialized")
		return nil
	}

	return p.client
}

// Close 关闭Redis连接池
// 功能:
//   - 关闭所有活动连接
//   - 将连接池状态标记为已关闭
//
// 注意:
//   - 关闭后需要重新调用LoadConfig才能继续使用
//   - 建议在应用程序退出时调用
//   - 多次调用是安全的
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		err := p.client.Close()
		if err != nil {
			mlog.Error("close redis pool error: %v", err)
		}
	}

	p.closed = true
	mlog.Info("Redis connection pool closed")
}

// IsHealthy 检查Redis连接池健康状态
// 返回:
//   - bool: 连接池健康返回true，否则返回false
//
// 功能:
//   - 检查客户端是否已初始化
//   - 执行PING命令测试连接
//
// 注意:
//   - 超时时间为2秒
//   - 如果连接池未初始化或PING失败则返回false
func (p *Pool) IsHealthy() bool {
	client := p.GetClient()
	if client == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return client.Ping(ctx).Err() == nil
}

// GetPoolStats 获取Redis连接池统计信息
// 返回:
//   - *goredis.PoolStats: 连接池统计信息对象，如果客户端未初始化则返回nil
//
// 统计信息包括:
//   - Hits: 从连接池获取连接的成功次数
//   - Misses: 从连接池获取连接的失败次数
//   - Timeouts: 获取连接超时次数
//   - TotalConns: 总连接数
//   - IdleConns: 空闲连接数
//   - StaleConns: 过期连接数
func (p *Pool) GetPoolStats() *goredis.PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.client == nil {
		return nil
	}

	return p.client.PoolStats()
}

// IsClosed 检查连接池是否已关闭
// 返回:
//   - bool: 连接池已关闭返回true，否则返回false
//
// 注意:
//   - 此方法是线程安全的
func (p *Pool) IsClosed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.closed
}
