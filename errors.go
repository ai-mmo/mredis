package mredis

import "errors"

// 定义常见错误
var (
	// ErrClientNil Redis客户端未初始化
	ErrClientNil = errors.New("redis client is nil")

	// ErrPoolClosed 连接池已关闭
	ErrPoolClosed = errors.New("connection pool is closed")

	// ErrInvalidConfig 无效的配置
	ErrInvalidConfig = errors.New("invalid redis configuration")

	// ErrKeyNotFound 键不存在
	ErrKeyNotFound = errors.New("key not found")

	// ErrInvalidType 无效的数据类型
	ErrInvalidType = errors.New("invalid data type")

	// ErrTimeout 操作超时
	ErrTimeout = errors.New("operation timeout")

	// ErrEmptyKey 空键
	ErrEmptyKey = errors.New("empty key")
)
