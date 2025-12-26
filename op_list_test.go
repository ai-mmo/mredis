package mredis

import (
	"testing"
)

func TestListCommandWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandList{}

	t.Run("BLPop", func(t *testing.T) {
		_, err := cmd.BLPop("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("BLPopByte", func(t *testing.T) {
		_, err := cmd.BLPopByte("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("BRPop", func(t *testing.T) {
		_, err := cmd.BRPop("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("BRPopByte", func(t *testing.T) {
		_, err := cmd.BRPopByte("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LPush", func(t *testing.T) {
		err := cmd.LPush("list:test", "value1", "value2")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("RPush", func(t *testing.T) {
		err := cmd.RPush("list:test", "value1", "value2")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LPop", func(t *testing.T) {
		_, err := cmd.LPop("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("RPop", func(t *testing.T) {
		_, err := cmd.RPop("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LLen", func(t *testing.T) {
		_, err := cmd.LLen("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LRange", func(t *testing.T) {
		_, err := cmd.LRange("list:test", 0, -1)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LTrim", func(t *testing.T) {
		err := cmd.LTrim("list:test", 0, 10)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LIndex", func(t *testing.T) {
		_, err := cmd.LIndex("list:test", 0)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LSet", func(t *testing.T) {
		err := cmd.LSet("list:test", 0, "value")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LRem", func(t *testing.T) {
		_, err := cmd.LRem("list:test", 1, "value")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("DelKey", func(t *testing.T) {
		err := cmd.DelKey("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetListLen", func(t *testing.T) {
		_, err := cmd.GetListLen("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetListRange", func(t *testing.T) {
		_, err := cmd.GetListRange("list:test", 0, -1)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("PushListValue", func(t *testing.T) {
		err := cmd.PushListValue("list:test", "value1", "value2")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("PopListValue", func(t *testing.T) {
		_, err := cmd.PopListValue("list:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("LRangeBytes", func(t *testing.T) {
		_, err := cmd.LRangeBytes("list:test", 0, -1)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})
}
