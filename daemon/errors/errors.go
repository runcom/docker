package errors

import (
	"fmt"
	"time"
)

type ErrWaitTimeout struct {
	Timeout time.Duration
}

func (err ErrWaitTimeout) Error() string {
	return fmt.Sprintf("Timed out: %v", err.Timeout)
}

func NewErrWaitTimeout(timeout time.Duration) ErrWaitTimeout {
	return ErrWaitTimeout{Timeout: timeout}
}

type ErrRemovalInProgress struct{}

func (err ErrRemovalInProgress) Error() string {
	return "Status is already RemovalInProgress"
}
