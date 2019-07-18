package pool

import (
	"errors"
	"sync"
)

type LocalPool struct {
	capacity int        // 容量
	running  int        // 运行中的数量
	free     int        // 空闲中的数量
	mux      sync.Mutex // 锁
}

func NewLocalPool(capacity int) *LocalPool {

	if capacity == 0 {
		capacity = MAX_CAPACITY
	}
	return &LocalPool{capacity: capacity, free: capacity}
}

func (p *LocalPool) Submit(task func()) error {

	p.mux.Lock()
	defer p.mux.Unlock()

	if p.running >= p.capacity {
		return errors.New("capacity limit")
	}

	p.running += 1
	p.free -= 1
	go func() {
		task()
		p.running -= 1
		p.free += 1
	}()

	return nil
}

func (p *LocalPool) Free() int {
	return p.free
}

func (p *LocalPool) Running() int {
	return p.running
}

func (p *LocalPool) Capacity() int {
	return p.capacity
}

func (p *LocalPool) NewCounter(count int) *Counter {
	if count == 0 {
		count = TRIGGER_COUNT
	}
	return &Counter{}
}
