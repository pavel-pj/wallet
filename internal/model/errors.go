package model

import "errors"

var (
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrInvalidUUID          = errors.New("invalid wallet uuid")
	ErrInvalidAmount        = errors.New("amount must be positive")
	ErrInvalidOperationType = errors.New("invalid operation type")
)
