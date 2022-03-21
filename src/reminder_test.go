package reminder

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Reminder(t *testing.T) {

	// Setup Service
	r := NewReminder("redis://redis:6379")

	// Listen for notifications
	r.Listen(func(n Notification, err error) {
		if err != nil {
			fmt.Println("🔥Error🔥", err)
		} else {
			fmt.Println("✨Event✨", n)
		}
	})

	t.Run("remind_me_in", func(t *testing.T) {
		t.Parallel()
		err := r.Remind("test_remind_me_in").In(time.Second * 1)
		assert.NoError(t, err)
		time.Sleep(time.Second * 2)
	})

	t.Run("remind_me_at", func(t *testing.T) {
		t.Parallel()
		err := r.Remind("test_remind_me_at").At(time.Now().Add(time.Second * 1))
		assert.NoError(t, err)
		time.Sleep(time.Second * 2)
	})

}
