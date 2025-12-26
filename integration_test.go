package mredis

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

// TestUtilsSplitFromScore 测试分数拆分工具函数
func TestUtilsSplitFromScore(t *testing.T) {
	tests := []struct {
		name   string
		score  uint64
		first  uint64
		second uint64
	}{
		{
			name:   "zero",
			score:  0,
			first:  0,
			second: 0xFFFFFFFF,
		},
		{
			name:   "max",
			score:  0xFFFFFFFFFFFFFFFF,
			first:  0xFFFFFFFF,
			second: 0,
		},
		{
			name:   "mixed",
			score:  0x123456789ABCDEF0,
			first:  0x12345678,
			second: 0x6543210F,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first, second := SplitFromScore(tt.score)
			if first != tt.first {
				t.Errorf("first = %x, want %x", first, tt.first)
			}
			if second != tt.second {
				t.Errorf("second = %x, want %x", second, tt.second)
			}
		})
	}
}

// TestWrapperFunctions 测试包装器函数
func TestWrapperFunctions(t *testing.T) {
	// 测试未初始化状态
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	t.Run("Close", func(t *testing.T) {
		// 关闭未初始化的池不应该panic
		Close()
	})

	t.Run("IsHealthy", func(t *testing.T) {
		healthy := IsHealthy()
		if healthy {
			t.Error("Uninitialized pool should not be healthy")
		}
	})

	t.Run("GetPoolStats", func(t *testing.T) {
		stats := GetPoolStats()
		if stats != nil {
			t.Error("Uninitialized pool should return nil stats")
		}
	})
}

// TestHashCommandAdditionalMethods 测试Hash命令的额外方法
func TestHashCommandAdditionalMethods(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandHash{}

	t.Run("GetHashValueProtoMsg", func(t *testing.T) {
		err := cmd.GetHashValueProtoMsg("hash:test", "field1", nil)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("SetMultiHashValues", func(t *testing.T) {
		dataMap := make(map[uint64]*proto.Message)
		success := cmd.SetMultiHashValues("hash:test", dataMap)
		if success {
			t.Error("SetMultiHashValues should return false when client is nil")
		}
	})

	t.Run("GetAllBytesHashValues", func(t *testing.T) {
		_, err := cmd.GetAllBytesHashValues("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetAllHashStringValues", func(t *testing.T) {
		_, err := cmd.GetAllHashStringValues("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetMultiHashValues", func(t *testing.T) {
		_, err := cmd.GetMultiHashValues("hash:test", []uint64{1, 2, 3})
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetMultiHashStringValues", func(t *testing.T) {
		_, err := cmd.GetMultiHashStringValues("hash:test", []any{"field1", "field2"})
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashValueIntExceptNil", func(t *testing.T) {
		_, err := cmd.GetHashValueIntExceptNil("hash:test", "field1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})
}

// TestPoolClose 测试连接池关闭
func TestPoolClose(t *testing.T) {
	pool := NewPool()

	// 关闭未初始化的池
	pool.Close()

	if !pool.IsClosed() {
		t.Error("Pool should be closed after Close()")
	}
}

// TestHashCommandSetHashValueProtoMsgTypes 测试不同类型的值设置
func TestHashCommandSetHashValueProtoMsgTypes(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandHash{}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"bytes", []byte("test")},
		{"string", "test"},
		{"int", 123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.SetHashValueProtoMsg("hash:test", "field", tt.value)
			if err != ErrClientNil {
				t.Errorf("Expected ErrClientNil, got %v", err)
			}
		})
	}
}

// TestStringCommandGetLockRetry 测试分布式锁重试逻辑
func TestStringCommandGetLockRetry(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandString{}

	// 测试超过最大重试次数
	success, firstTime, err := cmd.GetLock("lock:test", 1, 1, 1)
	if success {
		t.Error("GetLock should fail when client is nil")
	}
	if firstTime {
		t.Error("firstTime should be false when lock fails")
	}
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// TestZSetCommandGetZSetRankListEdgeCases 测试排行榜边界情况
func TestZSetCommandGetZSetRankListEdgeCases(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandZSet{}

	// 测试相等的排名范围
	_, _, _, err := cmd.GetZSetRankList("zset:test", 5, 5)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// TestHashCommandIncHashValueEdgeCases 测试自增的边界情况
func TestHashCommandIncHashValueEdgeCases(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandHash{}

	// 测试负数增量
	_, err := cmd.IncHashValue("hash:test", "field", int(-10))
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	// 测试大数值
	_, err = cmd.IncHashValue("hash:test", "field", int64(9223372036854775807))
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// TestSetValueKeyLockEdgeCases 测试键锁的边界情况
func TestSetValueKeyLockEdgeCases(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandHash{}

	// 测试零过期时间
	_, err := cmd.SetValueKeyLock("lock:test", int(0))
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	// 测试大过期时间
	_, err = cmd.SetValueKeyLock("lock:test", int64(86400))
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// TestExecuteWithArgs 测试Execute函数的参数处理
func TestExecuteWithArgs(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	// 测试带多个参数的命令
	_, err := Execute("SET", "key", "value", "EX", "60")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	// 测试不带参数的命令
	_, err = Execute("PING")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// TestGetRedisKeyTypeEdgeCases 测试键类型获取的边界情况
func TestGetRedisKeyTypeEdgeCases(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	// 测试空键
	_, err := GetRedisKeyType("")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	// 测试特殊字符键
	_, err = GetRedisKeyType("key:with:colons")
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// TestScanRedisKeysCallback 测试SCAN回调函数
func TestScanRedisKeysCallback(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	callCount := 0
	err := scanRedisKeys("test*", 100, func(keys []string) bool {
		callCount++
		return callCount < 3 // 只处理前3批
	})

	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	if callCount != 0 {
		t.Errorf("Callback should not be called when client is nil, got %d calls", callCount)
	}
}

// TestListCommandEdgeCases 测试List命令的边界情况
func TestListCommandEdgeCases(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandList{}

	// 测试空列表操作
	_, err := cmd.LRange("list:empty", 0, -1)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	// 测试负索引
	_, err = cmd.LIndex("list:test", -1)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	// 测试大范围
	_, err = cmd.LRange("list:test", 0, 10000)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// TestZSetCommandScoreOrdering 测试ZSet分数排序
func TestZSetCommandScoreOrdering(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandZSet{}

	// 测试相同分数
	_, _, err := cmd.GetZSetRangeByScoreWithScoresLimit("zset:test", 100.0, 100.0)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	// 测试负分数
	_, _, err = cmd.GetZSetRangeByScoreWithScoresLimit("zset:test", -100.0, 0.0)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}
