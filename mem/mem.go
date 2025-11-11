package mem

import (
	"sync"
)

type Pool struct {
	pool *sync.Pool
	size int
}

func NewPool(size int) *Pool {
	return &Pool{
		size: size,
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		},
	}
}

func (p *Pool) Get() []byte {
	b := p.pool.Get().([]byte)
	// ensure returned slice has the expected length
	if cap(b) >= p.size {
		return b[:p.size]
	}
	// fallback: allocate fresh slice of required size
	nb := make([]byte, p.size)
	copy(nb, b)
	return nb
}

func (p *Pool) Put(b []byte) {
	// normalize length to capacity for reuse consistency
	if cap(b) >= p.size {
		p.pool.Put(b[:p.size])
		return
	}
	// if slice is smaller than pool size, allocate a properly sized one
	nb := make([]byte, p.size)
	copy(nb, b)
	p.pool.Put(nb)
}
