package telnet

import (
	"time"
)

type (
	Option interface {
		isTelnetOption() // Attempt to stop someone from implementing this interface.
	}

	baseOption struct{}

	timeoutOption struct {
		baseOption
		timeout time.Duration
	}

	noTimeoutOption struct {
		baseOption
	}
)

func (_ baseOption) isTelnetOption() {}

var _ Option = new(baseOption)

func WithTimeout(d time.Duration) Option {
	return timeoutOption{
		timeout: d,
	}
}

func WithNoTimeout() Option {
	return noTimeoutOption{}
}

func handleOptions(options ...Option) (bool, time.Duration) {
	withTimeout := true
	timeout := time.Duration(60 * time.Second)

	for _, option := range options {
		switch op := option.(type) {
		case timeoutOption:
			timeout = op.timeout
			withTimeout = true
		case noTimeoutOption:
			withTimeout = false
		}
	}

	return withTimeout, timeout
}
