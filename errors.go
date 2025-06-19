package session

import "github.com/gflydev/core/errors"

var (
	ErrNotSetProvider = errors.New("Not set a session provider")
	ErrEmptySessionID = errors.New("Empty session id")
)
