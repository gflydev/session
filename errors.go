package session

import "github.com/gflydev/core/errors"

var (
	ErrNotSetProvider = errors.New("not set a session provider")
	ErrEmptySessionID = errors.New("empty session id")
)
