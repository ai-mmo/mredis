package mredis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

type RedisCommandList struct{}

// BLPop 从列表左侧弹出元素
// 参数:
//   - key: 列表键名
//
// 返回:
//   - string: 弹出的元素值，列表为空时返回空字符串
//   - error: 操作失败时返回错误
//
// 注意:
//   - 此方法实际使用LPop实现，非阻塞操作
func (mmo *RedisCommandList) BLPop(key any) (string, error) {
	client := GetClient()
	if client == nil {
		return "", ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	result, err := client.LPop(ctx, keyStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// BLPopByte 从列表左侧弹出元素（字节格式）
// 参数:
//   - key: 列表键名
//
// 返回:
//   - []byte: 弹出的元素字节数组，列表为空时返回nil
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) BLPopByte(key any) ([]byte, error) {
	result, err := mmo.BLPop(key)
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

// BRPop 从列表右侧弹出元素
// 参数:
//   - key: 列表键名
//
// 返回:
//   - string: 弹出的元素值，列表为空时返回空字符串
//   - error: 操作失败时返回错误
//
// 注意:
//   - 此方法实际使用RPop实现，非阻塞操作
func (mmo *RedisCommandList) BRPop(key any) (string, error) {
	client := GetClient()
	if client == nil {
		return "", ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	result, err := client.RPop(ctx, keyStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// BRPopByte 从列表右侧弹出元素（字节格式）
// 参数:
//   - key: 列表键名
//
// 返回:
//   - []byte: 弹出的元素字节数组，列表为空时返回nil
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) BRPopByte(key any) ([]byte, error) {
	result, err := mmo.BRPop(key)
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

// LPush 将一个或多个值插入到列表头部（左侧）
// 参数:
//   - key: 列表键名
//   - values: 要插入的值，可变参数
//
// 返回:
//   - error: 操作失败时返回错误
//
// 注意:
//   - 如果列表不存在，会创建一个新列表
//   - 多个值按顺序插入，最后一个值会成为列表头部
func (mmo *RedisCommandList) LPush(key any, values ...any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.LPush(ctx, keyStr, values...).Err()
}

// RPush 将一个或多个值插入到列表尾部（右侧）
// 参数:
//   - key: 列表键名
//   - values: 要插入的值，可变参数
//
// 返回:
//   - error: 操作失败时返回错误
//
// 注意:
//   - 如果列表不存在，会创建一个新列表
//   - 多个值按顺序插入到列表尾部
func (mmo *RedisCommandList) RPush(key any, values ...any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.RPush(ctx, keyStr, values...).Err()
}

// LPop 从列表左侧弹出一个元素
// 参数:
//   - key: 列表键名
//
// 返回:
//   - string: 弹出的元素值，列表为空时返回空字符串
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) LPop(key any) (string, error) {
	client := GetClient()
	if client == nil {
		return "", ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	result, err := client.LPop(ctx, keyStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// RPop 从列表右侧弹出一个元素
// 参数:
//   - key: 列表键名
//
// 返回:
//   - string: 弹出的元素值，列表为空时返回空字符串
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) RPop(key any) (string, error) {
	client := GetClient()
	if client == nil {
		return "", ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	result, err := client.RPop(ctx, keyStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// LLen 获取列表的长度
// 参数:
//   - key: 列表键名
//
// 返回:
//   - int64: 列表长度，列表不存在时返回0
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) LLen(key any) (int64, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.LLen(ctx, keyStr).Result()
}

// LRange 获取列表指定范围内的元素
// 参数:
//   - key: 列表键名
//   - start: 起始索引（0表示第一个元素）
//   - stop: 结束索引（-1表示最后一个元素）
//
// 返回:
//   - []string: 指定范围内的元素列表
//   - error: 操作失败时返回错误
//
// 注意:
//   - 索引可以是负数，-1表示最后一个元素，-2表示倒数第二个元素
func (mmo *RedisCommandList) LRange(key any, start, stop int64) ([]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.LRange(ctx, keyStr, start, stop).Result()
}

// LTrim 修剪列表，只保留指定区间内的元素
// 参数:
//   - key: 列表键名
//   - start: 起始索引
//   - stop: 结束索引
//
// 返回:
//   - error: 操作失败时返回错误
//
// 注意:
//   - 区间外的元素会被删除
//   - 索引可以是负数
func (mmo *RedisCommandList) LTrim(key any, start, stop int64) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.LTrim(ctx, keyStr, start, stop).Err()
}

// LIndex 获取列表指定位置的元素
// 参数:
//   - key: 列表键名
//   - index: 元素索引（0表示第一个元素）
//
// 返回:
//   - string: 指定位置的元素值，索引超出范围时返回空字符串
//   - error: 操作失败时返回错误
//
// 注意:
//   - 索引可以是负数，-1表示最后一个元素
func (mmo *RedisCommandList) LIndex(key any, index int64) (string, error) {
	client := GetClient()
	if client == nil {
		return "", ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	result, err := client.LIndex(ctx, keyStr, index).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// LSet 设置列表指定位置的元素值
// 参数:
//   - key: 列表键名
//   - index: 元素索引
//   - value: 新的元素值
//
// 返回:
//   - error: 操作失败时返回错误
//
// 注意:
//   - 索引必须在有效范围内，否则返回错误
//   - 索引可以是负数
func (mmo *RedisCommandList) LSet(key any, index int64, value any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	valueStr := fmt.Sprintf("%v", value)

	return client.LSet(ctx, keyStr, index, valueStr).Err()
}

// LRem 移除列表中与给定值相等的元素
// 参数:
//   - key: 列表键名
//   - count: 移除数量，>0从头到尾移除，<0从尾到头移除，=0移除所有
//   - value: 要移除的元素值
//
// 返回:
//   - int64: 实际移除的元素数量
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) LRem(key any, count int64, value any) (int64, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	valueStr := fmt.Sprintf("%v", value)

	return client.LRem(ctx, keyStr, count, valueStr).Result()
}

// DelKey 删除整个列表键
// 参数:
//   - key: 列表键名
//
// 返回:
//   - error: 删除失败时返回错误
func (mmo *RedisCommandList) DelKey(key any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.Del(ctx, keyStr).Err()
}

// GetListLen 获取列表长度（兼容旧API）
// 参数:
//   - key: 列表键名
//
// 返回:
//   - int: 列表长度
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) GetListLen(key any) (int, error) {
	result, err := mmo.LLen(key)
	return int(result), err
}

// GetListRange 获取列表指定范围内的元素（兼容旧API）
// 参数:
//   - key: 列表键名
//   - start: 起始索引
//   - stop: 结束索引
//
// 返回:
//   - []string: 指定范围内的元素列表
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) GetListRange(key any, start, stop int64) ([]string, error) {
	return mmo.LRange(key, start, stop)
}

// PushListValue 向列表尾部推入值（兼容旧API）
// 参数:
//   - key: 列表键名
//   - values: 要推入的值，可变参数
//
// 返回:
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) PushListValue(key any, values ...any) error {
	return mmo.RPush(key, values...)
}

// PopListValue 从列表头部弹出值（兼容旧API）
// 参数:
//   - key: 列表键名
//
// 返回:
//   - string: 弹出的元素值
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) PopListValue(key any) (string, error) {
	return mmo.LPop(key)
}

// LRangeBytes 获取列表指定范围内的元素（字节格式）
// 参数:
//   - key: 列表键名
//   - start: 起始索引
//   - stop: 结束索引
//
// 返回:
//   - [][]byte: 指定范围内的元素字节数组列表
//   - error: 操作失败时返回错误
func (mmo *RedisCommandList) LRangeBytes(key any, start, stop int64) ([][]byte, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	// 使用LRANGE获取列表范围内的元素
	result, err := client.LRange(ctx, keyStr, start, stop).Result()
	if err != nil {
		return nil, err
	}

	// 转换字符串数组为字节数组
	byteResults := make([][]byte, len(result))
	for i, str := range result {
		byteResults[i] = []byte(str)
	}

	return byteResults, nil
}
