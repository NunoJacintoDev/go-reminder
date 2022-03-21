package reminder

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

// RedisServiceOptions options for redisService
type RedisServiceOptions struct {
	Url    string
	Prefix string
	Exp    time.Duration
	Hashed bool
}

// NewRedisServiceOptions returns RedisServiceOptions
func NewRedisServiceOptions(url string, prefix string, exp time.Duration, hashed bool) RedisServiceOptions {
	return RedisServiceOptions{
		Url: url, Prefix: prefix, Exp: exp, Hashed: hashed,
	}
}

// redisService provides cache functionality with redis
type redisService struct {
	// Client Client is a Redis client representing a pool of zero or more underlying connections. It's safe for concurrent use by multiple goroutines.
	Client *redis.Client
	// opts options
	opts RedisServiceOptions
}

// NewRedisService creates redisService
func NewRedisService(opts RedisServiceOptions) (rc redisService, err error) {
	rc.opts = opts
	options, err := redis.ParseURL(opts.Url)
	if err != nil {
		return
	}
	rc.Client = redis.NewClient(options)
	return
}

// Get key from cache
func (rc redisService) Get(key string) (value interface{}, err error) {
	return rc.Client.Get(rc.hashKey(key)).Result()
}

// SetWithExp Set key with expiration time
func (rc redisService) SetWithExp(key string, value interface{}, exp time.Duration) (err error) {
	_, err = rc.Client.Set(rc.hashKey(key), value, exp).Result()
	return
}

// SetWithoutExp Set key with no expiration time
func (rc redisService) SetWithoutExp(key string, value interface{}) (err error) {
	_, err = rc.Client.Set(rc.hashKey(key), value, 0).Result()
	return
}

// Set key
func (rc redisService) Set(key string, value interface{}) (err error) {
	_, err = rc.Client.Set(rc.hashKey(key), value, rc.opts.Exp).Result()
	return
}

// Unset key
func (rc redisService) Unset(key string) (err error) {
	_, err = rc.Client.Del(rc.hashKey(key)).Result()
	return
}

// HandleExpire handler for key expiration event
func (rc redisService) HandleExpire(handler func(string)) (err error) {
	// this is telling redis to publish events since it's off by default.
	_, err = rc.Client.Do("CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err != nil {
		fmt.Printf("unable to set keyspace events %v", err.Error())
		os.Exit(1)
	}
	// this is telling redis to subscribe to events published in the keyevent channel, specifically for expired events
	pubsub := rc.Client.PSubscribe("__keyevent@0__:expired")
	go func(redis.PubSub) {
		for { // infinite loop
			// this listens in the background for messages.
			message, err := pubsub.ReceiveMessage()
			if err != nil {
				fmt.Printf("error message - %v", err.Error())
				break
			}
			if message != nil {
				handler(message.Payload)
			}
		}
	}(*pubsub)
	return
}

// hashKey hashes key using md5
func (rc redisService) hashKey(key string) string {
	if !rc.opts.Hashed {
		return key
	}
	/* #nosec */
	hasher := md5.New()
	_, _ = hasher.Write([]byte(fmt.Sprintf("%s-%s", rc.opts.Prefix, key)))
	return hex.EncodeToString(hasher.Sum(nil))
}
