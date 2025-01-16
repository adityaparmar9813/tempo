package cache

import (
	"errors"
	"time"

	"github.com/adityaparmar9813/tempo/internal/eviction"
)

type Item struct {
    Value      interface{}
    ExpiresAt  time.Time
    LastAccess time.Time
}

type EvictionPolicy interface {
    Put(key string, value interface{})
    Get(key string) (interface{}, bool)
    Remove(key string) bool
    Clear()
    Len() int
}

type EvictionAlgorithm string

const (
    LRU EvictionAlgorithm = "LRU"
    LFU EvictionAlgorithm = "LFU"
)

type TempoCache struct {
    policy     EvictionPolicy
    defaultTTL time.Duration
}

type CacheConfig struct {
    EvictionAlgo EvictionAlgorithm
    MaxItems     int
    DefaultTTL   time.Duration
}

func NewTempoCache(config CacheConfig) (*TempoCache, error) {
    if config.MaxItems <= 0 {
        return nil, errors.New("maximum items must be greater than 0")
    }

    var policy EvictionPolicy
    switch config.EvictionAlgo {
    case LRU:
        policy = eviction.NewLRU(config.MaxItems)
    case LFU:
        policy = eviction.NewLFU(config.MaxItems)
    default:
        return nil, errors.New("unsupported eviction algorithm")
    }

    return &TempoCache{
        policy:     policy,
        defaultTTL: config.DefaultTTL,
    }, nil
}

func (c *TempoCache) Set(key string, value interface{}, ttl *time.Duration) error {
    if key == "" {
        return errors.New("key cannot be empty")
    }

    expirationTime := time.Now().Add(c.defaultTTL)
    if ttl != nil {
        expirationTime = time.Now().Add(*ttl)
    }

    item := &Item{
        Value:      value,
        ExpiresAt:  expirationTime,
        LastAccess: time.Now(),
    }

    c.policy.Put(key, item)
    return nil
}

func (c *TempoCache) Get(key string) (interface{}, error) {
    if key == "" {
        return nil, errors.New("key cannot be empty")
    }

    value, exists := c.policy.Get(key)
    if !exists {
        return nil, errors.New("key not found")
    }

    item, ok := value.(*Item)
    if !ok {
        return nil, errors.New("invalid cache item")
    }

    if time.Now().After(item.ExpiresAt) {
        c.policy.Remove(key)
        return nil, errors.New("key expired")
    }

    item.LastAccess = time.Now()
    return item.Value, nil
}

func (c *TempoCache) Delete(key string) error {
    if key == "" {
        return errors.New("key cannot be empty")
    }

    if !c.policy.Remove(key) {
        return errors.New("key not found")
    }
    return nil
}

func (c *TempoCache) Clear() {
    c.policy.Clear()
}

func (c *TempoCache) StartCleanup(interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for range ticker.C {
            c.cleanupExpired()
        }
    }()
}

func (c *TempoCache) cleanupExpired() {
    // now := time.Now()
	
    
    // Get all items and check expiration
    // Note: This would require adding a method to iterate over items in the EvictionPolicy interface
    // For now, this is a placeholder for the cleanup logic
}