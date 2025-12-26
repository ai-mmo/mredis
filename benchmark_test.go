package mredis

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkPoolCreation 测试连接池创建性能
func BenchmarkPoolCreation(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pool := NewPool()
		_ = pool
	}
}

// BenchmarkContextCreation 测试上下文创建性能
func BenchmarkContextCreation(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx, cancel := ContextWithDefaultTimeout()
		cancel()
		_ = ctx
	}
}

// BenchmarkContextWithTimeout 测试带超时的上下文创建性能
func BenchmarkContextWithTimeout(b *testing.B) {
	b.ReportAllocs()
	timeout := 5 * time.Second
	for i := 0; i < b.N; i++ {
		ctx, cancel := ContextWithTimeout(timeout)
		cancel()
		_ = ctx
	}
}

// BenchmarkErrorComparison 测试错误比较性能
func BenchmarkErrorComparison(b *testing.B) {
	b.ReportAllocs()
	err := ErrClientNil
	for i := 0; i < b.N; i++ {
		_ = err == ErrClientNil
	}
}

// BenchmarkPoolIsClosed 测试连接池状态检查性能
func BenchmarkPoolIsClosed(b *testing.B) {
	pool := NewPool()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pool.IsClosed()
	}
}

// BenchmarkGetVersion 测试获取版本号性能
func BenchmarkGetVersion(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = GetVersion()
	}
}

// BenchmarkFormatKey 测试键格式化性能
func BenchmarkFormatKey(b *testing.B) {
	b.ReportAllocs()
	key := "test:key"
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", key)
	}
}

// BenchmarkFormatKeyWithInt 测试整数键格式化性能
func BenchmarkFormatKeyWithInt(b *testing.B) {
	b.ReportAllocs()
	key := 12345
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", key)
	}
}

// BenchmarkSplitFromScore 测试分数拆分性能
func BenchmarkSplitFromScore(b *testing.B) {
	b.ReportAllocs()
	score := uint64(0x123456789ABCDEF0)
	for i := 0; i < b.N; i++ {
		_, _ = SplitFromScore(score)
	}
}

// BenchmarkPoolGetClientNil 测试获取nil客户端性能
func BenchmarkPoolGetClientNil(b *testing.B) {
	pool := NewPool()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pool.GetClient()
	}
}

// BenchmarkLoadConfigValidation 测试配置验证性能
func BenchmarkLoadConfigValidation(b *testing.B) {
	pool := NewPool()
	config := &RedisConfig{
		Addr:        "localhost:6379",
		PoolSize:    10,
		DialTimeout: 5 * time.Second,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// 注意：这会尝试连接Redis，在benchmark中可能失败
		_ = pool.LoadConfig(config)
	}
}

// BenchmarkParallelPoolIsClosed 并发测试连接池状态检查
func BenchmarkParallelPoolIsClosed(b *testing.B) {
	pool := NewPool()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.IsClosed()
		}
	})
}

// BenchmarkParallelContextCreation 并发测试上下文创建
func BenchmarkParallelContextCreation(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := ContextWithDefaultTimeout()
			cancel()
			_ = ctx
		}
	})
}

// BenchmarkParallelGetVersion 并发测试获取版本号
func BenchmarkParallelGetVersion(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetVersion()
		}
	})
}

// BenchmarkMemoryAllocation 测试内存分配
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// 模拟常见操作的内存分配
		_ = make([]string, 0, 10)
		_ = make(map[string]string, 10)
		_ = fmt.Sprintf("key:%d", i)
	}
}

// BenchmarkStringConcatenation 测试字符串拼接性能
func BenchmarkStringConcatenation(b *testing.B) {
	b.Run("Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("prefix:%d:suffix", i)
		}
	})

	b.Run("Plus", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = "prefix:" + fmt.Sprint(i) + ":suffix"
		}
	})
}

// BenchmarkSliceAppend 测试切片追加性能
func BenchmarkSliceAppend(b *testing.B) {
	b.Run("WithCapacity", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			slice := make([]string, 0, 100)
			for j := 0; j < 100; j++ {
				slice = append(slice, "value")
			}
		}
	})

	b.Run("WithoutCapacity", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var slice []string
			for j := 0; j < 100; j++ {
				slice = append(slice, "value")
			}
		}
	})
}

// BenchmarkMapAccess 测试Map访问性能
func BenchmarkMapAccess(b *testing.B) {
	m := make(map[string]string, 1000)
	for i := 0; i < 1000; i++ {
		m[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = m["key500"]
	}
}

// BenchmarkMapIteration 测试Map遍历性能
func BenchmarkMapIteration(b *testing.B) {
	m := make(map[string]string, 1000)
	for i := 0; i < 1000; i++ {
		m[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for k, v := range m {
			_ = k
			_ = v
		}
	}
}
