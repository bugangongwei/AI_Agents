package outfit_recommender

import (
	"hash/fnv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 分段锁
type ShardedRateLimiter struct {
	shards []*RateLimiterSegment
	mask   uint64 // 掩码，位运算比数学运算性能更强
}

func NewShardedRateLimiter(shardNum uint64) *ShardedRateLimiter {
	l := &ShardedRateLimiter{
		shards: make([]*RateLimiterSegment, shardNum),
		mask:   shardNum - 1,
	}

	for i := 0; i < int(shardNum); i++ {
		l.shards[i] = newRateLimitSegment()
	}

	return l
}

func (l *ShardedRateLimiter) getShard(key string) *RateLimiterSegment {
	h := fnv.New64a()
	h.Write([]byte(key))
	shardIdx := h.Sum64() & l.mask
	return l.shards[shardIdx]
}

func (l *ShardedRateLimiter) Allow(key string) bool {
	segment := l.getShard(key)

	return segment.Allow(key)
}

func (l *ShardedRateLimiter) SetBucket(key string, rate float64, cap uint64) {
	segment := l.getShard(key)

	newBucket := &TokenBucket{
		capacity: cap,
		rate:     rate,
	}

	segment.rwmu.Lock()
	defer segment.rwmu.Unlock()

	segment.buckets[key] = newBucket
}

// 不同纬度的令牌桶
type RateLimiterSegment struct {
	rwmu    sync.RWMutex            // 保护buckets的读写并发安全
	buckets map[string]*TokenBucket // key是不同的限流维度，比如IP/UserID等
}

func newRateLimitSegment() *RateLimiterSegment {
	return &RateLimiterSegment{
		rwmu:    sync.RWMutex{},
		buckets: make(map[string]*TokenBucket),
	}
}

func (s *RateLimiterSegment) Allow(key string) bool {
	s.rwmu.RLock()
	bucket, ok := s.buckets[key]
	if !ok {
		// fmt.Println("key:", key, " not exist, allowed!!!")
		// return true
		bucket = &TokenBucket{
			rate:     0.001,
			capacity: 1000,
		}
	}
	s.rwmu.RUnlock()

	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	now := time.Now()
	// 如果时钟拨乱，则重置时钟
	if now.Before(bucket.lastTime) {
		bucket.lastTime = now
		bucket.current = 0
	}
	elapsed := now.Sub(bucket.lastTime).Seconds()
	newTokens := elapsed * bucket.rate
	if newTokens > 0 {
		bucket.current = min(bucket.capacity, bucket.current+uint64(newTokens))
		bucket.lastTime = now
	}

	if bucket.current > 0 {
		bucket.current--
		return true
	}

	return false
}

// 单独的令牌桶
type TokenBucket struct {
	capacity uint64    // 桶最大容量
	current  uint64    // 当前桶中令牌个数
	rate     float64   // 每秒放回桶中的令牌数
	lastTime time.Time // 上次拿token时间
}

// API网关限流配置
type GatewayRateLimiter struct {
	limiter *ShardedRateLimiter
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		g := &GatewayRateLimiter{
			limiter: NewShardedRateLimiter(16),
		}

		allowd := true

		if userID := c.GetHeader("X-User-ID"); userID != "" {
			allowd = allowd && g.limiter.Allow(userID)
		} else if apiKey := c.GetHeader("X-API-Key"); apiKey != "" {
			allowd = allowd && g.limiter.Allow(apiKey)
		} else {
			allowd = allowd && g.limiter.Allow(c.ClientIP())
		}

		if !allowd {
			c.AbortWithStatusJSON(429, gin.H{
				"err": "requests limited",
			})
			return
		}

		c.Next()

	}
}
