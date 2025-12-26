package mredis

import (
	"fmt"
	"mlog"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func SetRedisKey(key, value any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx, cancel := ContextWithDefaultTimeout()
	defer cancel()

	keyStr := fmt.Sprintf("%v", key)
	valueStr := fmt.Sprintf("%v", value)
	return client.Set(ctx, keyStr, valueStr, 0).Err()
}

func GetRedisKeyType(key string) (string, error) {
	client := GetClient()
	if client == nil {
		return "", ErrClientNil
	}

	ctx, cancel := ContextWithDefaultTimeout()
	defer cancel()

	result, err := client.Type(ctx, key).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

func DeleteRedisKey(key any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx, cancel := ContextWithDefaultTimeout()
	defer cancel()

	keyStr := fmt.Sprintf("%v", key)
	return client.Del(ctx, keyStr).Err()
}

// GetRedisKeysByPattern 获取匹配模式的所有Redis键
func GetRedisKeysByPattern(pattern string) ([]string, error) {
	var allKeys []string

	err := scanRedisKeys(pattern, 100, func(keys []string) bool {
		allKeys = append(allKeys, keys...)
		return true // 继续扫描
	})

	return allKeys, err
}

// DeleteRedisKeysByPattern 删除匹配模式的所有Redis键
func DeleteRedisKeysByPattern(pattern string) ([]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx, cancel := ContextWithTimeout(30 * time.Second) // 删除操作可能需要更长时间
	defer cancel()

	var deletedKeys []string

	err := scanRedisKeys(pattern, 100, func(keys []string) bool {
		// 批量删除键
		if len(keys) > 0 {
			deleted, err := client.Del(ctx, keys...).Result()
			if err != nil {
				mlog.Warn("[DeleteRedisKeysByPattern] batch delete failed: %v", err)
				// 尝试逐个删除
				for _, key := range keys {
					if err = client.Del(ctx, key).Err(); err != nil {
						mlog.Warn("[DeleteRedisKeysByPattern] delete key [%s] failed: %v", key, err)
					} else {
						deletedKeys = append(deletedKeys, key)
					}
				}
			} else {
				// 批量删除成功
				if deleted > 0 {
					deletedKeys = append(deletedKeys, keys...)
				}
			}
		}
		return true // 继续扫描
	})

	return deletedKeys, err
}

// scanRedisKeys 通用的Redis键扫描方法
// pattern: 匹配模式，如 "*game_impl*:123*"
// batchSize: 每次扫描的批次大小，建议100-1000
// callback: 对每批键的处理回调函数，返回true继续扫描，false停止
func scanRedisKeys(pattern string, batchSize int64, callback func(keys []string) bool) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx, cancel := ContextWithTimeout(30 * time.Second) // SCAN操作可能需要较长时间
	defer cancel()

	var cursor uint64

	if batchSize <= 0 {
		batchSize = 100 // 默认批次大小
	}

	for {
		// 使用SCAN命令分批获取键，避免阻塞
		keys, nextCursor, err := client.Scan(ctx, cursor, pattern, batchSize).Result()
		if err != nil {
			return fmt.Errorf("scan keys failed: %w", err)
		}

		// 如果有键，调用回调函数处理
		if len(keys) > 0 {
			if !callback(keys) {
				break // 回调返回false，停止扫描
			}
		}

		// 检查是否扫描完成
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
