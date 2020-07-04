package cache_proxy

import "github.com/treeyh/soc-go-common/library/redis"

func GetRedis() *redis.RedisProxy {
	return redis.GetProxy()
}
