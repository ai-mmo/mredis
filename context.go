package mredis

import (
	"context"
	"time"
)

// DefaultTimeout 默认操作超时时间
const DefaultTimeout = 5 * time.Second

// ContextWithTimeout 创建带指定超时时间的上下文
// 参数:
//   - timeout: 超时时间，如果<=0则使用默认超时时间(5秒)
//
// 返回:
//   - context.Context: 上下文对象
//   - context.CancelFunc: 取消函数，用于提前取消上下文
//
// 注意:
//   - 使用完毕后应调用返回的cancel函数释放资源
//   - 建议使用defer cancel()确保资源释放
func ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return context.WithTimeout(context.Background(), timeout)
}

// ContextWithDefaultTimeout 创建带默认超时时间的上下文
// 返回:
//   - context.Context: 上下文对象
//   - context.CancelFunc: 取消函数，用于提前取消上下文
//
// 注意:
//   - 默认超时时间为5秒
//   - 使用完毕后应调用返回的cancel函数释放资源
//   - 建议使用defer cancel()确保资源释放
func ContextWithDefaultTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}
