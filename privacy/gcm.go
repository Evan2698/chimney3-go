package privacy

import (
	"chimney3-go/utils"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type gcm struct {
	iv []byte
}

const (
	gcmName = "AES-GCM"
	gcmCode = 0x1234
)

func (g *gcm) Compress(src []byte, key []byte, out []byte) (int, error) {
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
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil
	}
	return nonce
}

func (g *gcm) GetIV() []byte {
	return g.iv
}

func (g *gcm) SetIV(iv []byte) {
	g.iv = make([]byte, len(iv))
	copy(g.iv, iv)
}

func (g *gcm) GetSize() int {
	return 2 + 1 + len(g.iv)
}

func (g *gcm) ToBytes() []byte {
	return methodToBytes(gcmCode, g.iv)
}

// From bytes
func (g *gcm) FromBytes(v []byte) error {
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
