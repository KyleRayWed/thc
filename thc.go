package thc

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

var removedID = uuid.NewString() // so little truly matters

type dataMap map[string]struct {
	value        any
	timeModified time.Time
}

type container struct {
	identity string
	data     dataMap
	mut      sync.RWMutex // goroutine safety compliance
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
func NewTHC() container {
	return container{
		identity: uuid.NewString(),
		data:     make(dataMap),
	}
}

// Create a key that has an instantiated type of the input itself.
func Store[T any](c *container, input T) (Key[T], error) {
	switch any(input).(type) {
	case container:
		if any(input).(container).identity == c.identity {
			var zero Key[T]
			return zero, fmt.Errorf("container may not store itself")
		}
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
		return zero, fmt.Errorf("deleted value at key")
	}
	if c.identity != key.identity {
		return zero, fmt.Errorf("container/key identity mismatch")
	}

	c.mut.RLock()
	defer c.mut.RUnlock()

	val, ok := c.data[key.mapKey]
	if !ok {
		return zero, fmt.Errorf("value not found")
	}

	casted, ok := val.value.(T)
	if !ok {
		return zero, fmt.Errorf("type-casting error")
	}
	return casted, nil
}

// Update the value at the key within the container. Types must match.
func Update[T any](c *container, key Key[T], input T) error {
	switch any(input).(type) {
	case container:
		if any(input).(container).identity == c.identity {
			return fmt.Errorf("container may not store itself")
		}
	}
	if key.identity == removedID {
		return fmt.Errorf("deleted value at key")
	}
	if c.identity != key.identity {
		return fmt.Errorf("container/key identity mismatch")
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
		return fmt.Errorf("deleted value at key")
	}
	if c.identity != key.identity {
		return fmt.Errorf("container/key identity mismatch")
	}

	c.mut.Lock()
	defer c.mut.Unlock()

	_, ok := c.data[key.mapKey]
	if !ok {
		return fmt.Errorf("no value to remove at key")
	}

	key.identity = removedID
	delete(c.data, key.mapKey)

	return nil
}
