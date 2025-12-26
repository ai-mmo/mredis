package mredis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestZSetCommandWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandZSet{}

	t.Run("ZAdd", func(t *testing.T) {
		success := cmd.ZAdd("zset:test", 100.0, "member1")
		if success {
			t.Error("ZAdd should return false when client is nil")
		}
	})

	t.Run("GetZSetCount", func(t *testing.T) {
		_, err := cmd.GetZSetCount("zset:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRangeByScore", func(t *testing.T) {
		_, err := cmd.GetZSetRangeByScore("zset:test", "0", "100")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRankList", func(t *testing.T) {
		_, _, _, err := cmd.GetZSetRankList("zset:test", 1, 10)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRangeByScoreWithScores", func(t *testing.T) {
		_, _, err := cmd.GetZSetRangeByScoreWithScores("zset:test", "0", "100", nil)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRangeByScoreWithScoresLimit", func(t *testing.T) {
		_, _, err := cmd.GetZSetRangeByScoreWithScoresLimit("zset:test", 0.0, 100.0)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ZRem", func(t *testing.T) {
		err := cmd.ZRem("zset:test", "member1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ZRemRangeByScore", func(t *testing.T) {
		err := cmd.ZRemRangeByScore("zset:test", "0", "100")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ZScore", func(t *testing.T) {
		_, err := cmd.ZScore("zset:test", "member1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ZRank", func(t *testing.T) {
		_, err := cmd.ZRank("zset:test", "member1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ZRevRank", func(t *testing.T) {
		_, err := cmd.ZRevRank("zset:test", "member1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ZRange", func(t *testing.T) {
		_, err := cmd.ZRange("zset:test", 0, -1)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ZRevRange", func(t *testing.T) {
		_, err := cmd.ZRevRange("zset:test", 0, -1)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ZIncrBy", func(t *testing.T) {
		_, err := cmd.ZIncrBy("zset:test", 10.0, "member1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("DelKey", func(t *testing.T) {
		err := cmd.DelKey("zset:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})
}

func TestZSetCommandCompatibilityMethods(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandZSet{}

	t.Run("ZSetAddValue", func(t *testing.T) {
		err := cmd.ZSetAddValue("zset:test", "member1", 100)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetValueByRange", func(t *testing.T) {
		_, err := cmd.GetZSetValueByRange("zset:test", 1, 10)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRandomValueByRange", func(t *testing.T) {
		_, _, err := cmd.GetZSetRandomValueByRange("zset:test", 0.0, 100.0, nil, 10)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRangeByValue", func(t *testing.T) {
		_, _, err := cmd.GetZSetRangeByValue("zset:test", 0.0, 100.0, nil, 10)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetValueByRank", func(t *testing.T) {
		_, err := cmd.GetZSetValueByRank("zset:test", 1, 10)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRank", func(t *testing.T) {
		_, err := cmd.GetZSetRank("zset:test", "member1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRankByScore", func(t *testing.T) {
		_, err := cmd.GetZSetRankByScore("zset:test", 100.0)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetZSetRankAndScore", func(t *testing.T) {
		_, _, _, err := cmd.GetZSetRankAndScore("zset:test", "member1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("RemZSetElement", func(t *testing.T) {
		err := cmd.RemZSetElement("zset:test", "member1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("ResetZSet", func(t *testing.T) {
		err := cmd.ResetZSet("zset:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("DelZSetKey", func(t *testing.T) {
		err := cmd.DelZSetKey("zset:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})
}

func TestGetZSetRankListInvalidRange(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandZSet{}

	// 测试 rankMin > rankMax 的情况
	_, _, _, err := cmd.GetZSetRankList("zset:test", 10, 1)
	if err == nil {
		t.Error("Expected error when rankMin > rankMax")
	}
}

func TestGetZSetRangeByScoreWithScoresLimitOrdering(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandZSet{}

	// 测试正序
	_, _, err := cmd.GetZSetRangeByScoreWithScoresLimit("zset:test", 0.0, 100.0)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}

	// 测试逆序
	_, _, err = cmd.GetZSetRangeByScoreWithScoresLimit("zset:test", 100.0, 0.0)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}

// setupTestRedis 创建一个miniredis实例用于测试
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *Pool) {
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

func TestZSetRankListWithRealClient(t *testing.T) {
	_, _ = setupTestRedis(t)

	cmd := &RedisCommandZSet{}
	key := "test:zset:ranklist"

	// 清理测试数据
	_ = cmd.DelKey(key)

	// 添加测试数据
	cmd.ZAdd(key, 100.0, "member1")
	cmd.ZAdd(key, 200.0, "member2")
	cmd.ZAdd(key, 300.0, "member3")
	cmd.ZAdd(key, 400.0, "member4")
	cmd.ZAdd(key, 500.0, "member5")

	// 测试正常范围
	members, ranks, scores, err := cmd.GetZSetRankList(key, 1, 3)
	if err != nil {
		t.Fatalf("GetZSetRankList failed: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(members))
	}
	if len(ranks) != 3 {
		t.Errorf("Expected 3 ranks, got %d", len(ranks))
	}
	if len(scores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(scores))
	}

	// 验证排名顺序（从高到低）
	if members[0] != "member5" {
		t.Errorf("Expected first member to be member5, got %s", members[0])
	}
	if ranks[0] != 1 {
		t.Errorf("Expected first rank to be 1, got %d", ranks[0])
	}

	// 清理
	_ = cmd.DelKey(key)
}

func TestZSetRangeByScoreWithScoresLimitWithRealClient(t *testing.T) {
	_, _ = setupTestRedis(t)

	cmd := &RedisCommandZSet{}
	key := "test:zset:scorelimit"

	// 清理测试数据
	_ = cmd.DelKey(key)

	// 添加测试数据
	cmd.ZAdd(key, 10.0, "member1")
	cmd.ZAdd(key, 20.0, "member2")
	cmd.ZAdd(key, 30.0, "member3")
	cmd.ZAdd(key, 40.0, "member4")
	cmd.ZAdd(key, 50.0, "member5")

	// 测试正序（param1 <= param2）
	members, scores, err := cmd.GetZSetRangeByScoreWithScoresLimit(key, 20.0, 40.0)
	if err != nil {
		t.Fatalf("GetZSetRangeByScoreWithScoresLimit (ascending) failed: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(members))
	}
	if len(scores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(scores))
	}
	// 验证正序
	if members[0] != "member2" {
		t.Errorf("Expected first member to be member2, got %s", members[0])
	}

	// 测试逆序（param1 > param2）
	members, scores, err = cmd.GetZSetRangeByScoreWithScoresLimit(key, 40.0, 20.0)
	if err != nil {
		t.Fatalf("GetZSetRangeByScoreWithScoresLimit (descending) failed: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(members))
	}
	// 验证逆序
	if members[0] != "member4" {
		t.Errorf("Expected first member to be member4, got %s", members[0])
	}

	// 清理
	_ = cmd.DelKey(key)
}

func TestZSetRangeByScoreWithScoresWithRealClient(t *testing.T) {
	_, _ = setupTestRedis(t)

	cmd := &RedisCommandZSet{}
	key := "test:zset:scorewithscores"

	// 清理测试数据
	_ = cmd.DelKey(key)

	// 添加测试数据
	cmd.ZAdd(key, 10.0, "member1")
	cmd.ZAdd(key, 20.0, "member2")
	cmd.ZAdd(key, 30.0, "member3")

	// 测试带matched参数
	matched := make(map[string]bool)
	members, scores, err := cmd.GetZSetRangeByScoreWithScores(key, "10", "30", matched)
	if err != nil {
		t.Fatalf("GetZSetRangeByScoreWithScores failed: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(members))
	}
	if len(scores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(scores))
	}
	if !matched["member1"] || !matched["member2"] || !matched["member3"] {
		t.Error("Expected all members to be marked as matched")
	}

	// 测试空结果
	members, scores, err = cmd.GetZSetRangeByScoreWithScores(key, "100", "200", nil)
	if err != nil {
		t.Fatalf("GetZSetRangeByScoreWithScores (empty) failed: %v", err)
	}
	if members != nil || scores != nil {
		t.Error("Expected nil result for empty range")
	}

	// 清理
	_ = cmd.DelKey(key)
}

func TestZSetRangeByScoreWithRealClient(t *testing.T) {
	_, _ = setupTestRedis(t)

	cmd := &RedisCommandZSet{}
	key := "test:zset:rangebyscore"

	// 清理测试数据
	_ = cmd.DelKey(key)

	// 添加测试数据
	cmd.ZAdd(key, 10.0, "member1")
	cmd.ZAdd(key, 20.0, "member2")
	cmd.ZAdd(key, 30.0, "member3")

	// 测试正常范围
	members, err := cmd.GetZSetRangeByScore(key, "10", "30")
	if err != nil {
		t.Fatalf("GetZSetRangeByScore failed: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(members))
	}

	// 测试空结果
	members, err = cmd.GetZSetRangeByScore(key, "100", "200")
	if err != nil {
		t.Fatalf("GetZSetRangeByScore (empty) failed: %v", err)
	}
	if members != nil {
		t.Error("Expected nil result for empty range")
	}

	// 清理
	_ = cmd.DelKey(key)
}
