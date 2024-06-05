package store

import (
	"github.com/patrickmn/go-cache"
	"time"
)

//内存缓存
var C *cache.Cache

func InitMemoryCache() {
	C = cache.New(time.Second, time.Minute)
}
