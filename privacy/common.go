package privacy

import (
	"bytes"
	"chimney3-go/utils"
	"crypto/rand"
	"errors"
	"io"
)

func cloneBytes(src []byte) []byte {
	if len(src) == 0 {
		return nil
	}
	out := make([]byte, len(src))
	copy(out, src)
	return out
}

func randomBytes(size int) []byte {
	if size <= 0 {
		return nil
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return nil
	}
	return buf
}

// helper to encode method code + iv into bytes
func methodToBytes(code uint16, iv []byte) []byte {
	var op bytes.Buffer
	op.Write(utils.Uint162Bytes(code))
	lv := byte(len(iv))
	op.WriteByte(lv)
	if lv > 0 {
		op.Write(iv)
	}
	return op.Bytes()
}

// helper to decode iv from bytes (expects bytes after the 2-byte code)
func methodFromBytes(v []byte) ([]byte, error) {
	if v == nil {
		return nil, errors.New("nil input")
	}
	op := bytes.NewBuffer(v)
	lvl := op.Next(1)
	if len(lvl) < 1 {
		return nil, errors.New("out of length")
	}
	value := int(lvl[0])
	if value > 0 {
		iv := op.Next(value)
		if len(iv) != value {
			return nil, errors.New("iv length mismatch")
		}
		return iv, nil
	}
	return nil, nil
}
