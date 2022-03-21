package reminder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RedisCacheService(t *testing.T) {
	var err error
	var c redisService
	testKey := "test_key"

	expiredN := 0

	t.Run("connect_to_redis", func(t *testing.T) {
		c, err = NewRedisService(RedisServiceOptions{
			Url:    "redis://redis:6379",
			Prefix: "test",
			Exp:    time.Second,
			// Hashed: true,
		})

		c.HandleExpire(func(key string) {
			if key == testKey {
				expiredN++
			}
		})

		assert.NoError(t, err)
		assert.NotNil(t, c)
	})

	t.Run("get_not_in_cache", func(t *testing.T) {
		_, err = c.Get(testKey)
		assert.Error(t, err)
	})

	t.Run("set_in_cache", func(t *testing.T) {
		err = c.Set(testKey, "string_value")
		assert.NoError(t, err)
	})

	t.Run("get_in_cache", func(t *testing.T) {
		v, err := c.Get(testKey)
		assert.NoError(t, err)
		assert.Equal(t, "string_value", v)
	})

	t.Run("unset_and_get_not_in_cache", func(t *testing.T) {
		err = c.Unset(testKey)
		assert.NoError(t, err)

		_, err = c.Get(testKey)
		assert.Error(t, err)
	})

	t.Run("set_in_cache_wait_expire", func(t *testing.T) {
		err = c.Set(testKey, "string_value")
		assert.NoError(t, err)

		time.Sleep(time.Second * 2)

		_, err = c.Get(testKey)
		assert.Error(t, err)

		assert.Equal(t, 1, expiredN)
	})

}
