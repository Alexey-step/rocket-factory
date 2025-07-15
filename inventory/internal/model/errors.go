package model

import "errors"

var (
	ErrPartNotFound       = errors.New("part not found")
	ErrPartsNotFound      = errors.New("parts not found")
	ErrPartsInternalError = errors.New("internal error while getting parts")
)
