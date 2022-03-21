package reminder

import "time"

type Service interface {
	RemindAt(message interface{}, at time.Time) error
	RemindIn(message interface{}, in time.Duration) error
}
