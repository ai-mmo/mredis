package mredis

import (
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	pool := NewPool()
	if pool == nil {
		t.Fatal("NewPool returned nil")
	}

	if !pool.IsClosed() {
		t.Error("New pool should be in closed state")
	}
}

func TestPoolLoadConfig(t *testing.T) {
	pool := NewPool()

	// 测试nil配置
	err := pool.LoadConfig(nil)
	if err != ErrInvalidConfig {
		t.Errorf("Expected ErrInvalidConfig, got %v", err)
	}

	// 测试空地址配置
	err = pool.LoadConfig(&RedisConfig{})
	if err == nil {
		t.Error("Expected error for empty addr, got nil")
	}

	// 测试有效配置（但可能连接失败）
	config := &RedisConfig{
		Addr:         "localhost:6379",
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  5 * time.Second,
	}

	err = pool.LoadConfig(config)
	// 注意：这里可能会失败如果没有Redis服务器运行
	// 在实际测试中，应该使用mock或者跳过这个测试
	if err != nil {
		t.Logf("LoadConfig failed (expected if no Redis server): %v", err)
	}
}

func TestPoolGetClient(t *testing.T) {
	pool := NewPool()

	// 测试未初始化的池
	client := pool.GetClient()
	if client != nil {
		t.Error("GetClient should return nil for uninitialized pool")
	}
}

func TestVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("Version should not be empty")
	}
	t.Logf("Current version: %s", version)
}

func TestErrors(t *testing.T) {
	// 测试错误常量是否定义
	if ErrClientNil == nil {
		t.Error("ErrClientNil should not be nil")
	}
	if ErrPoolClosed == nil {
		t.Error("ErrPoolClosed should not be nil")
	}
	if ErrInvalidConfig == nil {
		t.Error("ErrInvalidConfig should not be nil")
	}
}

func TestContextWithTimeout(t *testing.T) {
	ctx, cancel := ContextWithTimeout(1 * time.Second)
	defer cancel()

	if ctx == nil {
		t.Error("Context should not be nil")
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Context should have a deadline")
	}

	if time.Until(deadline) > 2*time.Second {
		t.Error("Deadline is too far in the future")
	}
}

func TestContextWithDefaultTimeout(t *testing.T) {
	ctx, cancel := ContextWithDefaultTimeout()
	defer cancel()

	if ctx == nil {
		t.Error("Context should not be nil")
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Context should have a deadline")
	}

	expectedTimeout := DefaultTimeout
	actualTimeout := time.Until(deadline)

	// 允许一些误差
	if actualTimeout > expectedTimeout+time.Second || actualTimeout < expectedTimeout-time.Second {
		t.Errorf("Expected timeout around %v, got %v", expectedTimeout, actualTimeout)
	}
}
