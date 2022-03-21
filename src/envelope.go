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
	ID         string
	Reference  string
	CreatedAt  time.Time
	Expiration time.Duration
}

// WithReference adds a reference to your remind note
func (r *envelope) WithReference(reference string) *envelope {
	r.Reference = reference
	return r
}

// At sets the envelope notification to "date"
func (r *envelope) At(date time.Time) (err error) {
	r.generateID()
	return r.In(time.Until(date))
}

// In sets the envelope notification from now to "duration"
func (r *envelope) In(duration time.Duration) (err error) {
	r.generateID()

	r.Expiration = duration

	// Add shadow key with expiration
	shadowKey := getShadowKey(r.ID)
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
	err = r.s.rcs.SetWithoutExp(r.ID, data)
	if err != nil {
		return
	}
	return
}

func (r *envelope) encode() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	toCache := &Notification{
		ID:         r.ID,
		Message:    r.message,
		Reference:  r.Reference,
		CreatedAt:  r.CreatedAt,
		Expiration: r.Expiration,
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
	r.ID = uuid.NewString()
}
