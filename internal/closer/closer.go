package closer

import "sync"

type closer []func() error

var c closer
var closerMux sync.Mutex

func Add(fn func() error) {
	closerMux.Lock()
	defer closerMux.Unlock()
	c = append(c, fn)
}

func New() {
	c = make(closer, 0, 10)
}

func Close() []error {
	closerMux.Lock()
	defer closerMux.Unlock()
	res := make([]error, 0, len(c))
	for i := len(c) - 1; i >= 0; i-- {
		i := i
		if err := c[i](); err != nil {
			res = append(res, err)
		}
	}
	return res
}
