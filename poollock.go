package poollock

import (
	"fmt"
)

var lm map[*PoolLock]chan bool

func init() {
	lm = make(map[*PoolLock]chan bool)
}

func New(size int) (l *PoolLock) {
	l = &PoolLock{
		size: size,
	}
	l.ch() // initialize the channel anyway
	return
}

type PoolLock struct {
	size int
}

func (l *PoolLock) ch() (ch chan bool) {

	// try to accquire lock channel
	ch, ok := lm[l]
	if ok {
		return
	}

	// if lock channel doesn't exists, create one
	if l.size < 0 {
		panic(fmt.Errorf("Pool size must be unsigned integer. "+
			"Cannot be %d", l.size))
	}
	ch = make(chan bool, l.size)
	lm[l] = ch
	return
}

func (l *PoolLock) Lock() {
	// pass or lock
	l.ch() <- true
}

func (l *PoolLock) Unlock() {
	// remove lock once
	<-l.ch()
}
