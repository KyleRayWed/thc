package main

import (
	"testing"
)

func TestStoreFetchUpdateRemove(t *testing.T) {
	c := NewTHC(nil)

	// Store int
	intKey, err := Store(&c, 42)
	if err != nil {
		t.Fatalf("unexpected error storing int: %v", err)
	}

	// Store string
	strKey, err := Store(&c, "hello")
	if err != nil {
		t.Fatalf("unexpected error storing string: %v", err)
	}

	// Fetch int
	valInt, err := Fetch(&c, intKey)
	if err != nil {
		t.Fatalf("unexpected error fetching int: %v", err)
	}
	if valInt != 42 {
		t.Errorf("expected 42, got %v", valInt)
	}

	// Fetch string
	valStr, err := Fetch(&c, strKey)
	if err != nil {
		t.Fatalf("unexpected error fetching string: %v", err)
	}
	if valStr != "hello" {
		t.Errorf("expected 'hello', got %v", valStr)
	}

	// Update int
	if err := Update(&c, intKey, 99); err != nil {
		t.Fatalf("unexpected error updating int: %v", err)
	}
	valInt, _ = Fetch(&c, intKey)
	if valInt != 99 {
		t.Errorf("expected 99 after update, got %v", valInt)
	}

	// Remove string
	if err := Remove(&c, &strKey); err != nil {
		t.Fatalf("unexpected error removing string: %v", err)
	}

	// Fetch after remove should error
	_, err = Fetch(&c, strKey)
	if err == nil {
		t.Errorf("expected error fetching removed key, got nil")
	}

	// Update after remove should error
	if err := Update(&c, strKey, "world"); err == nil {
		t.Errorf("expected error updating removed key, got nil")
	}

	// Remove again should error
	if err := Remove(&c, &strKey); err == nil {
		t.Errorf("expected error removing twice, got nil")
	}
}

func TestIdentityMismatch(t *testing.T) {
	c1 := NewTHC(nil)
	c2 := NewTHC(nil)

	key, err := Store(&c1, "data")
	if err != nil {
		t.Fatalf("unexpected error storing: %v", err)
	}

	// Try to fetch with wrong container
	_, err = Fetch(&c2, key)
	if err == nil {
		t.Errorf("expected identity mismatch error, got nil")
	}

	// Try to update with wrong container
	if err := Update(&c2, key, "new"); err == nil {
		t.Errorf("expected identity mismatch error, got nil")
	}

	// Try to remove with wrong container
	if err := Remove(&c2, &key); err == nil {
		t.Errorf("expected identity mismatch error, got nil")
	}
}

func TestStringAndLen(t *testing.T) {
	c := NewTHC(nil)

	if c.Len() != 0 {
		t.Errorf("expected len=0, got %d", c.Len())
	}

	_, _ = Store(&c, 123)
	_, _ = Store(&c, "abc")

	if c.Len() != 2 {
		t.Errorf("expected len=2, got %d", c.Len())
	}

	str := c.String()
	if str != "Length: 2" {
		t.Errorf("unexpected String(): %s", str)
	}
}
