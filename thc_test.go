package thc_test

import (
	"testing"

	"github.com/kyleraywed/thc"
	"github.com/kyleraywed/thc/thc_errs"
)

func TestContainerBasicOperations(t *testing.T) {
	hooksCalled := make(map[string]bool)

	c := thc.NewTHC(thc.FuncMap{
		"Store":  func() { hooksCalled["Store"] = true },
		"Fetch":  func() { hooksCalled["Fetch"] = true },
		"Update": func() { hooksCalled["Update"] = true },
		"Remove": func() { hooksCalled["Remove"] = true },
	})

	// Test storing int
	intKey, err := thc.Store(c, 42)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}
	if !hooksCalled["Store"] {
		t.Errorf("Store hook was not called")
	}

	// Test fetching int
	val, err := thc.Fetch(c, intKey)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}
	if val != 42 {
		t.Errorf("Expected 42, got %v", val)
	}
	if !hooksCalled["Fetch"] {
		t.Errorf("Fetch hook was not called")
	}

	// Test updating int
	err = thc.Update(c, intKey, 100)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if !hooksCalled["Update"] {
		t.Errorf("Update hook was not called")
	}

	// Fetch updated value
	val2, err := thc.Fetch(c, intKey)
	if err != nil {
		t.Fatalf("Fetch after update failed: %v", err)
	}
	if val2 != 100 {
		t.Errorf("Expected 100, got %v", val2)
	}

	// Test storing string
	strKey, err := thc.Store(c, "hello")
	if err != nil {
		t.Fatalf("Store string failed: %v", err)
	}

	// Fetch string
	strVal, err := thc.Fetch(c, strKey)
	if err != nil {
		t.Fatalf("Fetch string failed: %v", err)
	}
	if strVal != "hello" {
		t.Errorf("Expected 'hello', got %v", strVal)
	}

	// Test removing int
	err = thc.Remove(c, &intKey)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	if !hooksCalled["Remove"] {
		t.Errorf("Remove hook was not called")
	}

	// Fetch removed key should fail
	_, err = thc.Fetch(c, intKey)
	if err != thc_errs.ErrDeletedValue {
		t.Errorf("Expected ErrDeletedValue after remove, got %v", err)
	}

	// Test key identity mismatch
	secondContainer := thc.NewTHC(nil)
	secondKey, _ := thc.Store(secondContainer, 42)

	_, err = thc.Fetch(c, secondKey)
	if err != thc_errs.ErrIdentMismatch {
		t.Errorf("Expected ErrIdentMismatch, got %v", err)
	}
}
