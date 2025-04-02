package pkg

import (
	"errors"
	"time"
)

func Retry(operation func() error, maxRetries int, baseDelay time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		time.Sleep(baseDelay)

		baseDelay *= 2
	}

	return errors.New("operation failed after max retries")

}

func Timeout(operation func() error, timeout time.Duration) error {
	c := make(chan error, 1)

	go func() {
		c <- operation()
	}()

	select {
	case err := <-c:
		return err
	case <-time.After(timeout * time.Millisecond):
		return errors.New("Timeout")
	}
}

type DeadLetterQueue struct {
	mass []string
}

func (dlq *DeadLetterQueue) GetMessages() []string {
	return dlq.mass
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{
		mass: make([]string, 0),
	}
}

func ProcessWithDLQ(messages []string, operation func(string) error, dlq *DeadLetterQueue) {
	for _, msg := range messages {
		err := operation(msg)
		if err != nil {
			dlq.mass = append(dlq.mass, msg)
		}
	}
}
