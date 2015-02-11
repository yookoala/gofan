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

func (l *PoolLock) Lock() {

	// accquire lock channel, or create one
	ch, ok := lm[l]
	if !ok {
		ch = make(chan bool, l.size+4)
		lm[l] = ch
	}

	// pass or lock
	ch <- true
}
