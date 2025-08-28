package model

import "errors"

var (
	// Session errors
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionBadRequest  = errors.New("bad request")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// User errors
	ErrUserNotFound = errors.New("user not found")
)
