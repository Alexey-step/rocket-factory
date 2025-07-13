package model

import "errors"

// Order errors
var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderInternalError    = errors.New("internal error while get order")
	ErrOrderConflict         = errors.New("order conflict")
	ErrOrderAlreadyPaid      = errors.New("order already paid, cannot be cancelled")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled, cannot be cancelled again")
)

// Parts errors
var (
	ErrPartsNotFound = errors.New("some parts not found")
)

// Payment errors
var (
	ErrPaymentInternalError = errors.New("internal error while processing payment")
	ErrPaymentConflict      = errors.New("payment conflict")
	ErrPaymentNotFound      = errors.New("payment not found")
)
