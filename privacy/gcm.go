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
	gcmName = "AES-GCM"
	gcmCode = 0x1234
)

func (g *gcm) Compress(src []byte, key []byte, out []byte) (int, error) {
	if g == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	defer utils.Trace("Compress")()

	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, err
	}

	ciphertext := aesgcm.Seal(nil, g.iv, src, nil)
	n := len(ciphertext)
	if n == 0 {
		return 0, errors.New("compressed failed")
	}

	if len(out) < n {
		return 0, errors.New("out of buffer")
	}

	m := copy(out, ciphertext)

	return m, nil
}

func (g *gcm) Uncompress(src []byte, key []byte, out []byte) (int, error) {
	if g == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	defer utils.Trace("Uncompress")()

	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, err
	}

	plaintext, err := aesgcm.Open(nil, g.iv, src, nil)
	n := len(plaintext)

	if n == 0 {
		return 0, errors.New("compressed failed")
	}

	if len(out) < n {
		return 0, errors.New("out of buffer")
	}

	m := copy(out, plaintext)

	return m, err
}

func (g *gcm) MakeSalt() []byte {
	if g == nil {
		return nil
	}
	return randomBytes(12)
}

func (g *gcm) GetIV() []byte {
	if g == nil {
		return nil
	}
	return cloneBytes(g.iv)
}

func (g *gcm) SetIV(iv []byte) {
	if g == nil {
		return
	}
	g.iv = cloneBytes(iv)
}

func (g *gcm) GetSize() int {
	if g == nil {
		return 0
	}
	return 2 + 1 + len(g.iv)
}

func (g *gcm) ToBytes() []byte {
	if g == nil {
		return nil
	}
	return methodToBytes(gcmCode, g.iv)
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
		g.SetIV(iv)
	}
	return nil
}

func init() {
	register(gcmName, gcmCode, &gcm{})
}
