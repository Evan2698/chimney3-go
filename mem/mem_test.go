package mem

import "testing"

func TestPoolGetPutRoundTrip(t *testing.T) {
	p := NewPool(16)
	buf := p.Get()
	if len(buf) != 16 {
		t.Fatalf("Get() len = %d, want 16", len(buf))
	}

	copy(buf, []byte("hello"))
	p.Put(buf)

	got := p.Get()
	if len(got) != 16 {
		t.Fatalf("Get() len after Put = %d, want 16", len(got))
	}
}

func TestPoolHandlesZeroSize(t *testing.T) {
	p := NewPool(0)
	if got := p.Get(); len(got) != 0 {
		t.Fatalf("Get() len = %d, want 0", len(got))
	}
	p.Put(nil)
}

func TestBufferHolderNilSafe(t *testing.T) {
	var b *bufferHolder
	if got := b.GetLarge(); got != nil {
		t.Fatalf("nil bufferHolder.GetLarge() = %v, want nil", got)
	}
	if got := b.GetSmall(); got != nil {
		t.Fatalf("nil bufferHolder.GetSmall() = %v, want nil", got)
	}
	b.PutLarge(nil)
	b.PutSmall(nil)
}
