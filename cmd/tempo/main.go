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

	for i := 0; i < 10; i++ {
		tempoCache.Set("key", i, nil)
	}
	tempoCache.Set("key5", "yummy5", nil)
	tempoCache.Clear()
	// tempoCache.Delete("key3")

	value, err := tempoCache.Get("key")
	if err != nil {
		fmt.Println("Error getting value:", err)
		return
	}
	fmt.Println(value)
}
