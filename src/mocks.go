package reminder

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockedService struct {
	mock.Mock
}

// RemindAt implements reminder.Service interface for mocking
func (m *MockedService) RemindAt(message interface{}, date time.Time) error {
	args := m.Called(message, date)
	return args.Error(1)
}

// RemindIn implements reminder.Service interface for mocking
func (m *MockedService) RemindIn(message interface{}, duration time.Duration) error {
	args := m.Called(message, duration)
	return args.Error(1)
}
