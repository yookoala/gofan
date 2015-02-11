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

// generate a stream of integer for testing
func generate(a int) (out chan int) {
	out = make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= a; i++ {
			out <- i
		}
	}()

	return
}

// a simple function to check test result
func addTo(n int) int {
	if n == 0 {
		return 0
	}
	return n + addTo(n-1)
}

// a simple fan out / fan in funciton for test
func pipe(b int, in chan int) (out chan int) {
	out = make(chan int)
	fg := NewGroup(b)

	for i := range in {

		// clone variable for fan out
		j := i

		// store the function and run until ready
		fg.Run(func() {
			out <- j // just pass the number out
		})
	}

	// start a goroutine to wait
	fg.OnFinish(func() {
		close(out)
	})

	return
}

// test the fan out / fan in pattern against
// given pipe scale and input scale
func testPipe(t *testing.T, pscale, tscale int) {

	s := 0
	c := 0
	for n := range pipe(pscale, generate(tscale)) {
		c++
		s += n
	}

	// check count of numbers
	if c != tscale {
		t.Errorf("Error: `generate()` is expected to yield "+
			"%d numbers. Get %d", tscale, c)
	} else {
		t.Logf("`generate()` has yield %d numbers", tscale)
	}

	// test sum of numbers
	es := addTo(tscale) // expected sum
	if s != es {
		t.Errorf("Error: Sum of generated numbers "+
			"is expected to be %d. Get %d", es, s)
	} else {
		t.Logf("Sum of generated numbers is %d", es)
	}

}

// test the generated series
func TestPipe(t *testing.T) {

	// test when pipe scale > input scale
	testPipe(t, 20, 5)

	// test when pipe scale == input scale
	testPipe(t, 10, 10)

	// test when pipe scale < input scale
	testPipe(t, 5, 20)

}
