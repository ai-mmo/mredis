package mredis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// setupMiniRedis 创建测试用的Redis实例
func setupMiniRedis(t *testing.T) (*miniredis.Miniredis, *Pool) {
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

// TestHashOperationsWithRealClient 测试Hash操作
func TestHashOperationsWithRealClient(t *testing.T) {
	mr, _ := setupMiniRedis(t)
	cmd := &RedisCommandHash{}

	t.Run("GetMultiHashValues", func(t *testing.T) {
		mr.HSet("test:hash", "1", "value1")
		mr.HSet("test:hash", "2", "value2")

		values, err := cmd.GetMultiHashValues("test:hash", []uint64{1, 2, 3})
		if err != nil {
			t.Fatalf("GetMultiHashValues failed: %v", err)
		}
		// 函数返回所有非nil值，包括空字符串
		if len(values) < 2 {
			t.Errorf("Expected at least 2 values, got %d", len(values))
		}
	})

	t.Run("GetMultiHashStringValues", func(t *testing.T) {
		mr.HSet("test:hash2", "field1", "value1")
		mr.HSet("test:hash2", "field2", "value2")

		values, err := cmd.GetMultiHashStringValues("test:hash2", []any{"field1", "field2", "field3"})
		if err != nil {
			t.Fatalf("GetMultiHashStringValues failed: %v", err)
		}
		if len(values) != 2 {
			t.Errorf("Expected 2 values, got %d", len(values))
		}
	})

	t.Run("GetAllHashValues", func(t *testing.T) {
		mr.HSet("test:hash3", "f1", "v1")
		mr.HSet("test:hash3", "f2", "v2")

		values, err := cmd.GetAllHashValues("test:hash3")
		if err != nil {
			t.Fatalf("GetAllHashValues failed: %v", err)
		}
		// 函数返回所有字段和值
		if len(values) < 2 {
			t.Errorf("Expected at least 2 values, got %d", len(values))
		}
	})

	t.Run("GetAllBytesHashValues", func(t *testing.T) {
		mr.HSet("test:hash4", "f1", "v1")

		values, err := cmd.GetAllBytesHashValues("test:hash4")
		if err != nil {
			t.Fatalf("GetAllBytesHashValues failed: %v", err)
		}
		if len(values) == 0 {
			t.Error("Expected non-empty values")
		}
	})

	t.Run("IncHashValue", func(t *testing.T) {
		mr.HSet("test:inc", "counter", "10")

		newVal, err := cmd.IncHashValue("test:inc", "counter", int64(5))
		if err != nil {
			t.Fatalf("IncHashValue failed: %v", err)
		}
		if newVal != 15 {
			t.Errorf("Expected 15, got %d", newVal)
		}
	})

	t.Run("GetHashValueString", func(t *testing.T) {
		mr.HSet("test:str", "field", "value")

		val, err := cmd.GetHashValueString("test:str", "field")
		if err != nil {
			t.Fatalf("GetHashValueString failed: %v", err)
		}
		if val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("GetHashValueByte", func(t *testing.T) {
		mr.HSet("test:byte", "field", "data")

		val, err := cmd.GetHashValueByte("test:byte", "field")
		if err != nil {
			t.Fatalf("GetHashValueByte failed: %v", err)
		}
		if string(val) != "data" {
			t.Errorf("Expected 'data', got '%s'", string(val))
		}
	})

	t.Run("GetHashValueInt", func(t *testing.T) {
		mr.HSet("test:int", "num", "42")

		val, err := cmd.GetHashValueInt("test:int", "num")
		if err != nil {
			t.Fatalf("GetHashValueInt failed: %v", err)
		}
		if val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	})
}

// TestListOperationsWithRealClient 测试List操作
func TestListOperationsWithRealClient(t *testing.T) {
	mr, _ := setupMiniRedis(t)
	cmd := &RedisCommandList{}

	t.Run("BLPop", func(t *testing.T) {
		mr.Lpush("test:list", "value")

		val, err := cmd.BLPop("test:list")
		if err != nil {
			t.Fatalf("BLPop failed: %v", err)
		}
		if val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("BRPop", func(t *testing.T) {
		mr.Lpush("test:list2", "value")

		val, err := cmd.BRPop("test:list2")
		if err != nil {
			t.Fatalf("BRPop failed: %v", err)
		}
		if val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("LPop", func(t *testing.T) {
		mr.Lpush("test:list3", "value")

		val, err := cmd.LPop("test:list3")
		if err != nil {
			t.Fatalf("LPop failed: %v", err)
		}
		if val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("RPop", func(t *testing.T) {
		mr.Lpush("test:list4", "value")

		val, err := cmd.RPop("test:list4")
		if err != nil {
			t.Fatalf("RPop failed: %v", err)
		}
		if val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("LIndex", func(t *testing.T) {
		mr.Lpush("test:list5", "value1")
		mr.Lpush("test:list5", "value2")

		val, err := cmd.LIndex("test:list5", 0)
		if err != nil {
			t.Fatalf("LIndex failed: %v", err)
		}
		if val != "value2" {
			t.Errorf("Expected 'value2', got '%s'", val)
		}
	})

	t.Run("LRangeBytes", func(t *testing.T) {
		mr.Lpush("test:list6", "v1")
		mr.Lpush("test:list6", "v2")
		mr.Lpush("test:list6", "v3")

		vals, err := cmd.LRangeBytes("test:list6", 0, -1)
		if err != nil {
			t.Fatalf("LRangeBytes failed: %v", err)
		}
		if len(vals) != 3 {
			t.Errorf("Expected 3 values, got %d", len(vals))
		}
	})

	t.Run("LSet", func(t *testing.T) {
		mr.Lpush("test:list7", "old")

		err := cmd.LSet("test:list7", 0, "new")
		if err != nil {
			t.Fatalf("LSet failed: %v", err)
		}

		// 使用LIndex验证
		val, _ := cmd.LIndex("test:list7", 0)
		if val != "new" {
			t.Errorf("Expected 'new', got '%s'", val)
		}
	})
}

// TestStringOperationsWithRealClient 测试String操作
func TestStringOperationsWithRealClient(t *testing.T) {
	mr, _ := setupMiniRedis(t)
	cmd := &RedisCommandString{}

	t.Run("GetStringValue", func(t *testing.T) {
		mr.Set("test:str", "value")

		val, err := cmd.GetStringValue("test:str")
		if err != nil {
			t.Fatalf("GetStringValue failed: %v", err)
		}
		if val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("GetByteValue", func(t *testing.T) {
		mr.Set("test:byte", "data")

		val, err := cmd.GetByteValue("test:byte")
		if err != nil {
			t.Fatalf("GetByteValue failed: %v", err)
		}
		if string(val) != "data" {
			t.Errorf("Expected 'data', got '%s'", string(val))
		}
	})

	t.Run("GetLock", func(t *testing.T) {
		success, firstTime, err := cmd.GetLock("test:lock", 10, 100, 3)
		if err != nil {
			t.Fatalf("GetLock failed: %v", err)
		}
		if !success {
			t.Error("Expected to get lock")
		}
		if !firstTime {
			t.Error("Expected firstTime to be true")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		mr.Set("test:exists", "value")

		exists, err := cmd.Exists("test:exists")
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if exists == 0 {
			t.Error("Expected key to exist")
		}
	})
}

// TestWrapperExecute 测试Execute函数
func TestWrapperExecute(t *testing.T) {
	mr, _ := setupMiniRedis(t)

	t.Run("Execute SET", func(t *testing.T) {
		result, err := Execute("SET", "test:key", "value")
		if err != nil {
			t.Fatalf("Execute SET failed: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result")
		}
	})

	t.Run("Execute GET", func(t *testing.T) {
		mr.Set("test:get", "value")

		result, err := Execute("GET", "test:get")
		if err != nil {
			t.Fatalf("Execute GET failed: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result")
		}
	})

	t.Run("Execute DEL", func(t *testing.T) {
		mr.Set("test:del", "value")

		result, err := Execute("DEL", "test:del")
		if err != nil {
			t.Fatalf("Execute DEL failed: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result")
		}
	})
}
