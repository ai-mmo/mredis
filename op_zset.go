package mredis

import (
	"context"
	"errors"
	"fmt"
	"mlog"

	goredis "github.com/redis/go-redis/v9"
)

type RedisCommandZSet struct{}

// ZAdd 向有序集合添加一个成员或更新已存在成员的分数
// 参数:
//   - key: 有序集合键名
//   - score: 成员的分数
//   - member: 成员值
//
// 返回:
//   - bool: 添加成功返回true，失败返回false
func (cm *RedisCommandZSet) ZAdd(key any, score float64, member any) bool {
	client := GetClient()
	if client == nil {
		mlog.Error("Redis client is nil")
		return false
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	z := goredis.Z{Score: score, Member: member}
	err := client.ZAdd(ctx, keyStr, z).Err()
	if err != nil {
		mlog.Error("ZAdd failed: %v", err)
		return false
	}
	return true
}

// GetZSetCount 获取有序集合的成员数量
// 参数:
//   - key: 有序集合键名
//
// 返回:
//   - uint32: 成员数量
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) GetZSetCount(key any) (uint32, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	result, err := client.ZCard(ctx, keyStr).Result()
	if err != nil {
		mlog.Error("[RedisCommandZSet] GetZSetCount zcard error,key:%s err:%v", keyStr, err)
		return 0, err
	}

	return uint32(result), nil
}

// GetZSetRangeByScore 根据分数范围获取有序集合成员
// 参数:
//   - key: 有序集合键名
//   - min: 最小分数（字符串格式，支持"-inf"表示负无穷）
//   - max: 最大分数（字符串格式，支持"+inf"表示正无穷）
//
// 返回:
//   - []string: 指定分数范围内的成员列表
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) GetZSetRangeByScore(key any, min, max string) ([]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	opt := &goredis.ZRangeBy{
		Min: min,
		Max: max,
	}

	result, err := client.ZRangeByScore(ctx, keyStr, opt).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

// GetZSetRankList 返回有序集中指定排名区间内的成员（分数从高到低）
// 参数:
//   - key: 有序集合键名
//   - rankMin: 最小排名（从1开始）
//   - rankMax: 最大排名
//
// 返回:
//   - []string: 成员ID列表
//   - []uint64: 排名列表
//   - [][]uint64: 分数列表（每个元素包含两个uint64值）
//   - error: 操作失败时返回错误
//
// 注意:
//   - 排名从1开始计数
//   - 分数从高到低排序
func (cm *RedisCommandZSet) GetZSetRankList(key any, rankMin, rankMax uint64) ([]string, []uint64, [][]uint64, error) {
	client := GetClient()
	if client == nil {
		return nil, nil, nil, ErrClientNil
	}

	if rankMin > rankMax {
		return nil, nil, nil, errors.New("rankMin > rankMax")
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	// 使用ZREVRANGE获取逆序排列的成员和分数
	result, err := client.ZRevRangeWithScores(ctx, keyStr, int64(rankMin-1), int64(rankMax-1)).Result()
	if err != nil {
		return nil, nil, nil, err
	}

	size := len(result)
	if size == 0 {
		return nil, nil, nil, errors.New("rank data empty")
	}

	startRank := rankMin
	slcUniqueId := make([]string, 0, size)
	rankList := make([]uint64, 0, size)
	slcSortValue := make([][]uint64, 0, size)

	for _, z := range result {
		// 添加成员ID
		memberStr := fmt.Sprintf("%v", z.Member)
		slcUniqueId = append(slcUniqueId, memberStr)

		// 添加排名
		rankList = append(rankList, startRank)
		startRank++

		// 处理分数，转换为uint64并拆分
		scoreUint64 := uint64(z.Score)
		first, second := SplitFromScore(scoreUint64)
		slcSortValue = append(slcSortValue, []uint64{first, second})
	}

	return slcUniqueId, rankList, slcSortValue, nil
}

// GetZSetRangeByScoreWithScores 根据分数范围获取有序集合成员和分数
// 参数:
//   - key: 有序集合键名
//   - min: 最小分数（字符串格式）
//   - max: 最大分数（字符串格式）
//   - matched: 用于记录已匹配成员的map，可为nil
//
// 返回:
//   - []string: 成员列表
//   - []float64: 对应的分数列表
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) GetZSetRangeByScoreWithScores(key any, min, max string, matched map[string]bool) ([]string, []float64, error) {
	client := GetClient()
	if client == nil {
		return nil, nil, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	if matched == nil {
		matched = make(map[string]bool)
	}

	opt := &goredis.ZRangeBy{
		Min: min,
		Max: max,
	}

	result, err := client.ZRangeByScoreWithScores(ctx, keyStr, opt).Result()
	if err != nil {
		return nil, nil, err
	}

	if len(result) == 0 {
		return nil, nil, nil
	}

	members := make([]string, len(result))
	scores := make([]float64, len(result))

	for i, z := range result {
		memberStr := fmt.Sprintf("%v", z.Member)
		members[i] = memberStr
		scores[i] = z.Score
		matched[memberStr] = true
	}

	return members, scores, nil
}

// GetZSetRangeByScoreWithScoresLimit 根据分数范围获取有序集合成员和分数（自动判断排序方向）
// 参数:
//   - key: 有序集合键名
//   - param1: 第一个分数参数
//   - param2: 第二个分数参数
//
// 返回:
//   - []string: 成员列表
//   - []float64: 对应的分数列表
//   - error: 操作失败时返回错误
//
// 注意:
//   - 如果param1<=param2，按升序返回
//   - 如果param1>param2，按降序返回
func (cm *RedisCommandZSet) GetZSetRangeByScoreWithScoresLimit(key any, param1, param2 float64) ([]string, []float64, error) {
	client := GetClient()
	if client == nil {
		return nil, nil, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	var opt *goredis.ZRangeBy
	if param1 <= param2 {
		opt = &goredis.ZRangeBy{
			Min: fmt.Sprintf("%f", param1),
			Max: fmt.Sprintf("%f", param2),
		}
	} else {
		opt = &goredis.ZRangeBy{
			Min: fmt.Sprintf("%f", param2),
			Max: fmt.Sprintf("%f", param1),
		}
	}

	var result []goredis.Z
	var err error

	if param1 <= param2 {
		result, err = client.ZRangeByScoreWithScores(ctx, keyStr, opt).Result()
	} else {
		result, err = client.ZRevRangeByScoreWithScores(ctx, keyStr, opt).Result()
	}

	if err != nil {
		return nil, nil, err
	}

	if len(result) == 0 {
		return nil, nil, nil
	}

	members := make([]string, len(result))
	scores := make([]float64, len(result))

	for i, z := range result {
		members[i] = fmt.Sprintf("%v", z.Member)
		scores[i] = z.Score
	}

	return members, scores, nil
}

// ZRem 移除有序集合中的一个或多个成员
// 参数:
//   - key: 有序集合键名
//   - members: 要移除的成员，可变参数
//
// 返回:
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) ZRem(key any, members ...any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.ZRem(ctx, keyStr, members...).Err()
}

// ZRemRangeByScore 移除有序集合中指定分数区间的所有成员
// 参数:
//   - key: 有序集合键名
//   - min: 最小分数（字符串格式）
//   - max: 最大分数（字符串格式）
//
// 返回:
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) ZRemRangeByScore(key any, min, max string) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.ZRemRangeByScore(ctx, keyStr, min, max).Err()
}

// ZScore 获取有序集合中指定成员的分数值
// 参数:
//   - key: 有序集合键名
//   - member: 成员值
//
// 返回:
//   - float64: 成员的分数
//   - error: 成员不存在或操作失败时返回错误
func (cm *RedisCommandZSet) ZScore(key, member any) (float64, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	memberStr := fmt.Sprintf("%v", member)

	return client.ZScore(ctx, keyStr, memberStr).Result()
}

// ZRank 获取有序集合中指定成员的排名（按分数从低到高）
// 参数:
//   - key: 有序集合键名
//   - member: 成员值
//
// 返回:
//   - int64: 成员的排名（从0开始）
//   - error: 成员不存在或操作失败时返回错误
func (cm *RedisCommandZSet) ZRank(key, member any) (int64, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	memberStr := fmt.Sprintf("%v", member)

	return client.ZRank(ctx, keyStr, memberStr).Result()
}

// ZRevRank 获取有序集合中指定成员的逆序排名（按分数从高到低）
// 参数:
//   - key: 有序集合键名
//   - member: 成员值
//
// 返回:
//   - int64: 成员的逆序排名（从0开始）
//   - error: 成员不存在或操作失败时返回错误
func (cm *RedisCommandZSet) ZRevRank(key, member any) (int64, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	memberStr := fmt.Sprintf("%v", member)

	return client.ZRevRank(ctx, keyStr, memberStr).Result()
}

// ZRange 获取有序集合中指定区间内的成员（按分数从低到高）
// 参数:
//   - key: 有序集合键名
//   - start: 起始索引（从0开始）
//   - stop: 结束索引（-1表示最后一个）
//
// 返回:
//   - []string: 成员列表
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) ZRange(key any, start, stop int64) ([]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.ZRange(ctx, keyStr, start, stop).Result()
}

// ZRevRange 获取有序集合中指定区间内的成员（按分数从高到低）
// 参数:
//   - key: 有序集合键名
//   - start: 起始索引（从0开始）
//   - stop: 结束索引（-1表示最后一个）
//
// 返回:
//   - []string: 成员列表
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) ZRevRange(key any, start, stop int64) ([]string, error) {
	client := GetClient()
	if client == nil {
		return nil, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.ZRevRange(ctx, keyStr, start, stop).Result()
}

// ZIncrBy 为有序集合中指定成员的分数加上增量
// 参数:
//   - key: 有序集合键名
//   - increment: 增量值（可以是负数）
//   - member: 成员值
//
// 返回:
//   - float64: 增加后的分数
//   - error: 操作失败时返回错误
//
// 注意:
//   - 如果成员不存在，会先创建成员并设置分数为increment
func (cm *RedisCommandZSet) ZIncrBy(key any, increment float64, member any) (float64, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	memberStr := fmt.Sprintf("%v", member)

	return client.ZIncrBy(ctx, keyStr, increment, memberStr).Result()
}

// DelKey 删除整个有序集合键
// 参数:
//   - key: 有序集合键名
//
// 返回:
//   - error: 删除失败时返回错误
func (cm *RedisCommandZSet) DelKey(key any) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	return client.Del(ctx, keyStr).Err()
}

// ZSetAddValue 向有序集合添加成员（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - member: 成员值
//   - score: 成员的分数（uint64类型）
//
// 返回:
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) ZSetAddValue(key any, member any, score uint64) error {
	client := GetClient()
	if client == nil {
		return ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)

	z := goredis.Z{Score: float64(score), Member: member}
	return client.ZAdd(ctx, keyStr, z).Err()
}

// GetZSetValueByRange 根据排名范围获取有序集合成员（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - start: 起始排名（从1开始）
//   - stop: 结束排名
//
// 返回:
//   - []string: 成员列表（按分数从高到低）
//   - error: 操作失败时返回错误
//
// 注意:
//   - 排名从1开始，会自动转换为0基索引
func (cm *RedisCommandZSet) GetZSetValueByRange(key any, start, stop int64) ([]string, error) {
	return cm.ZRevRange(key, start-1, stop-1) // 转换为0基索引
}

// GetZSetRandomValueByRange 根据分数范围获取有序集合成员（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - min: 最小分数
//   - max: 最大分数
//   - matched: 用于记录已匹配成员的map
//   - count: 未使用的参数（保留用于兼容性）
//
// 返回:
//   - []string: 成员列表
//   - []float64: 对应的分数列表
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) GetZSetRandomValueByRange(key any, min, max float64, matched map[string]bool, count int) ([]string, []float64, error) {
	minStr := fmt.Sprintf("%f", min)
	maxStr := fmt.Sprintf("%f", max)
	return cm.GetZSetRangeByScoreWithScores(key, minStr, maxStr, matched)
}

// GetZSetRangeByValue 根据分数范围获取有序集合成员（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - min: 最小分数
//   - max: 最大分数
//   - matched: 用于记录已匹配成员的map
//   - count: 未使用的参数（保留用于兼容性）
//
// 返回:
//   - []string: 成员列表
//   - []float64: 对应的分数列表
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) GetZSetRangeByValue(key any, min, max float64, matched map[string]bool, count int) ([]string, []float64, error) {
	minStr := fmt.Sprintf("%f", min)
	maxStr := fmt.Sprintf("%f", max)
	return cm.GetZSetRangeByScoreWithScores(key, minStr, maxStr, matched)
}

// GetZSetValueByRank 根据排名获取有序集合成员（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - start: 起始排名（从1开始）
//   - stop: 结束排名
//
// 返回:
//   - []string: 成员列表（按分数从高到低）
//   - error: 操作失败时返回错误
//
// 注意:
//   - 排名从1开始，会自动转换为0基索引
func (cm *RedisCommandZSet) GetZSetValueByRank(key any, start, stop int64) ([]string, error) {
	return cm.ZRevRange(key, start-1, stop-1) // 转换为0基索引
}

// GetZSetRank 获取有序集合中指定成员的排名（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - member: 成员值
//
// 返回:
//   - int64: 成员的排名（从1开始，按分数从高到低）
//   - error: 成员不存在或操作失败时返回错误
//
// 注意:
//   - 返回的排名从1开始（与Redis的0基索引不同）
func (cm *RedisCommandZSet) GetZSetRank(key, member any) (int64, error) {
	rank, err := cm.ZRevRank(key, member)
	if err != nil {
		return 0, err
	}
	return rank + 1, nil // 转换为1基索引
}

// GetZSetRankByScore 根据分数获取排名（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - score: 分数值
//
// 返回:
//   - int64: 分数大于等于指定分数的成员数量
//   - error: 操作失败时返回错误
//
// 注意:
//   - 返回的是分数大于等于指定分数的成员数量，而非具体排名
func (cm *RedisCommandZSet) GetZSetRankByScore(key any, score float64) (int64, error) {
	client := GetClient()
	if client == nil {
		return 0, ErrClientNil
	}

	ctx := context.Background()
	keyStr := fmt.Sprintf("%v", key)
	scoreStr := fmt.Sprintf("%f", score)

	// 使用ZCOUNT获取分数大于等于指定分数的成员数量
	count, err := client.ZCount(ctx, keyStr, scoreStr, "+inf").Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetZSetRankAndScore 获取有序集合中指定成员的排名和分数（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - member: 成员值
//
// 返回:
//   - int64: 成员的排名（从1开始）
//   - float64: 成员的分数（float64格式）
//   - uint64: 成员的分数（uint64格式）
//   - error: 成员不存在或操作失败时返回错误
func (cm *RedisCommandZSet) GetZSetRankAndScore(key, member any) (int64, float64, uint64, error) {
	rank, err := cm.GetZSetRank(key, member)
	if err != nil {
		return 0, 0, 0, err
	}

	score, err := cm.ZScore(key, member)
	if err != nil {
		return 0, 0, 0, err
	}

	scoreUint64 := uint64(score)
	return rank, score, scoreUint64, nil
}

// RemZSetElement 移除有序集合中的指定成员（兼容旧API）
// 参数:
//   - key: 有序集合键名
//   - member: 要移除的成员值
//
// 返回:
//   - error: 操作失败时返回错误
func (cm *RedisCommandZSet) RemZSetElement(key, member any) error {
	return cm.ZRem(key, member)
}

// ResetZSet 清空有序集合（兼容旧API）
// 参数:
//   - key: 有序集合键名
//
// 返回:
//   - error: 操作失败时返回错误
//
// 注意:
//   - 此操作会删除整个键
func (cm *RedisCommandZSet) ResetZSet(key any) error {
	return cm.DelKey(key)
}

// DelZSetKey 删除有序集合键（兼容旧API）
// 参数:
//   - key: 有序集合键名
//
// 返回:
//   - error: 删除失败时返回错误
func (cm *RedisCommandZSet) DelZSetKey(key any) error {
	return cm.DelKey(key)
}
