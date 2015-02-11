package gofan

import (
	"sync"
	"testing"
	"time"
)

// Test if FanGroup implements sync.Locker
func TestFanGroupSync(t *testing.T) {
	var l sync.Locker = NewGroup(0)
	_ = l
	t.Log("FanGroup implements sync.Locker")
}

// Test error for size < 0
func TestFanGroupErr(t *testing.T) {

	// try to recover from panic
	defer func() {
		if r := recover(); r != nil {
			t.Log("FanGroup size < 0 triggers error")
		}
	}()
	NewGroup(-1)

	// should not have run till here
	t.Errorf("FanGroup size < 0 doesn't trigger error")
}

// Test of the lock can block
// Also waits until explicitly
func TestFanGroup(t *testing.T) {

	size := 4

	fg := NewGroup(size)

	// lock number = the pool size
	for i := 0; i < size; i++ {
		fg.Lock()
	}

	// try acquire new lock until timeout
	// should all be blocked
	notBlocked := make(chan bool)
	go func() {
		fg.Lock()
		notBlocked <- true
	}()

	select {
	case <-notBlocked:
		t.Error("Lock failed")
	case <-time.After(time.Second):
		t.Log("Lock success")
	}

}

// Test if the lock can unlock
func TestPoolUnlock(t *testing.T) {

	size := 4

	fg := NewGroup(size)

	// lock number = the pool size
	for i := 0; i < size; i++ {
		fg.Lock()
	}

	// unlock once
	fg.Unlock()

	// try acquire new lock until timeout
	// should not be blocked
	notBlocked := make(chan bool)
	go func() {
		fg.Lock()
		notBlocked <- true
	}()

	select {
	case <-notBlocked:
		t.Log("Unlock success")
	case <-time.After(time.Second):
		t.Error("Unlock failed")
	}

}
