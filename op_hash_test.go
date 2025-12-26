package mredis

import (
	"testing"
)

func TestHashCommandWithoutClient(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandHash{}

	t.Run("Exists", func(t *testing.T) {
		exists := cmd.Exists("hash:test", "field1")
		if exists {
			t.Error("Exists should return false when client is nil")
		}
	})

	t.Run("DeleteHashKeyList", func(t *testing.T) {
		err := cmd.DeleteHashKeyList("hash:test", []uint64{1, 2, 3})
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("DeleteHashStringKeyList", func(t *testing.T) {
		err := cmd.DeleteHashStringKeyList("hash:test", []string{"field1", "field2"})
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashLen", func(t *testing.T) {
		_, err := cmd.GetHashLen("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("HLen", func(t *testing.T) {
		_, err := cmd.HLen("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashKeys", func(t *testing.T) {
		_, err := cmd.GetHashKeys("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("HKEYS", func(t *testing.T) {
		_, err := cmd.HKEYS("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("DelKey", func(t *testing.T) {
		err := cmd.DelKey("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("SetHashValueProtoMsg", func(t *testing.T) {
		err := cmd.SetHashValueProtoMsg("hash:test", "field1", "value1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("SetHashValue", func(t *testing.T) {
		err := cmd.SetHashValue("hash:test", "field1", "value1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("SetHashExpire", func(t *testing.T) {
		err := cmd.SetHashExpire("hash:test", 60)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashValueString", func(t *testing.T) {
		_, err := cmd.GetHashValueString("hash:test", "field1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashValueInt", func(t *testing.T) {
		_, err := cmd.GetHashValueInt("hash:test", "field1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashValueByte", func(t *testing.T) {
		_, err := cmd.GetHashValueByte("hash:test", "field1")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("IncHashValue", func(t *testing.T) {
		_, err := cmd.IncHashValue("hash:test", "field1", 1)
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetAllHashValues", func(t *testing.T) {
		_, err := cmd.GetAllHashValues("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashValues", func(t *testing.T) {
		_, err := cmd.GetHashValues("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashSize", func(t *testing.T) {
		_, err := cmd.GetHashSize("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("HashExists", func(t *testing.T) {
		_, err := cmd.HashExists("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})
}

func TestHashCommandHelperMethods(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandHash{}

	t.Run("GetHashValuesAsJSON", func(t *testing.T) {
		_, err := cmd.GetHashValuesAsJSON("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashValuesForAPI", func(t *testing.T) {
		_, err := cmd.GetHashValuesForAPI("hash:test")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashValuesWithPrefix", func(t *testing.T) {
		_, err := cmd.GetHashValuesWithPrefix("hash:test", "prefix_")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})

	t.Run("GetHashValuesExcludePrefix", func(t *testing.T) {
		_, err := cmd.GetHashValuesExcludePrefix("hash:test", "prefix_")
		if err != ErrClientNil {
			t.Errorf("Expected ErrClientNil, got %v", err)
		}
	})
}

func TestIncHashValueTypes(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandHash{}

	tests := []struct {
		name    string
		incNum  interface{}
		wantErr bool
	}{
		{"int", int(1), true},
		{"int32", int32(1), true},
		{"int64", int64(1), true},
		{"uint32", uint32(1), true},
		{"invalid", "string", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmd.IncHashValue("hash:test", "field", tt.incNum)
			if tt.wantErr && err == nil {
				t.Error("Expected error")
			}
		})
	}
}

func TestSetValueKeyLock(t *testing.T) {
	pool := NewPool()
	oldPool := gPoolM
	gPoolM = pool
	defer func() { gPoolM = oldPool }()

	cmd := &RedisCommandHash{}

	tests := []struct {
		name       string
		expiredSec interface{}
		wantErr    bool
	}{
		{"uint32", uint32(10), true},
		{"int", int(10), true},
		{"int64", int64(10), true},
		{"invalid", "string", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmd.SetValueKeyLock("lock:test", tt.expiredSec)
			if tt.wantErr && err == nil {
				t.Error("Expected error")
			}
		})
	}
}
