package privacy

import (
	"chimney3-go/utils"

	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type gcm struct {
	iv []byte
}

const (
	gcmName       = "AES-GCM"
	gcmCode       = 0x1234
	gcmBlockSize  = 16
	gcmNonceSize  = 12
	gcmTagSize    = 16
	gcmMinTagSize = 12
)

func newAESGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if block.BlockSize() != gcmBlockSize {
		return nil, errors.New("privacy: AES-GCM requires 128-bit block size")
	}
	return cipher.NewGCM(block)
}

func (g *gcm) RecommendedBufferSize(srcLen int) int {
	if g == nil {
		return srcLen
	}
	return srcLen + gcmTagSize
}

func (g *gcm) Compress(src []byte, key []byte, out []byte) (int, error) {
	if g == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	defer utils.Trace("Compress")()

	aesgcm, err := newAESGCM(key)
	if err != nil {
		return 0, err
	}

	if len(g.iv) == 0 {
		return 0, errors.New("privacy: empty IV")
	}
	if len(out) < len(src)+aesgcm.Overhead() {
		return 0, errors.New("out of buffer")
	}

	ciphertext := aesgcm.Seal(out[:0], g.iv, src, nil)
	if len(ciphertext) == 0 {
		return 0, errors.New("compressed failed")
	}
	return len(ciphertext), nil
}

func (g *gcm) Uncompress(src []byte, key []byte, out []byte) (int, error) {
	if g == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	defer utils.Trace("Uncompress")()

	aesgcm, err := newAESGCM(key)
	if err != nil {
		return 0, err
	}

	if len(g.iv) == 0 {
		return 0, errors.New("privacy: empty IV")
	}
	if len(src) < aesgcm.Overhead() {
		return 0, errors.New("ciphertext too short")
	}

	plaintext, err := aesgcm.Open(nil, g.iv, src, nil)
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

func (g *gcm) GenerateIV() ([]byte, error) {
	if g == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	iv := randomBytes(gcmNonceSize)
	if iv == nil {
		return nil, errors.New("generate IV failed")
	}
	return iv, nil
}

func (g *gcm) GenerateSalt() ([]byte, error) {
	return g.GenerateIV()
}

func (g *gcm) GetIV() []byte {
	if g == nil {
		return nil
	}
	return cloneBytes(g.iv)
}

func (g *gcm) GetSalt() []byte {
	return g.GetIV()
}

func (g *gcm) SetIV(iv []byte) error {
	if g == nil {
		return errors.New("privacy: nil receiver")
	}
	if len(iv) != gcmNonceSize {
		return errors.New("IV length must be 12 bytes")
	}
	g.iv = cloneBytes(iv)
	return nil
}

func (g *gcm) PartialSerializeSize() int {
	if g == nil {
		return 0
	}
	return 2 + 1 + len(g.iv)
}

func (g *gcm) BlockSize() int {
	if g == nil {
		return 0
	}
	return gcmBlockSize
}

func (g *gcm) ToBytes() ([]byte, error) {
	if g == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	return methodToBytes(gcmCode, g.iv), nil
}

// From bytes
func (g *gcm) FromBytes(v []byte) error {
	if g == nil {
		return errors.New("privacy: nil receiver")
	}
	iv, err := methodFromBytes(v)
	if err != nil {
		return err
	}
	if iv != nil {
		if err := g.SetIV(iv); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	register(gcmName, gcmCode, &gcm{})
}
