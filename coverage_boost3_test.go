package mredis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
)

// setupTestRedis3 设置测试用的Redis环境
func setupTestRedis3(t *testing.T) (*miniredis.Miniredis, func()) {
	mr := miniredis.RunT(t)

	config := &RedisConfig{
		Addr:  mr.Addr(),
		Index: 0,
	}

	if err := gPoolM.LoadConfig(config); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	cleanup := func() {
		gPoolM.Close()
		mr.Close()
	}

	return mr, cleanup
}

// TestHashProtoMsgOperations 测试Hash的Protobuf消息操作
func TestHashProtoMsgOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	hash := &RedisCommandHash{}

	t.Run("SetHashValueProtoMsg_WithBytes", func(t *testing.T) {
		data := []byte("test_bytes")
		err := hash.SetHashValueProtoMsg("proto_hash", "field2", data)
		if err != nil {
			t.Errorf("SetHashValueProtoMsg with bytes failed: %v", err)
		}
	})

	t.Run("SetHashValueProtoMsg_WithString", func(t *testing.T) {
		err := hash.SetHashValueProtoMsg("proto_hash", "field3", "test_string")
		if err != nil {
			t.Errorf("SetHashValueProtoMsg with string failed: %v", err)
		}
	})

	t.Run("SetHashValue_Compatibility", func(t *testing.T) {
		err := hash.SetHashValue("compat_hash", "field1", "value1")
		if err != nil {
			t.Errorf("SetHashValue failed: %v", err)
		}
	})
}

// TestHashAdvancedOperations 测试Hash的高级操作
func TestHashAdvancedOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	hash := &RedisCommandHash{}

	t.Run("GetHashValuesAsJSON_Empty", func(t *testing.T) {
		json, err := hash.GetHashValuesAsJSON("empty_hash")
		if err != nil {
			t.Errorf("GetHashValuesAsJSON failed: %v", err)
		}
		if json != "{}" {
			t.Errorf("Expected empty JSON object, got: %s", json)
		}
	})

	t.Run("GetHashValuesAsJSON_WithData", func(t *testing.T) {
		// 先设置一些数据
		_ = hash.SetHashValueProtoMsg("json_hash", "field1", "value1")
		_ = hash.SetHashValueProtoMsg("json_hash", "field2", "value2")

		json, err := hash.GetHashValuesAsJSON("json_hash")
		if err != nil {
			t.Errorf("GetHashValuesAsJSON failed: %v", err)
		}
		if json == "" {
			t.Error("Expected non-empty JSON")
		}
	})

	t.Run("GetHashValuesWithFilter_NoFilter", func(t *testing.T) {
		_ = hash.SetHashValueProtoMsg("filter_hash", "field1", "value1")

		result, err := hash.GetHashValuesWithFilter("filter_hash", nil)
		if err != nil {
			t.Errorf("GetHashValuesWithFilter failed: %v", err)
		}
		if len(result) == 0 {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("GetHashValuesWithFilter_WithFilter", func(t *testing.T) {
		_ = hash.SetHashValueProtoMsg("filter_hash2", "prefix_field1", "value1")
		_ = hash.SetHashValueProtoMsg("filter_hash2", "other_field", "value2")

		filter := func(field string) bool {
			return len(field) >= 6 && field[:6] == "prefix"
		}

		result, err := hash.GetHashValuesWithFilter("filter_hash2", filter)
		if err != nil {
			t.Errorf("GetHashValuesWithFilter failed: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("Expected 1 filtered result, got: %d", len(result))
		}
	})

	t.Run("GetHashValuesForAPI_NotExists", func(t *testing.T) {
		json, err := hash.GetHashValuesForAPI("nonexistent_hash")
		if err != nil {
			t.Errorf("GetHashValuesForAPI failed: %v", err)
		}
		if json != "{}" {
			t.Errorf("Expected empty JSON, got: %s", json)
		}
	})

	t.Run("GetHashValuesForAPI_Exists", func(t *testing.T) {
		_ = hash.SetHashValueProtoMsg("api_hash", "field1", "value1")

		json, err := hash.GetHashValuesForAPI("api_hash")
		if err != nil {
			t.Errorf("GetHashValuesForAPI failed: %v", err)
		}
		if json == "" {
			t.Error("Expected non-empty JSON")
		}
	})

	t.Run("GetHashValuesWithPrefix", func(t *testing.T) {
		_ = hash.SetHashValueProtoMsg("prefix_hash", "user_name", "alice")
		_ = hash.SetHashValueProtoMsg("prefix_hash", "user_age", "30")
		_ = hash.SetHashValueProtoMsg("prefix_hash", "admin_role", "superuser")

		result, err := hash.GetHashValuesWithPrefix("prefix_hash", "user")
		if err != nil {
			t.Errorf("GetHashValuesWithPrefix failed: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 results with 'user' prefix, got: %d", len(result))
		}
	})

	t.Run("GetHashValuesExcludePrefix", func(t *testing.T) {
		_ = hash.SetHashValueProtoMsg("exclude_hash", "user_name", "alice")
		_ = hash.SetHashValueProtoMsg("exclude_hash", "admin_role", "superuser")

		result, err := hash.GetHashValuesExcludePrefix("exclude_hash", "user")
		if err != nil {
			t.Errorf("GetHashValuesExcludePrefix failed: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("Expected 1 result excluding 'user' prefix, got: %d", len(result))
		}
	})
}

// TestHashLockOperations 测试Hash的锁操作
func TestHashLockOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	hash := &RedisCommandHash{}

	t.Run("SetValueKeyLock_WithUint32", func(t *testing.T) {
		result, err := hash.SetValueKeyLock("lock_key1", uint32(10))
		if err != nil {
			t.Errorf("SetValueKeyLock failed: %v", err)
		}
		if !result {
			t.Error("Expected lock to be acquired")
		}
	})

	t.Run("SetValueKeyLock_WithInt", func(t *testing.T) {
		result, err := hash.SetValueKeyLock("lock_key2", 10)
		if err != nil {
			t.Errorf("SetValueKeyLock failed: %v", err)
		}
		if !result {
			t.Error("Expected lock to be acquired")
		}
	})

	t.Run("SetValueKeyLock_WithInt64", func(t *testing.T) {
		result, err := hash.SetValueKeyLock("lock_key3", int64(10))
		if err != nil {
			t.Errorf("SetValueKeyLock failed: %v", err)
		}
		if !result {
			t.Error("Expected lock to be acquired")
		}
	})

	t.Run("SetValueKeyLock_AlreadyLocked", func(t *testing.T) {
		// 先获取锁
		_, _ = hash.SetValueKeyLock("lock_key4", 10)

		// 再次尝试获取同一个锁
		result, err := hash.SetValueKeyLock("lock_key4", 10)
		if err != nil {
			t.Errorf("SetValueKeyLock failed: %v", err)
		}
		if result {
			t.Error("Expected lock to be already acquired")
		}
	})
}

// TestHashIncOperations 测试Hash的自增操作
func TestHashIncOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	hash := &RedisCommandHash{}

	t.Run("IncHashValue_WithInt32", func(t *testing.T) {
		result, err := hash.IncHashValue("inc_hash", "field1", int32(5))
		if err != nil {
			t.Errorf("IncHashValue failed: %v", err)
		}
		if result != 5 {
			t.Errorf("Expected 5, got: %d", result)
		}
	})

	t.Run("IncHashValue_WithUint32", func(t *testing.T) {
		result, err := hash.IncHashValue("inc_hash", "field2", uint32(10))
		if err != nil {
			t.Errorf("IncHashValue failed: %v", err)
		}
		if result != 10 {
			t.Errorf("Expected 10, got: %d", result)
		}
	})

	t.Run("IncHashValue_Multiple", func(t *testing.T) {
		_, _ = hash.IncHashValue("inc_hash", "field3", 1)
		result, err := hash.IncHashValue("inc_hash", "field3", 2)
		if err != nil {
			t.Errorf("IncHashValue failed: %v", err)
		}
		if result != 3 {
			t.Errorf("Expected 3, got: %d", result)
		}
	})
}

// TestHashDeleteOperations 测试Hash的删除操作
func TestHashDeleteOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	hash := &RedisCommandHash{}

	t.Run("DeleteHashStringKeyList", func(t *testing.T) {
		// 先设置一些数据
		_ = hash.SetHashValueProtoMsg("del_hash", "field1", "value1")
		_ = hash.SetHashValueProtoMsg("del_hash", "field2", "value2")
		_ = hash.SetHashValueProtoMsg("del_hash", "field3", "value3")

		err := hash.DeleteHashStringKeyList("del_hash", []string{"field1", "field2"})
		if err != nil {
			t.Errorf("DeleteHashStringKeyList failed: %v", err)
		}

		// 验证删除
		exists := hash.Exists("del_hash", "field1")
		if exists {
			t.Error("field1 should be deleted")
		}
	})

	t.Run("DelKey", func(t *testing.T) {
		_ = hash.SetHashValueProtoMsg("del_key_hash", "field1", "value1")

		err := hash.DelKey("del_key_hash")
		if err != nil {
			t.Errorf("DelKey failed: %v", err)
		}

		// 验证删除
		exists, _ := hash.HashExists("del_key_hash")
		if exists {
			t.Error("Hash should be deleted")
		}
	})
}

// TestHashGetAllOperations 测试Hash的获取所有值操作
func TestHashGetAllOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	hash := &RedisCommandHash{}

	t.Run("GetAllHashStringValues", func(t *testing.T) {
		_ = hash.SetHashValueProtoMsg("all_hash", "field1", "value1")
		_ = hash.SetHashValueProtoMsg("all_hash", "field2", "value2")

		result, err := hash.GetAllHashStringValues("all_hash")
		if err != nil {
			t.Errorf("GetAllHashStringValues failed: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 fields, got: %d", len(result))
		}
	})
}

// TestListPopOperations 测试List的Pop操作
func TestListPopOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	list := &RedisCommandList{}

	t.Run("BLPop_EmptyList", func(t *testing.T) {
		result, err := list.BLPop("empty_list")
		if err != nil {
			t.Errorf("BLPop failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty string, got: %s", result)
		}
	})

	t.Run("BLPop_WithData", func(t *testing.T) {
		_ = list.LPush("blpop_list", "value1", "value2")

		result, err := list.BLPop("blpop_list")
		if err != nil {
			t.Errorf("BLPop failed: %v", err)
		}
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("BLPopByte", func(t *testing.T) {
		_ = list.LPush("blpop_byte_list", "value1")

		result, err := list.BLPopByte("blpop_byte_list")
		if err != nil {
			t.Errorf("BLPopByte failed: %v", err)
		}
		if len(result) == 0 {
			t.Error("Expected non-empty byte result")
		}
	})

	t.Run("BRPop_EmptyList", func(t *testing.T) {
		result, err := list.BRPop("empty_list2")
		if err != nil {
			t.Errorf("BRPop failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty string, got: %s", result)
		}
	})

	t.Run("BRPop_WithData", func(t *testing.T) {
		_ = list.RPush("brpop_list", "value1", "value2")

		result, err := list.BRPop("brpop_list")
		if err != nil {
			t.Errorf("BRPop failed: %v", err)
		}
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("BRPopByte", func(t *testing.T) {
		_ = list.RPush("brpop_byte_list", "value1")

		result, err := list.BRPopByte("brpop_byte_list")
		if err != nil {
			t.Errorf("BRPopByte failed: %v", err)
		}
		if len(result) == 0 {
			t.Error("Expected non-empty byte result")
		}
	})

	t.Run("LPop_EmptyList", func(t *testing.T) {
		result, err := list.LPop("empty_list3")
		if err != nil {
			t.Errorf("LPop failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty string, got: %s", result)
		}
	})

	t.Run("RPop_EmptyList", func(t *testing.T) {
		result, err := list.RPop("empty_list4")
		if err != nil {
			t.Errorf("RPop failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty string, got: %s", result)
		}
	})
}

// TestListIndexOperations 测试List的索引操作
func TestListIndexOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	list := &RedisCommandList{}

	t.Run("LIndex_EmptyList", func(t *testing.T) {
		result, err := list.LIndex("empty_list", 0)
		if err != nil {
			t.Errorf("LIndex failed: %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty string, got: %s", result)
		}
	})

	t.Run("LIndex_WithData", func(t *testing.T) {
		_ = list.RPush("index_list", "value1", "value2", "value3")

		result, err := list.LIndex("index_list", 1)
		if err != nil {
			t.Errorf("LIndex failed: %v", err)
		}
		if result != "value2" {
			t.Errorf("Expected 'value2', got: %s", result)
		}
	})
}

// TestListDelKey 测试List的删除键操作
func TestListDelKey(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	list := &RedisCommandList{}

	t.Run("DelKey", func(t *testing.T) {
		_ = list.RPush("del_list", "value1")

		err := list.DelKey("del_list")
		if err != nil {
			t.Errorf("DelKey failed: %v", err)
		}

		// 验证删除
		length, _ := list.LLen("del_list")
		if length != 0 {
			t.Error("List should be deleted")
		}
	})
}

// TestStringGetLock 测试String的分布式锁操作
func TestStringGetLock(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	str := &RedisCommandString{}

	t.Run("GetLock_FirstTime", func(t *testing.T) {
		success, firstTime, err := str.GetLock("lock1", 2, 10, 3)
		if err != nil {
			t.Errorf("GetLock failed: %v", err)
		}
		if !success {
			t.Error("Expected lock to be acquired")
		}
		if !firstTime {
			t.Error("Expected first time lock")
		}
	})

	t.Run("GetLock_AlreadyLocked", func(t *testing.T) {
		// 先获取锁
		_, _, _ = str.GetLock("lock2", 2, 10, 1)

		// 再次尝试获取（应该失败）
		success, _, err := str.GetLock("lock2", 2, 10, 1)
		if err == nil {
			t.Error("Expected error for max retries")
		}
		if success {
			t.Error("Expected lock acquisition to fail")
		}
	})

	t.Run("GetLock_DefaultExpire", func(t *testing.T) {
		success, firstTime, err := str.GetLock("lock3", 0, 10, 3)
		if err != nil {
			t.Errorf("GetLock failed: %v", err)
		}
		if !success {
			t.Error("Expected lock to be acquired")
		}
		if !firstTime {
			t.Error("Expected first time lock")
		}
	})

	t.Run("UnLock", func(t *testing.T) {
		_, _, _ = str.GetLock("lock4", 2, 10, 3)

		err := str.UnLock("lock4")
		if err != nil {
			t.Errorf("UnLock failed: %v", err)
		}

		// 验证锁已释放
		success, _, _ := str.GetLock("lock4", 2, 10, 1)
		if !success {
			t.Error("Expected lock to be available after unlock")
		}
	})
}

// TestZSetRemOperations 测试ZSet的删除操作
func TestZSetRemOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	zset := &RedisCommandZSet{}

	t.Run("ZRem", func(t *testing.T) {
		_ = zset.ZSetAddValue("zset1", "member1", 100)
		_ = zset.ZSetAddValue("zset1", "member2", 200)

		err := zset.ZRem("zset1", "member1")
		if err != nil {
			t.Errorf("ZRem failed: %v", err)
		}

		// 验证删除
		_, err = zset.ZScore("zset1", "member1")
		if err == nil {
			t.Error("member1 should be removed")
		}
	})

	t.Run("ZRemRangeByScore", func(t *testing.T) {
		_ = zset.ZSetAddValue("zset2", "member1", 100)
		_ = zset.ZSetAddValue("zset2", "member2", 200)
		_ = zset.ZSetAddValue("zset2", "member3", 300)

		err := zset.ZRemRangeByScore("zset2", "100", "200")
		if err != nil {
			t.Errorf("ZRemRangeByScore failed: %v", err)
		}

		// 验证删除
		count, _ := zset.GetZSetCount("zset2")
		if count != 1 {
			t.Errorf("Expected 1 member remaining, got: %d", count)
		}
	})
}

// TestZSetRangeOperations 测试ZSet的范围操作
func TestZSetRangeOperations(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	zset := &RedisCommandZSet{}

	t.Run("ZRange", func(t *testing.T) {
		_ = zset.ZSetAddValue("zset_range", "member1", 100)
		_ = zset.ZSetAddValue("zset_range", "member2", 200)
		_ = zset.ZSetAddValue("zset_range", "member3", 300)

		result, err := zset.ZRange("zset_range", 0, 1)
		if err != nil {
			t.Errorf("ZRange failed: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 members, got: %d", len(result))
		}
	})

	t.Run("ZRevRange", func(t *testing.T) {
		_ = zset.ZSetAddValue("zset_revrange", "member1", 100)
		_ = zset.ZSetAddValue("zset_revrange", "member2", 200)
		_ = zset.ZSetAddValue("zset_revrange", "member3", 300)

		result, err := zset.ZRevRange("zset_revrange", 0, 1)
		if err != nil {
			t.Errorf("ZRevRange failed: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 members, got: %d", len(result))
		}
		if result[0] != "member3" {
			t.Errorf("Expected 'member3' first, got: %s", result[0])
		}
	})
}

// TestPoolOperations 测试连接池操作
func TestPoolOperations(t *testing.T) {
	mr := miniredis.RunT(t)
	defer mr.Close()

	pool := NewPool()

	t.Run("IsHealthy_BeforeInit", func(t *testing.T) {
		healthy := pool.IsHealthy()
		if healthy {
			t.Error("Pool should not be healthy before initialization")
		}
	})

	t.Run("IsHealthy_AfterInit", func(t *testing.T) {
		config := &RedisConfig{
			Addr:  mr.Addr(),
			Index: 0,
		}

		err := pool.LoadConfig(config)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		healthy := pool.IsHealthy()
		if !healthy {
			t.Error("Pool should be healthy after initialization")
		}
	})

	t.Run("IsHealthy_AfterClose", func(t *testing.T) {
		pool.Close()

		healthy := pool.IsHealthy()
		if healthy {
			t.Error("Pool should not be healthy after close")
		}
	})
}

// TestContextOperations 测试上下文操作
func TestContextOperations(t *testing.T) {
	t.Run("ContextWithTimeout_Default", func(t *testing.T) {
		ctx, cancel := ContextWithTimeout(0)
		defer cancel()

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("Expected deadline to be set")
		}

		expectedDeadline := time.Now().Add(DefaultTimeout)
		if deadline.Before(expectedDeadline.Add(-time.Second)) || deadline.After(expectedDeadline.Add(time.Second)) {
			t.Error("Deadline not within expected range")
		}
	})

	t.Run("ContextWithTimeout_Custom", func(t *testing.T) {
		customTimeout := 10 * time.Second
		ctx, cancel := ContextWithTimeout(customTimeout)
		defer cancel()

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("Expected deadline to be set")
		}

		expectedDeadline := time.Now().Add(customTimeout)
		if deadline.Before(expectedDeadline.Add(-time.Second)) || deadline.After(expectedDeadline.Add(time.Second)) {
			t.Error("Deadline not within expected range")
		}
	})

	t.Run("ContextWithDefaultTimeout", func(t *testing.T) {
		ctx, cancel := ContextWithDefaultTimeout()
		defer cancel()

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("Expected deadline to be set")
		}

		expectedDeadline := time.Now().Add(DefaultTimeout)
		if deadline.Before(expectedDeadline.Add(-time.Second)) || deadline.After(expectedDeadline.Add(time.Second)) {
			t.Error("Deadline not within expected range")
		}
	})
}

// TestHashGetValueIntExceptNil 测试GetHashValueIntExceptNil
func TestHashGetValueIntExceptNil(t *testing.T) {
	_, cleanup := setupTestRedis3(t)
	defer cleanup()

	hash := &RedisCommandHash{}

	t.Run("GetHashValueIntExceptNil_Exists", func(t *testing.T) {
		_ = hash.SetHashValueProtoMsg("int_hash", "field1", "123")

		result, err := hash.GetHashValueIntExceptNil("int_hash", "field1")
		if err != nil {
			t.Errorf("GetHashValueIntExceptNil failed: %v", err)
		}
		if result != 123 {
			t.Errorf("Expected 123, got: %d", result)
		}
	})

	t.Run("GetHashValueIntExceptNil_NotExists", func(t *testing.T) {
		_, err := hash.GetHashValueIntExceptNil("int_hash", "nonexistent")
		if err == nil {
			t.Error("Expected error for nonexistent field")
		}
	})
}

// TestNilClientScenarios 测试nil客户端场景
func TestNilClientScenarios(t *testing.T) {
	// 保存原始pool
	originalPool := gPoolM
	defer func() {
		gPoolM = originalPool
	}()

	// 创建一个关闭的pool
	gPoolM = NewPool()
	gPoolM.Close()

	t.Run("Hash_SetValueKeyLock_InvalidType", func(t *testing.T) {
		hash := &RedisCommandHash{}
		_, err := hash.SetValueKeyLock("key", "invalid_type")
		if err == nil {
			t.Error("Expected error for invalid expiration type")
		}
	})

	t.Run("Hash_IncHashValue_InvalidType", func(t *testing.T) {
		hash := &RedisCommandHash{}
		_, err := hash.IncHashValue("key", "field", "invalid")
		if err == nil {
			t.Error("Expected error for invalid increment type")
		}
	})
}

// TestRealClientOperations 使用真实客户端测试
func TestRealClientOperations(t *testing.T) {
	mr := miniredis.RunT(t)
	defer mr.Close()

	// 创建真实的Redis客户端
	client := goredis.NewClient(&goredis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	t.Run("DirectClientOperations", func(t *testing.T) {
		// 测试基本操作
		err := client.Set(ctx, "test_key", "test_value", 0).Err()
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		val, err := client.Get(ctx, "test_key").Result()
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if val != "test_value" {
			t.Errorf("Expected 'test_value', got: %s", val)
		}
	})
}
