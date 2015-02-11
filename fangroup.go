package gofan

import (
	"fmt"
	"sync"
)

var fgm map[*FanGroup]chan bool

func init() {
	fgm = make(map[*FanGroup]chan bool)
}

func NewGroup(size int) (fg *FanGroup) {
	fg = &FanGroup{
		size: size,
	}
	fg.ch() // initialize the channel anyway
	return
}

type FanGroup struct {
	size int
	wg   *sync.WaitGroup
}

func (fg *FanGroup) ch() (ch chan bool) {

	// try to accquire lock channel
	ch, ok := fgm[fg]
	if ok {
		return
	}

	// if lock channel doesn't exists, create one
	if fg.size < 0 {
		panic(fmt.Errorf("Pool size must be unsigned integer. "+
			"Cannot be %d", fg.size))
	}
	ch = make(chan bool, fg.size)
	fg.wg = &sync.WaitGroup{}

	fgm[fg] = ch
	return
}

func (fg *FanGroup) Lock() {
	// pass or lock
	fg.ch() <- true
	fg.wg.Add(1)
}

func (fg *FanGroup) Unlock() {
	// remove lock once
	<-fg.ch()
	fg.wg.Done()
}
