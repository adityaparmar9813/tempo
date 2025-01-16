package main

import (
	"fmt"
	"time"

	cache "github.com/adityaparmar9813/tempo/internal"
)

func main() {
	cacheConfig := cache.CacheConfig{
		EvictionAlgo: cache.LRU,
		MaxItems:     5,
		DefaultTTL:   10 * time.Second,
	}

	tempoCache, err := cache.NewTempoCache(cacheConfig)

	if err != nil {
		fmt.Println("Error creating cache:", err)
		return
	}

	tempoCache.Set("key", "value", nil)
	fmt.Println(tempoCache.Get("key"))
}
