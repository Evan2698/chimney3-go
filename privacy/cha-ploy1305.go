package privacy

import (
	"chimney3-go/utils"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

type ploy struct {
	iv []byte
}

const (
	ployName = "CHACHA-POLY1305"
	ployCode = 0x1236
)

func (p *ploy) Compress(src []byte, key []byte, out []byte) (int, error) {
	defer utils.Trace("Compress")()
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return 0, err
	}

	ciphertext := aead.Seal(nil, p.iv, src, nil)
	if len(ciphertext) == 0 {
		return 0, errors.New("compressed failed")
	}

	n := len(ciphertext)
	if len(out) < n {
		return 0, errors.New("out of buffer")
	}

	m := copy(out, ciphertext)

	return m, nil
}

func (p *ploy) Uncompress(src []byte, key []byte, out []byte) (int, error) {
	defer utils.Trace("Uncompress")()
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return 0, err
	}

	plaintext, err := aead.Open(nil, p.iv, src, nil)
	if len(plaintext) == 0 || err != nil {
		return 0, errors.New("compressed fail")
	}

	n := len(plaintext)
	if len(out) < n {
		return 0, errors.New("out of buffer")
	}

	m := copy(out, plaintext)

	return m, nil
}

func (p *ploy) MakeSalt() []byte {
	nonce := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil
	}
	return nonce
}

func (p *ploy) GetIV() []byte {
	return p.iv
}

func (p *ploy) SetIV(iv []byte) {
	p.iv = make([]byte, len(iv))
	copy(p.iv, iv)
}

func (p *ploy) GetSize() int {
	return 2 + 1 + len(p.iv)
}

func (p *ploy) ToBytes() []byte {
	return methodToBytes(ployCode, p.iv)
}

// From bytes
func (p *ploy) FromBytes(v []byte) error {
	iv, err := methodFromBytes(v)
	if err != nil {
		return err
	}
	if iv != nil {
		p.SetIV(iv)
	}
	return nil
}

func init() {
	register(ployName, ployCode, &ploy{})
}
