package poollock

import (
	"sync"
	"testing"
	"time"
)

// Test if PoolLock implements sync.Locker
func TestPoolLockSync(t *testing.T) {
	var l sync.Locker = New(0)
	_ = l
	t.Log("PoolLock implements sync.Locker")
}

// Test of the lock can block
func TestPoolLock(t *testing.T) {

	size := 4

	l := New(size)

	// lock number = the pool size
	for i := 0; i < size; i++ {
		l.Lock()
	}

	// try acquire new lock until timeout
	// should all be blocked
	notBlocked := make(chan bool)
	go func() {
		l.Lock()
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

	l := New(size)

	// lock number = the pool size
	for i := 0; i < size; i++ {
		l.Lock()
	}

	// unlock once
	l.Unlock()

	// try acquire new lock until timeout
	// should not be blocked
	notBlocked := make(chan bool)
	go func() {
		l.Lock()
		notBlocked <- true
	}()

	select {
	case <-notBlocked:
		t.Log("Unlock success")
	case <-time.After(time.Second):
		t.Error("Unlock failed")
	}

}
