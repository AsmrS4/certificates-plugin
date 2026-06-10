package models

import "errors"

var (
	ErrOrderNotPending   = errors.New("Order is not in pending status")
	ErrOrderNotInPrepare = errors.New("Order is not in prepare status")
	ErrAlreadyRejected   = errors.New("Order already rejected")
	ErrOrderNotFound     = errors.New("Order not found")
)
