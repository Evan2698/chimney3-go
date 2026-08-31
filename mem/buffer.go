package mem

import (
	"sync"
)

type bufferHolder struct {
	smallbuffer *Pool
	largebuffer *Pool
}

type Buffer interface {
	GetLarge() []byte
	PutLarge([]byte)
	GetSmall() []byte
	PutSmall([]byte)
}

var (
	instance *bufferHolder
	once     sync.Once
)

const (
	LARGE_BUFFER_SIZE = 4096
	SMALL_BUFFER_SIZE = 512
)

func newApplicationBuffer() Buffer {
	once.Do(func() {
		instance = &bufferHolder{
			smallbuffer: NewPool(SMALL_BUFFER_SIZE),
			largebuffer: NewPool(LARGE_BUFFER_SIZE),
		}
	})
	return instance
}

// Convenience package-level helpers so callers don't need to call NewApplicationBuffer().
func GetLarge() []byte  { return newApplicationBuffer().GetLarge() }
func PutLarge(b []byte) { newApplicationBuffer().PutLarge(b) }
func GetSmall() []byte  { return newApplicationBuffer().GetSmall() }
func PutSmall(b []byte) { newApplicationBuffer().PutSmall(b) }

func (b *bufferHolder) GetLarge() []byte {
	if b == nil || b.largebuffer == nil {
		return nil
	}
	return b.largebuffer.Get()
}

func (b *bufferHolder) PutLarge(t []byte) {
	if b == nil || b.largebuffer == nil {
		return
	}
	b.largebuffer.Put(t)
}

func (b *bufferHolder) GetSmall() []byte {
	if b == nil || b.smallbuffer == nil {
		return nil
	}
	return b.smallbuffer.Get()
}

func (b *bufferHolder) PutSmall(t []byte) {
	if b == nil || b.smallbuffer == nil {
		return
	}
	b.smallbuffer.Put(t)
}
