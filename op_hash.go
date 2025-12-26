package mredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mlog"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

type RedisCommandHash struct{}

// Exists 检查Hash中指定字段是否存在
// 参数:
//   - key: Hash键名
//   - field: 字段名
//
// 返回:
//   - bool: 字段存在返回true，否则返回false
func (cm *RedisCommandHash) Exists(key, field any) bool {
	client := GetClient()
	if client == nil {
		return false
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	fieldStr := fmt.Sprintf("%v", field)

	result, err := client.HExists(ctx, keyStr, fieldStr).Result()
	if err != nil {
		return false
	}
	return result
}

// DeleteHashKeyList 删除Hash中的多个uint64类型字段
// 参数:
//   - hashKey: Hash键名
//   - deleteList: 要删除的字段列表（uint64类型）
//
// 返回:
//   - error: 删除失败时返回错误
func (cm *RedisCommandHash) DeleteHashKeyList(hashKey string, deleteList []uint64) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()

	// 转换uint64列表为字符串列表
	fields := make([]string, len(deleteList))
	for i, key := range deleteList {
		fields[i] = strconv.FormatUint(key, 10)
	}

	return client.HDel(ctx, hashKey, fields...).Err()
}

// DeleteHashStringKeyList 删除Hash中的多个字符串类型字段
// 参数:
//   - hashKey: Hash键名
//   - deleteStrList: 要删除的字段列表（字符串类型）
//
// 返回:
//   - error: 删除失败时返回错误
func (cm *RedisCommandHash) DeleteHashStringKeyList(hashKey any, deleteStrList []string) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	return client.HDel(ctx, hashKeyStr, deleteStrList...).Err()
}

// HLen 获取Hash的字段数量（兼容旧API）
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - int: Hash中的字段数量
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) HLen(hashKey string) (int, error) {
	return cm.GetHashLen(hashKey)
}

// GetHashLen 获取Hash的字段数量
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - int: Hash中的字段数量
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashLen(hashKey any) (int, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	result, err := client.HLen(ctx, hashKeyStr).Result()
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

// HKEYS 获取Hash的所有字段名（兼容旧API）
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - []string: 所有字段名列表
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) HKEYS(hashKey string) ([]string, error) {
	return cm.GetHashKeys(hashKey)
}

// GetHashKeys 获取Hash的所有字段名
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - []string: 所有字段名列表
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashKeys(hashKey any) ([]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	return client.HKeys(ctx, hashKeyStr).Result()
}

// DelKey 删除整个Hash键
// 参数:
//   - key: Hash键名
//
// 返回:
//   - error: 删除失败时返回错误
func (cm *RedisCommandHash) DelKey(key any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.Del(ctx, keyStr).Err()
}

// GetHashValueProtoMsg 获取Hash字段的Protobuf消息并反序列化
// 参数:
//   - hashKey: Hash键名
//   - hashField: 字段名
//   - protoMsg: Protobuf消息对象指针，用于接收反序列化结果
//
// 返回:
//   - error: 获取或反序列化失败时返回错误
//
// 注意:
//   - 字段不存在时返回nil
func (cm *RedisCommandHash) GetHashValueProtoMsg(hashKey, hashField any, protoMsg proto.Message) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)
	hashFieldStr := fmt.Sprintf("%v", hashField)

	result, err := client.HGet(ctx, hashKeyStr, hashFieldStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil
		}
		return err
	}

	return proto.Unmarshal([]byte(result), protoMsg)
}

// SetHashValue 设置Hash字段值（兼容旧API）
// 参数:
//   - hashKey: Hash键名
//   - hashField: 字段名
//   - value: 字段值
//
// 返回:
//   - error: 设置失败时返回错误
func (cm *RedisCommandHash) SetHashValue(hashKey, hashField, value any) error {
	return cm.SetHashValueProtoMsg(hashKey, hashField, value)
}

// SetHashValueProtoMsg 设置Hash字段值（支持Protobuf消息、字节数组和其他类型）
// 参数:
//   - hashKey: Hash键名
//   - hashField: 字段名
//   - value: 字段值，支持proto.Message、[]byte或其他类型
//
// 返回:
//   - error: 设置失败时返回错误
//
// 注意:
//   - proto.Message类型会先序列化为字节数组
//   - 其他类型会转换为字符串后存储
func (cm *RedisCommandHash) SetHashValueProtoMsg(hashKey, hashField any, value any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)
	hashFieldStr := fmt.Sprintf("%v", hashField)

	var bytesData []byte
	var err error

	switch v := value.(type) {
	case proto.Message:
		bytesData, err = proto.Marshal(v)
		if err != nil {
			return err
		}
	case []byte:
		bytesData = v
	default:
		bytesData = []byte(fmt.Sprintf("%v", value))
	}

	return client.HSet(ctx, hashKeyStr, hashFieldStr, bytesData).Err()
}

// SetHashExpire 设置Hash的过期时间
// 参数:
//   - hashKey: Hash键名
//   - expirationSeconds: 过期时间（秒）
//
// 返回:
//   - error: 设置失败时返回错误
func (cm *RedisCommandHash) SetHashExpire(hashKey, expirationSeconds any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	seconds, err := strconv.Atoi(fmt.Sprintf("%v", expirationSeconds))
	if err != nil {
		return err
	}

	return client.Expire(ctx, hashKeyStr, time.Duration(seconds)*time.Second).Err()
}

// SetMultiHashValues 批量设置Hash的多个字段值（使用Pipeline）
// 参数:
//   - hashKey: Hash键名
//   - dataList: 字段值映射表，key为uint64类型字段名，value为Protobuf消息指针
//
// 返回:
//   - bool: 批量设置成功返回true，否则返回false
//
// 注意:
//   - 使用Pipeline批量操作，提高性能
//   - 如果某个消息序列化失败，会跳过该字段继续处理其他字段
func (cm *RedisCommandHash) SetMultiHashValues(hashKey any, dataList map[uint64]*proto.Message) bool {
	client := GetClient()
	if client == nil {
		mlog.Error("Redis client is nil")
		return false
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	// 使用Pipeline批量操作
	pipe := client.Pipeline()

	for key, v := range dataList {
		if v != nil {
			bytesData, err := proto.Marshal(*v)
			if err != nil {
				mlog.Error("SetMultiHashValues marshal failed, err:%v", err)
				continue
			}
			pipe.HSet(ctx, hashKeyStr, strconv.FormatUint(key, 10), bytesData)
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		mlog.Error("SetMultiHashValues pipeline exec failed, err:%v", err)
		return false
	}

	return true
}

// GetAllBytesHashValues 获取Hash的所有字段和值（字节格式）
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - map[string][]byte: 字段名到字节数组的映射
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetAllBytesHashValues(hashKey any) (map[string][]byte, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	result, err := client.HGetAll(ctx, hashKeyStr).Result()
	if err != nil {
		return nil, err
	}

	// 转换为字节格式
	byteResult := make(map[string][]byte)
	for k, v := range result {
		byteResult[k] = []byte(v)
	}

	return byteResult, nil
}

// GetAllHashValues 获取Hash的所有字段和值（字节列表格式，键值交替）
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - [][]byte: 字节数组列表，格式为[field1, value1, field2, value2, ...]
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetAllHashValues(hashKey any) ([][]byte, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	result, err := client.HGetAll(ctx, hashKeyStr).Result()
	if err != nil {
		return nil, err
	}

	// 转换为字节列表格式（键值交替）
	var bytesList [][]byte
	for k, v := range result {
		bytesList = append(bytesList, []byte(k))
		bytesList = append(bytesList, []byte(v))
	}

	return bytesList, nil
}

// GetAllHashStringValues 获取Hash的所有字段和值（字符串格式）
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - map[string]string: 字段名到字符串值的映射
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetAllHashStringValues(hashKey any) (map[string]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	return client.HGetAll(ctx, hashKeyStr).Result()
}

// GetMultiHashValues 批量获取Hash多个字段的值（字节格式）
// 参数:
//   - hashKey: Hash键名
//   - hashFields: 要获取的字段列表（uint64类型）
//
// 返回:
//   - [][]byte: 字段值列表，顺序与hashFields对应，不存在的字段返回nil
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetMultiHashValues(hashKey any, hashFields []uint64) ([][]byte, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	// 转换字段列表为字符串
	fields := make([]string, len(hashFields))
	for i, field := range hashFields {
		fields[i] = strconv.FormatUint(field, 10)
	}

	result, err := client.HMGet(ctx, hashKeyStr, fields...).Result()
	if err != nil {
		return nil, err
	}

	// 转换为字节格式
	var bytesList [][]byte
	for _, val := range result {
		if val != nil {
			bytesList = append(bytesList, []byte(fmt.Sprintf("%v", val)))
		} else {
			bytesList = append(bytesList, nil)
		}
	}

	return bytesList, nil
}

// GetMultiHashStringValues 批量获取Hash多个字段的字符串值
// 参数:
//   - hashKey: Hash键名
//   - hashFields: 要获取的字段列表
//
// 返回:
//   - map[string]string: 字段名到字符串值的映射，不存在的字段不包含在结果中
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetMultiHashStringValues(hashKey any, hashFields []any) (map[string]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	// 转换字段列表为字符串
	fields := make([]string, len(hashFields))
	for i, field := range hashFields {
		fields[i] = fmt.Sprintf("%v", field)
	}

	result, err := client.HMGet(ctx, hashKeyStr, fields...).Result()
	if err != nil {
		return nil, err
	}

	// 转换为map格式
	resultMap := make(map[string]string)
	for i, val := range result {
		if val != nil {
			resultMap[fields[i]] = fmt.Sprintf("%v", val)
		}
	}

	return resultMap, nil
}

// GetHashValueString 获取Hash指定字段的字符串值
// 参数:
//   - hashKey: Hash键名
//   - hashField: 字段名
//
// 返回:
//   - string: 字段值，字段不存在时返回空字符串
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashValueString(hashKey, hashField any) (string, error) {
	client := GetClient()
	if client == nil {
		return "", ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)
	hashFieldStr := fmt.Sprintf("%v", hashField)

	result, err := client.HGet(ctx, hashKeyStr, hashFieldStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// GetHashValueInt 获取Hash指定字段的整数值
// 参数:
//   - hashKey: Hash键名
//   - hashField: 字段名
//
// 返回:
//   - int: 字段的整数值
//   - error: 获取失败或值无法转换为整数时返回错误
func (cm *RedisCommandHash) GetHashValueInt(hashKey, hashField any) (int, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)
	hashFieldStr := fmt.Sprintf("%v", hashField)

	return client.HGet(ctx, hashKeyStr, hashFieldStr).Int()
}

// GetHashValueIntExceptNil 获取Hash指定字段的整数值（字段必须存在）
// 参数:
//   - hashKey: Hash键名
//   - hashField: 字段名
//
// 返回:
//   - int: 字段的整数值
//   - error: 字段不存在或获取失败时返回错误
func (cm *RedisCommandHash) GetHashValueIntExceptNil(hashKey, hashField any) (int, error) {
	result, err := cm.GetHashValueInt(hashKey, hashField)
	if err != nil {
		if err == goredis.Nil {
			return 0, errors.New("value is nil")
		}
		return 0, err
	}
	return result, nil
}

// GetHashValueByte 获取Hash指定字段的字节数组值
// 参数:
//   - hashKey: Hash键名
//   - hashField: 字段名
//
// 返回:
//   - []byte: 字段的字节数组值，字段不存在时返回nil
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashValueByte(hashKey, hashField any) ([]byte, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)
	hashFieldStr := fmt.Sprintf("%v", hashField)

	result, err := client.HGet(ctx, hashKeyStr, hashFieldStr).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return []byte(result), nil
}

// IncHashValue Hash字段值自增
// 参数:
//   - hashKey: Hash键名
//   - hashField: 字段名
//   - incNum: 增量值，支持int、int32、int64、uint32类型
//
// 返回:
//   - int: 自增后的值
//   - error: 操作失败时返回错误
//
// 注意:
//   - 如果字段不存在，会先初始化为0再执行自增
func (cm *RedisCommandHash) IncHashValue(hashKey, hashField, incNum any) (int, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)
	hashFieldStr := fmt.Sprintf("%v", hashField)

	// 转换增量值
	var incValue int64
	switch v := incNum.(type) {
	case int:
		incValue = int64(v)
	case int32:
		incValue = int64(v)
	case int64:
		incValue = v
	case uint32:
		incValue = int64(v)
	default:
		return 0, errors.New("invalid increment value type")
	}

	result, err := client.HIncrBy(ctx, hashKeyStr, hashFieldStr, incValue).Result()
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

// SetValueKeyLock 设置键锁（使用SetNX实现）
// 参数:
//   - hashKey: 锁的键名
//   - expiredSec: 锁的过期时间（秒），支持uint32、int、int64类型
//
// 返回:
//   - bool: 获取锁成功返回true，锁已存在返回false
//   - error: 操作失败时返回错误
//
// 注意:
//   - 使用SetNX保证原子性
//   - 锁会自动过期
func (cm *RedisCommandHash) SetValueKeyLock(hashKey, expiredSec any) (bool, error) {
	client := GetClient()
	if client == nil {
		return false, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	var expireDuration time.Duration
	switch v := expiredSec.(type) {
	case uint32:
		expireDuration = time.Duration(v) * time.Second
	case int:
		expireDuration = time.Duration(v) * time.Second
	case int64:
		expireDuration = time.Duration(v) * time.Second
	default:
		return false, errors.New("invalid expiration type")
	}

	result, err := client.SetNX(ctx, hashKeyStr, "locked", expireDuration).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

// GetHashValues 获取Hash的所有字段和值
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - map[string]string: 字段名到字符串值的映射
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashValues(hashKey any) (map[string]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	// 使用HGETALL获取所有字段和值
	result, err := client.HGetAll(ctx, hashKeyStr).Result()
	if err != nil {
		mlog.Warn("[RedisCommandHash] GetHashValues error for key [%s]: %v", hashKeyStr, err)
		return nil, err
	}

	return result, nil
}

// GetHashValuesAsJSON 获取Hash的所有字段和值并转换为JSON字符串
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - string: JSON格式的字符串，Hash为空时返回"{}"
//   - error: 获取或序列化失败时返回错误
func (cm *RedisCommandHash) GetHashValuesAsJSON(hashKey any) (string, error) {
	hashMap, err := cm.GetHashValues(hashKey)
	if err != nil {
		return "", err
	}

	if len(hashMap) == 0 {
		return "{}", nil
	}

	jsonBytes, jsonErr := json.Marshal(hashMap)
	if jsonErr != nil {
		mlog.Warn("[RedisCommandHash] GetHashValuesAsJSON Error marshaling hash to JSON for key [%v]: %v", hashKey, jsonErr)
		return "", jsonErr
	}

	return string(jsonBytes), nil
}

// GetHashValuesWithFilter 获取Hash的所有字段和值并支持自定义过滤
// 参数:
//   - hashKey: Hash键名
//   - fieldFilter: 字段过滤函数，返回true表示保留该字段，nil表示不过滤
//
// 返回:
//   - map[string]string: 过滤后的字段名到字符串值的映射
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashValuesWithFilter(hashKey any, fieldFilter func(field string) bool) (map[string]string, error) {
	allValues, err := cm.GetHashValues(hashKey)
	if err != nil {
		return nil, err
	}

	if fieldFilter == nil {
		return allValues, nil
	}

	filteredMap := make(map[string]string)
	for field, value := range allValues {
		if fieldFilter(field) {
			filteredMap[field] = value
		}
	}

	return filteredMap, nil
}

// GetHashSize 获取Hash的字段数量
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - int64: Hash中的字段数量
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashSize(hashKey any) (int64, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	return client.HLen(ctx, hashKeyStr).Result()
}

// HashExists 检查Hash键是否存在
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - bool: Hash存在返回true，否则返回false
//   - error: 检查失败时返回错误
func (cm *RedisCommandHash) HashExists(hashKey any) (bool, error) {
	client := GetClient()
	if client == nil {
		return false, ErrClientNil
	}

	ctx := context.Background()
	hashKeyStr := fmt.Sprintf("%v", hashKey)

	result, err := client.Exists(ctx, hashKeyStr).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// 使用示例和兼容性方法

// GetHashValuesForAPI 为API接口提供的便捷方法（返回JSON格式）
// 参数:
//   - hashKey: Hash键名
//
// 返回:
//   - string: JSON格式的字符串，Hash不存在或为空时返回"{}"
//   - error: 获取或序列化失败时返回错误
//
// 注意:
//   - 会先检查Hash是否存在
//   - 适合直接用于API响应
func (cm *RedisCommandHash) GetHashValuesForAPI(hashKey any) (string, error) {
	// 检查Hash是否存在
	exists, err := cm.HashExists(hashKey)
	if err != nil {
		return "", err
	}

	if !exists {
		return "{}", nil // 返回空JSON对象
	}

	// 获取JSON格式的Hash值
	return cm.GetHashValuesAsJSON(hashKey)
}

// GetHashValuesWithPrefix 获取Hash中指定前缀的字段
// 参数:
//   - hashKey: Hash键名
//   - prefix: 字段名前缀
//
// 返回:
//   - map[string]string: 匹配前缀的字段名到字符串值的映射
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashValuesWithPrefix(hashKey any, prefix string) (map[string]string, error) {
	return cm.GetHashValuesWithFilter(hashKey, func(field string) bool {
		return len(field) >= len(prefix) && field[:len(prefix)] == prefix
	})
}

// GetHashValuesExcludePrefix 获取Hash中排除指定前缀的字段
// 参数:
//   - hashKey: Hash键名
//   - prefix: 要排除的字段名前缀
//
// 返回:
//   - map[string]string: 不匹配前缀的字段名到字符串值的映射
//   - error: 获取失败时返回错误
func (cm *RedisCommandHash) GetHashValuesExcludePrefix(hashKey any, prefix string) (map[string]string, error) {
	return cm.GetHashValuesWithFilter(hashKey, func(field string) bool {
		return len(field) < len(prefix) || field[:len(prefix)] != prefix
	})
}
