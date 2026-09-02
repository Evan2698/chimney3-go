package privacy

import (
	"chimney3-go/utils"
	"crypto/cipher"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

type ploy struct {
	iv []byte
}

const (
	ployName      = "CHACHA-POLY1305"
	ployCode      = 0x1236
	ployNonceSize = 24
	ployTagSize   = 16
)

func newChaCha20Poly1305(key []byte) (cipher.AEAD, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	return aead, nil
}

func (p *ploy) RecommendedBufferSize(srcLen int) int {
	if p == nil {
		return srcLen
	}
	return srcLen + chacha20poly1305.Overhead
}

func (p *ploy) Compress(src []byte, key []byte, out []byte) (int, error) {
	if p == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	defer utils.Trace("Compress")()

	if len(key) != chacha20poly1305.KeySize {
		return 0, errors.New("key length must be 32 bytes")
	}
	if len(p.iv) != chacha20poly1305.NonceSizeX {
		return 0, errors.New("IV length must be 24 bytes")
	}

	aead, err := newChaCha20Poly1305(key)
	if err != nil {
		return 0, err
	}
	if len(out) < len(src)+aead.Overhead() {
		return 0, errors.New("out of buffer")
	}

	ciphertext := aead.Seal(out[:0], p.iv, src, nil)
	if len(ciphertext) == 0 {
		return 0, errors.New("compressed failed")
	}
	return len(ciphertext), nil
}

func (p *ploy) Uncompress(src []byte, key []byte, out []byte) (int, error) {
	if p == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	defer utils.Trace("Uncompress")()

	if len(key) != chacha20poly1305.KeySize {
		return 0, errors.New("key length must be 32 bytes")
	}
	if len(p.iv) != chacha20poly1305.NonceSizeX {
		return 0, errors.New("IV length must be 24 bytes")
	}
	if len(src) < chacha20poly1305.Overhead {
		return 0, errors.New("ciphertext too short")
	}

	aead, err := newChaCha20Poly1305(key)
	if err != nil {
		return 0, err
	}

	plaintext, err := aead.Open(nil, p.iv, src, nil)
	if err != nil {
		return 0, err
	}
	if len(plaintext) == 0 {
		return 0, errors.New("compressed failed")
	}
	if len(out) < len(plaintext) {
		return 0, errors.New("out of buffer")
	}

	m := copy(out, plaintext)
	return m, nil
}

func (p *ploy) GenerateIV() ([]byte, error) {
	if p == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	iv := randomBytes(ployNonceSize)
	if iv == nil {
		return nil, errors.New("generate IV failed")
	}
	return iv, nil
}

func (p *ploy) GenerateSalt() ([]byte, error) {
	return p.GenerateIV()
}

func (p *ploy) GetIV() []byte {
	if p == nil {
		return nil
	}
	return cloneBytes(p.iv)
}

func (p *ploy) SetIV(iv []byte) error {
	if p == nil {
		return errors.New("privacy: nil receiver")
	}
	if len(iv) != ployNonceSize {
		return errors.New("IV length must be 24 bytes")
	}
	p.iv = cloneBytes(iv)
	return nil
}

func (p *ploy) PartialSerializeSize() int {
	if p == nil {
		return 0
	}
	return 2 + 1 + len(p.iv)
}

func (p *ploy) ToBytes() ([]byte, error) {
	if p == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	return methodToBytes(ployCode, p.iv), nil
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
		if err := p.SetIV(iv); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	register(ployName, ployCode, &ploy{})
}
