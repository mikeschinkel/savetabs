package ui

import (
	"errors"
)

var (
	ErrInvalidKeyFormat = errors.New("invalid key format (expected '<type>-<key>')")
)
