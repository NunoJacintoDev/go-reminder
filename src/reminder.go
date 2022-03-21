package reminder

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	shadowPrefix = "shadow:"
)

type notifyFn = func(n Notification, err error)

// Notification reminder notification event
type Notification struct {
	ID, Reference string
	CreatedAt     time.Time
	Expiration    time.Duration
	Message       interface{}
}

type reminder struct {
	notify notifyFn
	rcs    *redisService
}

// NewReminder creates a reminder that can be used to create reminder envelopes and notify you later with them
func NewReminder(redisURL string) (s *reminder) {
	s = &reminder{}
	rcs, err := NewRedisService(RedisServiceOptions{
		Url:    redisURL,
		Hashed: false, // must be NOT be hashed! we use the keys to map events
	})
	if err != nil {
		log.Fatal("Error connecting to redis")
		return
	}
	s.rcs = &rcs
	return s
}

// Listen for reminder notifications
func (s *reminder) Listen(fn notifyFn) (err error) {
	s.notify = fn
	if s.rcs == nil {
		err = fmt.Errorf("redis cache service not initialized")
		return
	}
	s.rcs.HandleExpire(func(key string) { s.handlerShadowKeyExpiration(key) })
	return
}

// Remind creates a envelope with the incoming "message"
func (s *reminder) Remind(message interface{}) *envelope {
	return &envelope{
		s:         s,
		message:   message,
		CreatedAt: time.Now(),
	}
}

// ---------------------------
// ----private functions -----
// ---------------------------

// handlerShadowKeyExpiration handles shadow key expiration event
func (s *reminder) handlerShadowKeyExpiration(shadowKey string) {
	// Get Key from Shadow Key
	key, isShadowKey := getKey(shadowKey)
	if !isShadowKey {
		return // do nothing
	}
	// Delete key
	defer s.rcs.Unset(key)

	// Get event from Original Key
	data, err := s.rcs.Get(key)
	if err != nil {
		s.notify(Notification{}, err)
		return
	}
	// convert data to note
	note := []byte(data.(string))
	n, err := decode(note)
	if err != nil {
		s.notify(n, err)
		return
	}
	s.notify(n, err)

}

// getKey get original key from shadow key, and check if incoming key is a shadow key (ok)
func getKey(shadowKey string) (key string, ok bool) {
	ok = strings.HasPrefix(shadowKey, shadowPrefix)
	key = strings.TrimPrefix(shadowKey, shadowPrefix)
	return
}

// getKey get shadow key from original key
func getShadowKey(key string) (shadowKey string) {
	return shadowPrefix + key
}
