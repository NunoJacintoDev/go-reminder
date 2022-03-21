package reminder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Envelope(t *testing.T) {
	t.Parallel()

	r := NewReminder("redis://redis:6379")

	t.Run("at_should_define_id", func(t *testing.T) {
		testEnv := envelope{s: r}
		testEnv.At(time.Now().Add(time.Second))
		assert.NotEqual(t, "", testEnv.id)
	})

	t.Run("in_should_define_id", func(t *testing.T) {
		testEnv := envelope{s: r}
		testEnv.In(time.Second)
		assert.NotEqual(t, "", testEnv.id)
	})
}
