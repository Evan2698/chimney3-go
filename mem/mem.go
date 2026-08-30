package mem

import (
	"sync"
)

type Pool struct {
	pool *sync.Pool
	size int
}

func normalizePoolSize(size int) int {
	if size < 0 {
		return 0
	}
	return size
}

func normalizePoolBuffer(buf []byte, size int) []byte {
	if size <= 0 {
		return nil
	}
	if cap(buf) >= size {
		return buf[:size]
	}
	nb := make([]byte, size)
	copy(nb, buf)
	return nb
}

func NewPool(size int) *Pool {
	size = normalizePoolSize(size)
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
	if p == nil {
		return nil
	}
	if p.size <= 0 {
		return make([]byte, 0)
	}
	b := p.pool.Get().([]byte)
	return normalizePoolBuffer(b, p.size)
}

func (p *Pool) Put(b []byte) {
	if p == nil {
		return
	}
	if p.size <= 0 {
		return
	}
	p.pool.Put(normalizePoolBuffer(b, p.size))
}
