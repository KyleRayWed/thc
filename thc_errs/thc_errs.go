package thc_errs

import "errors"

var (
	ErrStoreSelf      = errors.New("a container may not store itself")
	ErrDeletedValue   = errors.New("value at key deleted")
	ErrConKeyMismatch = errors.New("container/key identity mismatch")
	ErrValNotFound    = errors.New("value not found")
	ErrTypeCast       = errors.New("type-casting error")
	ErrMissingValue   = errors.New("no value to remove at key")
)
