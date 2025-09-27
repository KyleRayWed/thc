/*
	TODO
		clear/reset/delete-all
		helper functions like helper.RemoveDups() to go into the FuncMap
*/

package thc

import (
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kyleraywed/thc/thc_errs"
)

var removedID = uuid.NewString() // so little truly matters

type FuncMap map[string]func()

type dataMap map[string]struct {
	value        any
	timeModified time.Time
}

type container struct {
	identity string
	data     dataMap
	mut      sync.RWMutex // goroutine safety compliance

	maintainMap FuncMap
	//maintainWait time.Duration
}

type Key[T any] struct {
	identity string
	mapKey   string
}

// Number of records in the underlying map as a string
func (c *container) String() string {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return "Length: " + strconv.Itoa(len(c.data))
}

// Number of records in the underlying map
func (c *container) Len() int {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return len(c.data)
}

// Initialize container with a unique identity and fresh dataMap
// as well as a handler function that runs after a successful transaction
func NewTHC(handler FuncMap) container {
	return container{
		identity: uuid.NewString(),
		data:     make(dataMap),

		maintainMap: handler,
		//maintainWait: wait,
	}
}

// Create a key that has an instantiated type of the input itself.
func Store[T any](c *container, input T) (Key[T], error) {
	switch any(input).(type) {
	case container:
		if any(input).(container).identity == c.identity {
			var zero Key[T]
			return zero, thc_errs.ErrStoreSelf
		}
	}

	// only run if you make it past the error checks
	if fn, ok := c.maintainMap["Store"]; ok {
		defer fn()
	}

	newKey := uuid.NewString()

	c.mut.Lock()
	defer c.mut.Unlock()

	c.data[newKey] = struct {
		value        any
		timeModified time.Time
	}{
		value:        input,
		timeModified: time.Now(),
	}

	return Key[T]{
		identity: c.identity,
		mapKey:   newKey,
	}, nil
}

// Return the value at the key within the container, typecasted the type of the key.
func Fetch[T any](c *container, key Key[T]) (T, error) {
	var zero T

	if key.identity == removedID {
		return zero, thc_errs.ErrDeletedValue
	}
	if c.identity != key.identity {
		return zero, thc_errs.ErrConKeyMismatch
	}

	c.mut.RLock()
	defer c.mut.RUnlock()

	val, ok := c.data[key.mapKey]
	if !ok {
		return zero, thc_errs.ErrValNotFound
	}

	casted, ok := val.value.(T)
	if !ok {
		return zero, thc_errs.ErrTypeCast
	}

	// only run if you make it past the error checks
	if fn, ok := c.maintainMap["Fetch"]; ok {
		defer fn()
	}

	return casted, nil
}

// Update the value at the key within the container. Types must match.
func Update[T any](c *container, key Key[T], input T) error {
	switch any(input).(type) {
	case container:
		if any(input).(container).identity == c.identity {
			return thc_errs.ErrStoreSelf
		}
	}
	if key.identity == removedID {
		return thc_errs.ErrDeletedValue
	}
	if c.identity != key.identity {
		return thc_errs.ErrConKeyMismatch
	}

	// only run if you make it past the error checks
	if fn, ok := c.maintainMap["Update"]; ok {
		defer fn()
	}

	c.mut.Lock()
	defer c.mut.Unlock()

	c.data[key.mapKey] = struct {
		value        any
		timeModified time.Time
	}{
		value:        input,
		timeModified: time.Now(),
	}
	return nil
}

// Delete the value at the key within the container and mark as removed.
func Remove[T any](c *container, key *Key[T]) error {
	if key.identity == removedID {
		return thc_errs.ErrDeletedValue
	}
	if c.identity != key.identity {
		return thc_errs.ErrConKeyMismatch
	}

	c.mut.Lock()
	defer c.mut.Unlock()

	_, ok := c.data[key.mapKey]
	if !ok {
		return thc_errs.ErrMissingValue
	}

	// only run if you make it past the error checks
	if fn, ok := c.maintainMap["Remove"]; ok {
		defer fn()
	}

	key.identity = removedID
	delete(c.data, key.mapKey)

	return nil
}
