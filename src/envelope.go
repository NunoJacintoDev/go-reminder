package reminder

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/google/uuid"
)

type envelope struct {
	s          *reminder
	message    interface{}
	id         string
	createdAt  time.Time
	expiration time.Duration
}

// At sets the envelope notification to "date"
func (r *envelope) At(date time.Time) (err error) {
	r.generateID()
	return r.In(time.Until(date))
}

// In sets the envelope notification from now to "duration"
func (r *envelope) In(duration time.Duration) (err error) {
	r.generateID()

	r.expiration = duration

	// Add shadow key with expiration
	shadowKey := getShadowKey(r.id)
	err = r.s.rcs.SetWithExp(shadowKey, nil, duration)
	if err != nil {
		return
	}

	// convert note into data to store
	data, err := r.encode()
	if err != nil {
		return
	}

	// Add original with event data
	err = r.s.rcs.SetWithoutExp(r.id, data)
	if err != nil {
		return
	}
	return
}

func (r *envelope) encode() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	toCache := &Notification{
		ID:         r.id,
		Message:    r.message,
		CreatedAt:  r.createdAt,
		Expiration: r.expiration,
	}
	err = enc.Encode(toCache)
	if err != nil {
		return
	}
	return buf.Bytes(), nil
}

func decode(data []byte) (n Notification, err error) {
	var buf bytes.Buffer
	_, err = buf.Write(data)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(&buf)
	var decoded Notification
	err = dec.Decode(&decoded)
	if err != nil {
		return
	}
	return decoded, nil
}

func (r *envelope) generateID() {
	r.id = uuid.NewString()
}
