module github.com/ai-mmo/mredis

go 1.25

replace mlog => github.com/ai-mmo/mlog v0.0.19

require (
	github.com/redis/go-redis/v9 v9.12.1
	google.golang.org/protobuf v1.36.7
	mlog v0.0.0-00010101000000-000000000000
)

require (
	github.com/ai-mmo/lumberjack v0.0.5 // indirect
	github.com/alicebob/miniredis/v2 v2.35.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
)
