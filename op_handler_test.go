package mredis

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestSetRedisKeyWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	err := SetRedisKey("test", "value")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

func TestGetRedisKeyTypeWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	_, err := GetRedisKeyType("test")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

func TestDeleteRedisKeyWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	err := DeleteRedisKey("test")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

func TestGetRedisKeysByPatternWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	_, err := GetRedisKeysByPattern("test*")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

func TestDeleteRedisKeysByPatternWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	_, err := DeleteRedisKeysByPattern("test*")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

func TestScanRedisKeysWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	called := false
	err := scanRedisKeys("test*", 100, func(keys []string) bool {
		called = true
		return true
	})

	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	if called {
		t.Error("Callback should not be called when client is nil")
	}
}

func TestScanRedisKeysWithDefaultBatchSize(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	// 测试负数批次大小会使用默认值
	err := scanRedisKeys("test*", -1, func(keys []string) bool {
		return true
	})

	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// setupTestRedisForHandler 创建测试用的Redis实例
func setupTestRedisForHandler(t *testing.T) (*miniredis.Miniredis, *Pool) {
	mr := miniredis.RunT(t)

	pool := NewPool()
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	pool.client = client
	pool.closed = false

	oldPool := gPoolM
	gPoolM = pool

	t.Cleanup(func() {
		gPoolM = oldPool
		pool.Close()
		mr.Close()
	})

	return mr, pool
}

func TestSetRedisKeyWithRealClient(t *testing.T) {
	mr, _ := setupTestRedisForHandler(t)

	// 测试设置键
	err := SetRedisKey("test:key", "test:value")
	if err != nil {
		t.Fatalf("SetRedisKey failed: %v", err)
	}

	// 验证键已设置
	val, err := mr.Get("test:key")
	if err != nil {
		t.Fatalf("Get key failed: %v", err)
	}
	if val != "test:value" {
		t.Errorf("Expected value 'test:value', got '%s'", val)
	}
}

func TestGetRedisKeyTypeWithRealClient(t *testing.T) {
	mr, _ := setupTestRedisForHandler(t)

	// 设置不同类型的键
	mr.Set("string:key", "value")
	mr.HSet("hash:key", "field", "value")
	mr.Lpush("list:key", "value")
	mr.ZAdd("zset:key", 1.0, "member")
	mr.SAdd("set:key", "member")

	tests := []struct {
		key      string
		expected string
	}{
		{"string:key", "string"},
		{"hash:key", "hash"},
		{"list:key", "list"},
		{"zset:key", "zset"},
		{"set:key", "set"},
		{"nonexistent", "none"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			keyType, err := GetRedisKeyType(tt.key)
			if err != nil {
				t.Fatalf("GetRedisKeyType failed: %v", err)
			}
			if keyType != tt.expected {
				t.Errorf("Expected type '%s', got '%s'", tt.expected, keyType)
			}
		})
	}
}

func TestDeleteRedisKeyWithRealClient(t *testing.T) {
	mr, _ := setupTestRedisForHandler(t)

	// 设置测试键
	mr.Set("test:delete", "value")

	// 删除键
	err := DeleteRedisKey("test:delete")
	if err != nil {
		t.Fatalf("DeleteRedisKey failed: %v", err)
	}

	// 验证键已删除
	exists := mr.Exists("test:delete")
	if exists {
		t.Error("Key should be deleted")
	}
}

func TestGetRedisKeysByPatternWithRealClient(t *testing.T) {
	mr, _ := setupTestRedisForHandler(t)

	// 设置多个匹配的键
	mr.Set("test:pattern:1", "value1")
	mr.Set("test:pattern:2", "value2")
	mr.Set("test:pattern:3", "value3")
	mr.Set("other:key", "value")

	// 获取匹配的键
	keys, err := GetRedisKeysByPattern("test:pattern:*")
	if err != nil {
		t.Fatalf("GetRedisKeysByPattern failed: %v", err)
	}

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// 验证键的内容
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	expectedKeys := []string{"test:pattern:1", "test:pattern:2", "test:pattern:3"}
	for _, expected := range expectedKeys {
		if !keyMap[expected] {
			t.Errorf("Expected key '%s' not found", expected)
		}
	}
}

func TestDeleteRedisKeysByPatternWithRealClient(t *testing.T) {
	mr, _ := setupTestRedisForHandler(t)

	// 设置多个匹配的键
	mr.Set("test:delete:1", "value1")
	mr.Set("test:delete:2", "value2")
	mr.Set("test:delete:3", "value3")
	mr.Set("keep:key", "value")

	// 删除匹配的键
	deletedKeys, err := DeleteRedisKeysByPattern("test:delete:*")
	if err != nil {
		t.Fatalf("DeleteRedisKeysByPattern failed: %v", err)
	}

	if len(deletedKeys) != 3 {
		t.Errorf("Expected 3 deleted keys, got %d", len(deletedKeys))
	}

	// 验证键已删除
	if mr.Exists("test:delete:1") || mr.Exists("test:delete:2") || mr.Exists("test:delete:3") {
		t.Error("Matched keys should be deleted")
	}

	// 验证其他键未删除
	if !mr.Exists("keep:key") {
		t.Error("Non-matched key should not be deleted")
	}
}

func TestScanRedisKeysWithRealClient(t *testing.T) {
	mr, _ := setupTestRedisForHandler(t)

	// 设置多个键
	for i := 0; i < 10; i++ {
		mr.Set(fmt.Sprintf("scan:test:%d", i), "value")
	}

	// 扫描键
	var scannedKeys []string
	err := scanRedisKeys("scan:test:*", 3, func(keys []string) bool {
		scannedKeys = append(scannedKeys, keys...)
		return true
	})

	if err != nil {
		t.Fatalf("scanRedisKeys failed: %v", err)
	}

	if len(scannedKeys) != 10 {
		t.Errorf("Expected 10 scanned keys, got %d", len(scannedKeys))
	}
}

func TestScanRedisKeysWithStopCallback(t *testing.T) {
	mr, _ := setupTestRedisForHandler(t)

	// 设置多个键
	for i := 0; i < 20; i++ {
		mr.Set(fmt.Sprintf("stop:test:%d", i), "value")
	}

	// 扫描键，但在第一批后停止
	callCount := 0
	err := scanRedisKeys("stop:test:*", 5, func(keys []string) bool {
		callCount++
		return callCount < 2 // 只处理前两批
	})

	if err != nil {
		t.Fatalf("scanRedisKeys failed: %v", err)
	}

	if callCount >= 4 {
		t.Errorf("Expected callback to stop early, but was called %d times", callCount)
	}
}
