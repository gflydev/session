package session

import "github.com/gflydev/core/errors"

var (
	ErrNotSetProvider = errors.New("session provider is not set")
	ErrEmptySessionID = errors.New("session id is empty")
)
