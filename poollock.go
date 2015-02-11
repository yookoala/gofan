package poollock

var lm map[*PoolLock]chan bool

func init() {
	lm = make(map[*PoolLock]chan bool)
}

func New(size int) *PoolLock {
	return &PoolLock{
		size: size,
	}
}

type PoolLock struct {
	size int
}

func (l *PoolLock) ch() (ch chan bool) {
	// accquire lock channel, or create one
	ch, ok := lm[l]
	if !ok {
		ch = make(chan bool, l.size)
		lm[l] = ch
	}
	return ch
}

func (l *PoolLock) Lock() {
	// pass or lock
	l.ch() <- true
}

func (l *PoolLock) Unlock() {
	// remove lock once
	<-l.ch()
}
