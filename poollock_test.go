package poollock

import (
	"testing"
	"time"
)

func TestPoollock(t *testing.T) {

	size := 4

	l := &PoolLock{
		size: size,
	}

	// lock number = the pool size
	for i := 0; i < size; i++ {
		l.Lock()
	}

	// try acquire new lock until timeout
	// should all be blocked
	strange := make(chan bool)
	go func() {
		l.Lock()
		strange <- true
	}()

	select {
	case <-strange:
	case <-time.After(time.Second):
		t.Error("Lock failed")
	}

	t.Log("Lock success")
}
