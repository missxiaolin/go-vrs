package pool

import (
	"github.com/panjf2000/ants"
	"sync"
)

type Pool struct {
	p   *ants.Pool
	mux sync.Mutex
}

func NewPool(capacity int) (*Pool, error) {

	if capacity == 0 {
		capacity = MAX_CAPACITY
	}
	p, err := ants.NewPool(capacity)
	if err != nil {
		return nil, err
	}
	return &Pool{p: p}, nil
}

func (p *Pool) Submit(task func()) error {
	return p.p.Submit(task)
}

func (p *Pool) Free() int {
	return p.p.Free()
}

func (p *Pool) TuneSize(size int) {
	p.p.Tune(size)
}

func (p *Pool) NewCounter(count int) *Counter {
	if count == 0 {
		count = TRIGGER_COUNT
	}
	return &Counter{pool: p, triggerCount: count}
}