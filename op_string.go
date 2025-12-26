package mredis

import (
	"context"
	"fmt"
	"mlog"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisCommandString struct{}

// DelStringKey 删除字符串类型的键
// 参数:
//   - key: 要删除的键名
//
// 返回:
//   - error: 删除失败时返回错误
func (cm *RedisCommandString) DelStringKey(key any) error {
	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	return client.Del(ctx, keyStr).Err()
}

// SetStringValue 设置字符串类型的键值
// 参数:
//   - key: 键名
//   - value: 要设置的值
//
// 返回:
//   - error: 设置失败时返回错误
//
// 注意:
//   - 不设置过期时间，键将永久存在
func (cm *RedisCommandString) SetStringValue(key, value any) error {
	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	valueStr := fmt.Sprintf("%v", value)
	return client.Set(ctx, keyStr, valueStr, 0).Err()
}

// GetStringValue 获取字符串类型的键值
// 参数:
//   - key: 键名
//
// 返回:
//   - string: 键对应的值，键不存在时返回空字符串
//   - error: 获取失败时返回错误
func (cm *RedisCommandString) GetStringValue(key any) (string, error) {
	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return "", ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	result, err := client.Get(ctx, keyStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// GetIntValue 获取字符串类型键的整数值
// 参数:
//   - key: 键名
//
// 返回:
//   - int: 键对应的整数值
//   - error: 获取失败或值无法转换为整数时返回错误
func (cm *RedisCommandString) GetIntValue(key any) (int, error) {
	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	return client.Get(ctx, keyStr).Int()
}

// Exists 检查键是否存在
// 参数:
//   - key: 键名
//
// 返回:
//   - int: 键存在返回1，不存在返回0
//   - error: 检查失败时返回错误
func (cm *RedisCommandString) Exists(key any) (int, error) {
	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	result, err := client.Exists(ctx, keyStr).Result()
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

// GetInt64Value 获取字符串类型键的64位整数值
// 参数:
//   - key: 键名
//
// 返回:
//   - int64: 键对应的64位整数值
//   - error: 获取失败或值无法转换为int64时返回错误
func (cm *RedisCommandString) GetInt64Value(key any) (int64, error) {
	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	return client.Get(ctx, keyStr).Int64()
}

// GetByteValue 获取字符串类型键的字节数组值
// 参数:
//   - key: 键名
//
// 返回:
//   - []byte: 键对应的字节数组，键不存在时返回nil
//   - error: 获取失败时返回错误
func (cm *RedisCommandString) GetByteValue(key any) ([]byte, error) {
	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	result, err := client.Get(ctx, keyStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return []byte(result), nil
}

// Incr 将键的整数值增加1
// 参数:
//   - key: 键名
//
// 返回:
//   - int64: 增加后的值
//   - error: 操作失败时返回错误
//
// 注意:
//   - 如果键不存在，会先初始化为0再执行增加操作
//   - 如果键的值不是整数，会返回错误
func (cm *RedisCommandString) Incr(key any) (int64, error) {
	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	// INCR 命令不会返回 NilReply，但如果 key 原本不存在，Redis 会初始化为 0 并自增为 1
	return client.Incr(ctx, keyStr).Result()
}

// GetLock 获取分布式锁
// 参数:
//   - key: 锁的键名
//   - expire: 锁的过期时间（秒），0表示使用默认值2秒
//   - skipTime: 获取锁失败后的重试间隔时间（毫秒）
//   - maxLimit: 最大重试次数，0表示无限重试
//
// 返回:
//   - bool: 是否成功获取锁
//   - bool: 是否首次尝试就获取到锁
//   - error: 获取锁失败时返回错误
//
// 注意:
//   - 使用SetNX实现，保证原子性
//   - 锁会自动过期，避免死锁
//   - 重试次数超过maxLimit时返回错误
func (cm *RedisCommandString) GetLock(key string, expire int, skipTime, maxLimit int) (bool, bool, error) {
	if expire == 0 {
		expire = 2 // 锁默认过期时间2s
	}

	// 获取Redis客户端
	client := GetClient()
	if client == nil {
		return false, false, ErrClientNil
	}

	ctx := context.Background()
	expireDuration := time.Duration(expire) * time.Second

	limit := 1 // 第limit次
	for {
		// 使用SetNX实现分布式锁
		result, err := client.SetNX(ctx, key, "locked", expireDuration).Result()
		if err != nil {
			// 拿锁失败
			mlog.Warn("GetLock err,err:%v", err)
			limit++
			if maxLimit != 0 && limit > maxLimit {
				return false, false, fmt.Errorf("%v lock too many retries", key)
			} else {
				// 需要继续拿锁时才延迟
				time.Sleep(time.Duration(skipTime) * time.Millisecond)
			}
		} else if result {
			// 拿锁成功
			if limit == 1 {
				return true, true, nil
			} else {
				return true, false, nil
			}
		} else {
			// 锁已存在，继续重试
			limit++
			if maxLimit != 0 && limit > maxLimit {
				return false, false, fmt.Errorf("%v lock too many retries", key)
			} else {
				time.Sleep(time.Duration(skipTime) * time.Millisecond)
			}
		}
	}
}

// UnLock 释放分布式锁
// 参数:
//   - key: 锁的键名
//
// 返回:
//   - error: 释放锁失败时返回错误
//
// 注意:
//   - 通过删除键来释放锁
func (cm *RedisCommandString) UnLock(key string) error {
	err := cm.DelStringKey(key)
	if err != nil {
		mlog.Warn("UnLock err,err:%v", err)
		return err
	}
	return nil
}
