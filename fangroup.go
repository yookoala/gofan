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

	// initialize wait group
	fg.wg = &sync.WaitGroup{}

	fgm[fg] = ch
	return
}

func (fg *FanGroup) Lock() {
	// increment wait group
	fg.wg.Add(1)

	// pass or lock
	fg.ch() <- true
}

func (fg *FanGroup) Unlock() {
	go func() {
		// remove lock once
		<-fg.ch()

		// when lock removed, signal done
		fg.wg.Done()
	}()
}

func (fg *FanGroup) Run(j job) {
	go func() {
		fg.Lock()
		defer fg.Unlock()
		if err := j(); err != nil {
			panic(err)
		}
	}()
}

// Wait. Works like sync.WaitGroup.Wait()
// Wait until all members in the group finish
func (fg *FanGroup) Wait() {
	fg.wg.Wait()
}

// wait until finish to do something
func (fg *FanGroup) OnFinish(j job) {
	go func() {
		fg.Wait()
		if err := j(); err != nil {
			panic(err)
		}
	}()
}
