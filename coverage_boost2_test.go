package mredis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// setupTestRedis2 创建测试用的Redis实例
func setupTestRedis2(t *testing.T) (*miniredis.Miniredis, *Pool) {
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

// TestMoreHashOperations 测试更多Hash操作
func TestMoreHashOperations(t *testing.T) {
	mr, _ := setupTestRedis2(t)
	cmd := &RedisCommandHash{}

	t.Run("Exists", func(t *testing.T) {
		mr.HSet("test:exists", "field", "value")

		exists := cmd.Exists("test:exists", "field")
		if !exists {
			t.Error("Expected field to exist")
		}
	})

	t.Run("DeleteHashKeyList", func(t *testing.T) {
		mr.HSet("test:del", "1", "v1")
		mr.HSet("test:del", "2", "v2")

		err := cmd.DeleteHashKeyList("test:del", []uint64{1, 2})
		if err != nil {
			t.Fatalf("DeleteHashKeyList failed: %v", err)
		}
	})

	t.Run("GetHashLen", func(t *testing.T) {
		mr.HSet("test:len", "f1", "v1")
		mr.HSet("test:len", "f2", "v2")

		length, err := cmd.GetHashLen("test:len")
		if err != nil {
			t.Fatalf("GetHashLen failed: %v", err)
		}
		if length != 2 {
			t.Errorf("Expected length 2, got %d", length)
		}
	})

	t.Run("GetHashKeys", func(t *testing.T) {
		mr.HSet("test:keys", "f1", "v1")
		mr.HSet("test:keys", "f2", "v2")

		keys, err := cmd.GetHashKeys("test:keys")
		if err != nil {
			t.Fatalf("GetHashKeys failed: %v", err)
		}
		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(keys))
		}
	})

	t.Run("SetHashExpire", func(t *testing.T) {
		mr.HSet("test:expire", "f1", "v1")

		err := cmd.SetHashExpire("test:expire", 60)
		if err != nil {
			t.Fatalf("SetHashExpire failed: %v", err)
		}
	})

	t.Run("GetHashSize", func(t *testing.T) {
		mr.HSet("test:size", "f1", "v1")

		size, err := cmd.GetHashSize("test:size")
		if err != nil {
			t.Fatalf("GetHashSize failed: %v", err)
		}
		if size == 0 {
			t.Error("Expected non-zero size")
		}
	})

	t.Run("HashExists", func(t *testing.T) {
		mr.HSet("test:hexists", "f1", "v1")

		exists, err := cmd.HashExists("test:hexists")
		if err != nil {
			t.Fatalf("HashExists failed: %v", err)
		}
		if !exists {
			t.Error("Expected hash to exist")
		}
	})
}

// TestMoreListOperations 测试更多List操作
func TestMoreListOperations(t *testing.T) {
	mr, _ := setupTestRedis2(t)
	cmd := &RedisCommandList{}

	t.Run("LPush", func(t *testing.T) {
		err := cmd.LPush("test:lpush", "v1", "v2")
		if err != nil {
			t.Fatalf("LPush failed: %v", err)
		}

		// 验证长度
		length, _ := cmd.LLen("test:lpush")
		if length != 2 {
			t.Errorf("Expected length 2, got %d", length)
		}
	})

	t.Run("RPush", func(t *testing.T) {
		err := cmd.RPush("test:rpush", "v1", "v2")
		if err != nil {
			t.Fatalf("RPush failed: %v", err)
		}
	})

	t.Run("LLen", func(t *testing.T) {
		mr.Lpush("test:llen", "v1")
		mr.Lpush("test:llen", "v2")

		length, err := cmd.LLen("test:llen")
		if err != nil {
			t.Fatalf("LLen failed: %v", err)
		}
		if length != 2 {
			t.Errorf("Expected length 2, got %d", length)
		}
	})

	t.Run("LRange", func(t *testing.T) {
		mr.Lpush("test:lrange", "v1")
		mr.Lpush("test:lrange", "v2")

		values, err := cmd.LRange("test:lrange", 0, -1)
		if err != nil {
			t.Fatalf("LRange failed: %v", err)
		}
		if len(values) != 2 {
			t.Errorf("Expected 2 values, got %d", len(values))
		}
	})

	t.Run("LTrim", func(t *testing.T) {
		mr.Lpush("test:ltrim", "v1")
		mr.Lpush("test:ltrim", "v2")
		mr.Lpush("test:ltrim", "v3")

		err := cmd.LTrim("test:ltrim", 0, 1)
		if err != nil {
			t.Fatalf("LTrim failed: %v", err)
		}
	})

	t.Run("LRem", func(t *testing.T) {
		mr.Lpush("test:lrem", "v1")
		mr.Lpush("test:lrem", "v2")
		mr.Lpush("test:lrem", "v1")

		count, err := cmd.LRem("test:lrem", 0, "v1")
		if err != nil {
			t.Fatalf("LRem failed: %v", err)
		}
		if count != 2 {
			t.Errorf("Expected 2 removed, got %d", count)
		}
	})
}

// TestMoreStringOperations 测试更多String操作
func TestMoreStringOperations(t *testing.T) {
	mr, _ := setupTestRedis2(t)
	cmd := &RedisCommandString{}

	t.Run("SetStringValue", func(t *testing.T) {
		err := cmd.SetStringValue("test:set", "value")
		if err != nil {
			t.Fatalf("SetStringValue failed: %v", err)
		}

		val, _ := mr.Get("test:set")
		if val != "value" {
			t.Errorf("Expected 'value', got '%s'", val)
		}
	})

	t.Run("DelStringKey", func(t *testing.T) {
		mr.Set("test:delstr", "value")

		err := cmd.DelStringKey("test:delstr")
		if err != nil {
			t.Fatalf("DelStringKey failed: %v", err)
		}
	})

	t.Run("GetIntValue", func(t *testing.T) {
		mr.Set("test:int", "42")

		val, err := cmd.GetIntValue("test:int")
		if err != nil {
			t.Fatalf("GetIntValue failed: %v", err)
		}
		if val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	})

	t.Run("GetInt64Value", func(t *testing.T) {
		mr.Set("test:int64", "9223372036854775807")

		val, err := cmd.GetInt64Value("test:int64")
		if err != nil {
			t.Fatalf("GetInt64Value failed: %v", err)
		}
		if val != 9223372036854775807 {
			t.Errorf("Expected max int64, got %d", val)
		}
	})

	t.Run("Incr", func(t *testing.T) {
		mr.Set("test:incr", "10")

		newVal, err := cmd.Incr("test:incr")
		if err != nil {
			t.Fatalf("Incr failed: %v", err)
		}
		if newVal != 11 {
			t.Errorf("Expected 11, got %d", newVal)
		}
	})
}

// TestMoreZSetOperations 测试更多ZSet操作
func TestMoreZSetOperations(t *testing.T) {
	mr, _ := setupTestRedis2(t)
	cmd := &RedisCommandZSet{}

	t.Run("GetZSetCount", func(t *testing.T) {
		mr.ZAdd("test:zcount", 1.0, "m1")
		mr.ZAdd("test:zcount", 2.0, "m2")

		count, err := cmd.GetZSetCount("test:zcount")
		if err != nil {
			t.Fatalf("GetZSetCount failed: %v", err)
		}
		if count != 2 {
			t.Errorf("Expected count 2, got %d", count)
		}
	})

	t.Run("ZScore", func(t *testing.T) {
		mr.ZAdd("test:zscore", 100.0, "member")

		score, err := cmd.ZScore("test:zscore", "member")
		if err != nil {
			t.Fatalf("ZScore failed: %v", err)
		}
		if score != 100.0 {
			t.Errorf("Expected score 100.0, got %f", score)
		}
	})

	t.Run("ZRank", func(t *testing.T) {
		mr.ZAdd("test:zrank", 1.0, "m1")
		mr.ZAdd("test:zrank", 2.0, "m2")
		mr.ZAdd("test:zrank", 3.0, "m3")

		rank, err := cmd.ZRank("test:zrank", "m2")
		if err != nil {
			t.Fatalf("ZRank failed: %v", err)
		}
		if rank != 1 {
			t.Errorf("Expected rank 1, got %d", rank)
		}
	})

	t.Run("ZRevRank", func(t *testing.T) {
		mr.ZAdd("test:zrevrank", 1.0, "m1")
		mr.ZAdd("test:zrevrank", 2.0, "m2")
		mr.ZAdd("test:zrevrank", 3.0, "m3")

		rank, err := cmd.ZRevRank("test:zrevrank", "m2")
		if err != nil {
			t.Fatalf("ZRevRank failed: %v", err)
		}
		if rank != 1 {
			t.Errorf("Expected rank 1, got %d", rank)
		}
	})

	t.Run("ZIncrBy", func(t *testing.T) {
		mr.ZAdd("test:zincrby", 10.0, "member")

		newScore, err := cmd.ZIncrBy("test:zincrby", 5.0, "member")
		if err != nil {
			t.Fatalf("ZIncrBy failed: %v", err)
		}
		if newScore != 15.0 {
			t.Errorf("Expected score 15.0, got %f", newScore)
		}
	})

	t.Run("ZSetAddValue", func(t *testing.T) {
		err := cmd.ZSetAddValue("test:zsetadd", "member", 100)
		if err != nil {
			t.Fatalf("ZSetAddValue failed: %v", err)
		}
	})

	t.Run("GetZSetRankByScore", func(t *testing.T) {
		mr.ZAdd("test:zrankbyscore", 10.0, "m1")
		mr.ZAdd("test:zrankbyscore", 20.0, "m2")
		mr.ZAdd("test:zrankbyscore", 30.0, "m3")

		count, err := cmd.GetZSetRankByScore("test:zrankbyscore", 20.0)
		if err != nil {
			t.Fatalf("GetZSetRankByScore failed: %v", err)
		}
		if count < 1 {
			t.Errorf("Expected at least 1, got %d", count)
		}
	})

	t.Run("GetZSetRankAndScore", func(t *testing.T) {
		mr.ZAdd("test:zrankandscore", 100.0, "member")

		rank, score, scoreUint, err := cmd.GetZSetRankAndScore("test:zrankandscore", "member")
		if err != nil {
			t.Fatalf("GetZSetRankAndScore failed: %v", err)
		}
		if rank != 1 {
			t.Errorf("Expected rank 1, got %d", rank)
		}
		if score != 100.0 {
			t.Errorf("Expected score 100.0, got %f", score)
		}
		if scoreUint != 100 {
			t.Errorf("Expected scoreUint 100, got %d", scoreUint)
		}
	})
}
