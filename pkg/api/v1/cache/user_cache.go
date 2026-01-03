package cache

import (
	"codim/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	userCacheKeyPrefix = "api:user:"
	userLockKeyPrefix  = "api:lock:user:"
	userCacheTTL       = 15 * time.Minute
	lockTTL            = 5 * time.Second
)

type UserCache struct {
	redis  *redis.Client
	db     *db.Queries
	logger Logger
}

type Logger interface {
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func NewUserCache(redis *redis.Client, db *db.Queries, logger Logger) *UserCache {
	return &UserCache{
		redis:  redis,
		db:     db,
		logger: logger,
	}
}

// GetUser retrieves a user from cache or database with distributed locking
func (c *UserCache) GetUser(ctx context.Context, userUUID uuid.UUID) (db.User, error) {
	cacheKey := c.userCacheKey(userUUID)
	lockKey := c.userLockKey(userUUID)

	cached, err := c.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user db.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return user, nil
		}

		c.logger.Warnf("Failed to unmarshal cached user %s: %v", userUUID, err)
	}

	lockAcquired, err := c.acquireLock(ctx, lockKey)
	if err != nil {
		c.logger.Errorf("Failed to acquire lock for user %s: %v", userUUID, err)
		// Continue without lock - might cause duplicate queries but won't block
	}

	// If lock was acquired, we're responsible for fetching
	deadline := time.Now().Add(500 * time.Millisecond)
	for !lockAcquired && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
		// Try cache again - another request might have populated it
		cached, err := c.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var user db.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				return user, nil
			}
		}
		// If still not in cache after waiting, we'll fetch from DB
		// This handles the case where the lock holder failed to cache
	}

	user, err := c.db.GetUser(ctx, userUUID)
	if err != nil {
		if lockAcquired {
			c.releaseLock(ctx, lockKey)
		}

		return db.User{}, err
	}

	if err := c.setUserCache(ctx, cacheKey, user); err != nil {
		c.logger.Warnf("Failed to cache user %s: %v", userUUID, err)
	}

	// Release lock if we acquired it
	if lockAcquired {
		c.releaseLock(ctx, lockKey)
	}

	return user, nil
}

// InvalidateUser removes a user from cache
func (c *UserCache) InvalidateUser(ctx context.Context, userUUID uuid.UUID) error {
	cacheKey := c.userCacheKey(userUUID)
	return c.redis.Del(ctx, cacheKey).Err()
}

// acquireLock attempts to acquire a distributed lock
func (c *UserCache) acquireLock(ctx context.Context, lockKey string) (bool, error) {
	result, err := c.redis.SetNX(ctx, lockKey, "1", lockTTL).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

// releaseLock releases a distributed lock
func (c *UserCache) releaseLock(ctx context.Context, lockKey string) error {
	return c.redis.Del(ctx, lockKey).Err()
}

// setUserCache stores a user in cache
func (c *UserCache) setUserCache(ctx context.Context, key string, user db.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, userCacheTTL).Err()
}

func (c *UserCache) userCacheKey(userUUID uuid.UUID) string {
	return fmt.Sprintf("%s%s", userCacheKeyPrefix, userUUID.String())
}

func (c *UserCache) userLockKey(userUUID uuid.UUID) string {
	return fmt.Sprintf("%s%s", userLockKeyPrefix, userUUID.String())
}
