package model

import "errors"

var (
	ErrPartNotFound       = errors.New("part not found")
	ErrPartsInternalError = errors.New("internal error while getting parts")
)
