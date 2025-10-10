package thc

import (
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/kyleraywed/thc/thc_errs"
)

type FuncMap map[string]func()

type container struct {
	identity  string
	removedID string
	data      sync.Map // concurrent safe map

	auditHook FuncMap
}

type Key[T any] struct {
	identity string
	mapKey   string
}

// Number of records in the underlying map
func (c *container) Len() int {
	count := 0
	c.data.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

// String representation
func (c *container) String() string {
	return "Length: " + strconv.Itoa(c.Len())
}

// Container constructor. Handler's keys are strings and correpsond with
// the 4 transactions. Don't forget to capitalize. Func is run only on
// sucessful transaction.
func NewTHC(handler FuncMap) *container {
	return &container{
		identity:  uuid.NewString(),
		removedID: uuid.NewString(),
		auditHook: handler,
	}
}

// Store a value, get a key
func Store[T any](c *container, input T) (Key[T], error) {
	switch any(input).(type) {
	case container:
		if any(input).(container).identity == c.identity {
			var zero Key[T]
			return zero, thc_errs.ErrStoreSelf
		}
	}

	newKey := uuid.NewString()
	c.data.Store(newKey, input)

	if fn, ok := c.auditHook["Store"]; ok {
		fn()
	}

	return Key[T]{identity: c.identity, mapKey: newKey}, nil
}

// Fetch a value with key, get type-casted value
func Fetch[T any](c *container, key Key[T]) (T, error) {
	var zero T

	if key.identity == c.removedID {
		return zero, thc_errs.ErrDeletedValue
	}
	if c.identity != key.identity {
		return zero, thc_errs.ErrIdentMismatch
	}

	val, ok := c.data.Load(key.mapKey)
	if !ok {
		return zero, thc_errs.ErrValNotFound
	}

	casted, ok := val.(T)
	if !ok {
		return zero, thc_errs.ErrTypeCast
	}

	if fn, ok := c.auditHook["Fetch"]; ok {
		fn()
	}

	return casted, nil
}

// Update a value (must be same type)
func Update[T any](c *container, key Key[T], input T) error {
	switch any(input).(type) {
	case container:
		if any(input).(container).identity == c.identity {
			return thc_errs.ErrStoreSelf
		}
	}

	if key.identity == c.removedID {
		return thc_errs.ErrDeletedValue
	}
	if c.identity != key.identity {
		return thc_errs.ErrIdentMismatch
	}

	c.data.Store(key.mapKey, input)

	if fn, ok := c.auditHook["Update"]; ok {
		fn()
	}

	return nil
}

// Remove a value, invalidate key
func Remove[T any](c *container, key *Key[T]) error {
	if key.identity == c.removedID {
		return thc_errs.ErrDeletedValue
	}
	if c.identity != key.identity {
		return thc_errs.ErrIdentMismatch
	}

	_, ok := c.data.Load(key.mapKey)
	if !ok {
		return thc_errs.ErrMissingValue
	}

	key.identity = c.removedID
	c.data.Delete(key.mapKey)

	if fn, ok := c.auditHook["Remove"]; ok {
		fn()
	}

	return nil
}
