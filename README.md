# THC

A small (**t**)ype-safe, (**h**)eterogeneous (**c**)ontainer. It allows you to store values, retrieve those values with typed keys, and delete stored values safely.

```go
func NewTHC(handler func()) container

func Store[T any](container *thc_container, input T) (thc_key[T], error)
func Fetch[T any](container *thc_container, key thc_key[T]) (T, error)
func Update[T any](container *thc_container, key thc_key[T], input T) error
func Remove[T](container *thc_container, key *thc_key[T]) error
```

Usage

```go
package main

import (
    "fmt"
    "github.com/kyleraywed/thc"
)

func main() {
    // Create a new container. Upon every successful transaction, it will log its success.
    c := thc.NewTHC(thc.FuncMap{
		"Store": func() {
			log.Println("Sucessful store.")
		},
		"Fetch": func() {
			log.Println("Sucessful fetch.")
		},
        "Update": func() {
            log.Println("Sucessful update.")
        },
        "Remove": func() {
            log.Println("Sucessful removal.")
        },
	})

    // Store a string (or anything) in the container, get a key
    k, _ := thc.Store(&c, "hello, world")

    // Use the key to Fetch the value back
    v, _ := thc.Fetch(&c, k)
    fmt.Println("value:", v)

    // Update value (must be same type)
    thc.Update(&c, k, "goodbye, world")

    // Delete value and invalidate key
    if err := thc.Remove(&c, &k); err != nil {
        panic(err)
    }
}
```

Notes and design

- Attempting to use a key on a container it isn't associated with will result in error.
- Attempting to store a container within itself will result in error.
- Goroutine safe.
