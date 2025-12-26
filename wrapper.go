package mredis

import (
	"mlog"

	goredis "github.com/redis/go-redis/v9"
)

var (
	gPoolM        *Pool
	HashCommand   *RedisCommandHash
	ZSetCommand   *RedisCommandZSet
	StringCommand *RedisCommandString
	ListCommand   *RedisCommandList
)

func init() {
	gPoolM = NewPool()
	HashCommand = new(RedisCommandHash)
	ZSetCommand = new(RedisCommandZSet)
	StringCommand = new(RedisCommandString)
	ListCommand = new(RedisCommandList)
}

// LoadConfig 加载Redis配置并初始化连接池
// 参数:
//   - cfg: Redis配置对象，包含连接地址、密码、数据库索引等信息
//
// 返回:
//   - error: 配置加载失败时返回错误，成功返回nil
//
// 注意:
//   - 此函数会关闭旧连接并创建新连接
//   - 配置中的Addr字段为必填项
//   - 其他参数如未设置将使用默认值
func LoadConfig(cfg *RedisConfig) error {
	mlog.Info("LoadConfig:%v", cfg)
	return gPoolM.LoadConfig(cfg)
}

// Execute 执行任意Redis命令
// 参数:
//   - command: Redis命令名称（如"GET"、"SET"、"HGET"等）
//   - args: 命令参数，可变参数列表
//
// 返回:
//   - result: Redis命令执行结果对象
//   - err: 执行失败时返回错误
//
// 示例:
//
//	result, err := Execute("SET", "key", "value")
//	result, err := Execute("GET", "key")
func Execute(command string, args ...any) (result *goredis.Cmd, err error) {
	// 获取Redis客户端
	client := gPoolM.GetClient()
	if client == nil {
		mlog.Error("[Execute] Redis client is nil")
		return nil, ErrClientNil
	}

	// 构建命令参数
	cmdArgs := make([]interface{}, 0, len(args)+1)
	cmdArgs = append(cmdArgs, command)
	cmdArgs = append(cmdArgs, args...)

	// 执行命令
	ctx, cancel := ContextWithDefaultTimeout()
	defer cancel()

	result = client.Do(ctx, cmdArgs...)
	if result.Err() != nil {
		mlog.Error("[Execute] Redis command execution error: %v", result.Err())
		return nil, result.Err()
	}

	// 返回结果
	mlog.Debug("[Execute] Command [%s] executed args:%v with result: %v", command, args, result)
	return result, nil
}

// GetClient 获取Redis客户端实例
// 返回:
//   - *goredis.Client: Redis客户端对象，如果连接池未初始化或已关闭则返回nil
//
// 注意:
//   - 返回的客户端由go-redis内置连接池管理，无需手动释放
//   - 使用前应检查返回值是否为nil
func GetClient() *goredis.Client {
	return gPoolM.GetClient()
}

// Close 关闭Redis连接池
// 注意:
//   - 此操作会关闭所有活动连接
//   - 关闭后需要重新调用LoadConfig才能继续使用
//   - 建议在应用程序退出时调用
func Close() {
	gPoolM.Close()
}

// IsHealthy 检查Redis连接池健康状态
// 返回:
//   - bool: 连接池健康返回true，否则返回false
//
// 注意:
//   - 此方法会执行PING命令测试连接
//   - 超时时间为2秒
func IsHealthy() bool {
	return gPoolM.IsHealthy()
}

// GetPoolStats 获取Redis连接池统计信息
// 返回:
//   - *goredis.PoolStats: 连接池统计信息，包含活动连接数、空闲连接数等
//   - 如果连接池未初始化则返回nil
//
// 统计信息包括:
//   - Hits: 命中次数
//   - Misses: 未命中次数
//   - Timeouts: 超时次数
//   - TotalConns: 总连接数
//   - IdleConns: 空闲连接数
//   - StaleConns: 过期连接数
func GetPoolStats() *goredis.PoolStats {
	return gPoolM.GetPoolStats()
}
