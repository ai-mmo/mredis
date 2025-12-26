package mredis

import "time"

type RedisConfig struct {
	// 基础连接配置
	Addr  string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Pass  string `mapstructure:"pass" json:"pass" yaml:"pass"`
	Index int    `mapstructure:"index" json:"index" yaml:"index"`

	// 连接池配置
	PoolSize     int `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	MinIdleConns int `mapstructure:"min_idle_conns" json:"min_idle_conns" yaml:"min_idle_conns"`
	MaxIdleConns int `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`

	// 超时配置
	DialTimeout  time.Duration `mapstructure:"dial_timeout" json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout" json:"pool_timeout" yaml:"pool_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" json:"idle_timeout" yaml:"idle_timeout"`

	// 重试配置
	MaxRetries      int           `mapstructure:"max_retries" json:"max_retries" yaml:"max_retries"`
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff" json:"min_retry_backoff" yaml:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff" json:"max_retry_backoff" yaml:"max_retry_backoff"`
}
