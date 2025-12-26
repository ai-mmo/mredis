package mredis

import (
	"testing"
)

func TestStringCommandWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandString{}

	t.Run("DelStringKey", func(t *testing.T) {
		err := cmd.DelStringKey("test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("SetStringValue", func(t *testing.T) {
		err := cmd.SetStringValue("test", "value")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetStringValue", func(t *testing.T) {
		_, err := cmd.GetStringValue("test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetIntValue", func(t *testing.T) {
		_, err := cmd.GetIntValue("test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		_, err := cmd.Exists("test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetInt64Value", func(t *testing.T) {
		_, err := cmd.GetInt64Value("test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetByteValue", func(t *testing.T) {
		_, err := cmd.GetByteValue("test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("Incr", func(t *testing.T) {
		_, err := cmd.Incr("test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetLock", func(t *testing.T) {
		_, _, err := cmd.GetLock("lock:test", 10, 100, 1)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("UnLock", func(t *testing.T) {
		err := cmd.UnLock("lock:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})
}

func TestGetLockDefaultExpire(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandString{}

	// 测试默认过期时间
	_, _, err := cmd.GetLock("lock:test", 0, 100, 1)
	if err != ErrClientNil {
		t.Errorf("Expected ErrClientNil, got %v", err)
	}
}
