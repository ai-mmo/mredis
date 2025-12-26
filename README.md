# mredis

高性能的 Redis 客户端封装库，基于 go-redis/v9 构建，提供简洁易用的 API。

> **包名**: `mredis`

## 特性

- ✅ 完整的 Redis 数据类型支持（String、Hash、List、ZSet）
- ✅ 连接池管理和健康检查
- ✅ 统一的错误处理
- ✅ 上下文超时控制
- ✅ Protobuf 消息支持
- ✅ 分布式锁实现
- ✅ 批量操作优化
- ✅ 完善的单元测试

## 安装

```bash
go get github.com/ai-mmo/mredis@latest
```

## 快速开始

### 基础配置

```go
package main

import (
    "time"
    "github.com/ai-mmo/mredis"
)

func main() {
    // 配置 Redis 连接
    config := &mredis.RedisConfig{
        Addr:         "localhost:6379",
        Pass:         "",
        Index:        0,
        PoolSize:     10,
        MinIdleConns: 2,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    }

    // 加载配置
    if err := mredis.LoadConfig(config); err != nil {
        panic(err)
    }
    defer mredis.Close()

    // 检查健康状态
    if !mredis.IsHealthy() {
        panic("Redis connection is not healthy")
    }
}
```

### String 操作

```go
// 设置值
err := mredis.StringCommand.SetStringValue("key", "value")

// 获取值
value, err := mredis.StringCommand.GetStringValue("key")

// 自增
count, err := mredis.StringCommand.Incr("counter")

// 分布式锁
success, firstTime, err := mredis.StringCommand.GetLock(
    "lock:key",
    10,    // 过期时间（秒）
    100,   // 重试间隔（毫秒）
    5,     // 最大重试次数
)
if success {
    defer mredis.StringCommand.UnLock("lock:key")
    // 执行业务逻辑
}
```

### Hash 操作

```go
// 设置 Hash 字段
err := mredis.HashCommand.SetHashValue("user:1", "name", "Alice")

// 获取 Hash 字段
name, err := mredis.HashCommand.GetHashValueString("user:1", "name")

// 获取所有字段
allFields, err := mredis.HashCommand.GetHashValues("user:1")

// Hash 字段自增
newValue, err := mredis.HashCommand.IncHashValue("user:1", "score", 10)

// 批量设置（使用 Pipeline）
dataMap := map[uint64]*proto.Message{
    1: &msg1,
    2: &msg2,
}
success := mredis.HashCommand.SetMultiHashValues("data:key", dataMap)
```

### List 操作

```go
// 推入元素
err := mredis.ListCommand.RPush("queue", "item1", "item2")

// 弹出元素
item, err := mredis.ListCommand.LPop("queue")

// 获取列表长度
length, err := mredis.ListCommand.LLen("queue")

// 获取范围内的元素
items, err := mredis.ListCommand.LRange("queue", 0, -1)
```

### ZSet 操作

```go
// 添加成员
success := mredis.ZSetCommand.ZAdd("leaderboard", 100.0, "player1")

// 获取排名
rank, err := mredis.ZSetCommand.GetZSetRank("leaderboard", "player1")

// 获取排行榜
members, ranks, scores, err := mredis.ZSetCommand.GetZSetRankList(
    "leaderboard",
    1,  // 起始排名
    10, // 结束排名
)

// 根据分数范围查询
members, scores, err := mredis.ZSetCommand.GetZSetRangeByScoreWithScores(
    "leaderboard",
    "0",
    "100",
    nil,
)
```

### 高级功能

#### 模式匹配删除

```go
// 删除匹配模式的所有键
deletedKeys, err := mredis.DeleteRedisKeysByPattern("cache:*")
```

#### 连接池统计

```go
stats := mredis.GetPoolStats()
if stats != nil {
    fmt.Printf("Hits: %d, Misses: %d, Timeouts: %d\n",
        stats.Hits, stats.Misses, stats.Timeouts)
}
```

## 版本发布

本项目使用 GitHub Actions 自动发布。

### 使用 release.sh 脚本发布

```bash
# 自动递增补丁版本（推荐用于 bug 修复）
./release.sh init "修复连接池泄漏问题"

# 自动递增次版本（用于新功能）
./release.sh minor "新增分布式锁功能"

# 自动递增主版本（用于重大更新）
./release.sh major "重构 API 接口"

# 手动指定版本号
./release.sh v1.2.3 "发布 1.2.3 版本"
```

### 发布流程

1. 确保所有更改已提交
2. 运行 `./release.sh` 脚本
3. 脚本会自动：
   - 运行测试
   - 检查代码格式
   - 创建 Git 标签
   - 推送到 GitHub
4. GitHub Actions 会自动：
   - 运行 CI 测试
   - 创建 GitHub Release
   - 生成变更日志

## 错误处理

库提供了统一的错误常量：

```go
var (
    ErrClientNil     = errors.New("redis client is nil")
    ErrPoolClosed    = errors.New("connection pool is closed")
    ErrInvalidConfig = errors.New("invalid redis configuration")
    ErrKeyNotFound   = errors.New("key not found")
    ErrInvalidType   = errors.New("invalid data type")
    ErrTimeout       = errors.New("operation timeout")
    ErrEmptyKey      = errors.New("empty key")
)
```

## 性能优化建议

1. **连接池配置**：根据并发量调整 `PoolSize` 和 `MinIdleConns`
2. **超时设置**：合理设置 `DialTimeout`、`ReadTimeout` 和 `WriteTimeout`
3. **批量操作**：使用 Pipeline 进行批量操作以减少网络往返
4. **键命名**：使用有意义的前缀，便于管理和清理

## 测试

```bash
# 运行所有测试
go test -v ./...

# 运行测试并查看覆盖率
go test -v -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 依赖

- Go 1.25+
- github.com/redis/go-redis/v9
- google.golang.org/protobuf
- github.com/ai-mmo/mlog

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

查看 [Releases](https://github.com/ai-mmo/mredis/releases) 页面了解详细的版本更新信息。