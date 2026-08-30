package privacy

import (
	"chimney3-go/utils"
	"errors"

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
	if p == nil {
		return 0, errors.New("privacy: nil receiver")
	}
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
	if p == nil {
		return 0, errors.New("privacy: nil receiver")
	}
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
	if p == nil {
		return nil
	}
	return randomBytes(24)
}

func (p *ploy) GetIV() []byte {
	if p == nil {
		return nil
	}
	return cloneBytes(p.iv)
}

func (p *ploy) SetIV(iv []byte) {
	if p == nil {
		return
	}
	p.iv = cloneBytes(iv)
}

func (p *ploy) GetSize() int {
	if p == nil {
		return 0
	}
	return 2 + 1 + len(p.iv)
}

func (p *ploy) ToBytes() []byte {
	if p == nil {
		return nil
	}
	return methodToBytes(ployCode, p.iv)
}

// From bytes
func (p *ploy) FromBytes(v []byte) error {
	if p == nil {
		return errors.New("privacy: nil receiver")
	}
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
