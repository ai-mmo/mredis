package mredis

import (
	"testing"
	"time"
)

func TestLoadConfigError(t *testing.T) {
	// 测试nil配置
	err := LoadConfig(nil)
	if err != ErrInvalidConfig {
		t.Errorf("Expected ErrInvalidConfig for nil config, got %v", err)
	}

	// 测试空地址
	err = LoadConfig(&RedisConfig{})
	if err == nil {
		t.Error("Expected error for empty addr")
	}
}

func TestGetClientBeforeInit(t *testing.T) {
	// 创建新的池实例测试
	pool := NewPool()
	client := pool.GetClient()
	if client != nil {
		t.Error("GetClient should return nil before initialization")
	}
}

func TestIsHealthyBeforeInit(t *testing.T) {
	pool := NewPool()
	if pool.IsHealthy() {
		t.Error("Pool should not be healthy before initialization")
	}
}

func TestGetPoolStatsBeforeInit(t *testing.T) {
	pool := NewPool()
	stats := pool.GetPoolStats()
	if stats != nil {
		t.Error("GetPoolStats should return nil before initialization")
	}
}

func TestIsClosed(t *testing.T) {
	pool := NewPool()
	if !pool.IsClosed() {
		t.Error("New pool should be closed")
	}

	// 尝试加载配置（可能失败）
	config := &RedisConfig{
		Addr:        "localhost:6379",
		PoolSize:    5,
		DialTimeout: 1 * time.Second,
	}

	err := pool.LoadConfig(config)
	if err != nil {
		t.Logf("LoadConfig failed (expected if no Redis): %v", err)
		// 即使失败，池也应该是关闭状态
		if !pool.IsClosed() {
			t.Error("Pool should be closed after failed LoadConfig")
		}
	}
}

func TestExecuteWithoutClient(t *testing.T) {
	// 创建一个未初始化的池
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	_, err := Execute("PING")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}
