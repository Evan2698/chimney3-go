package privacy

import (
	"errors"

	"golang.org/x/crypto/salsa20"
)

type salsa_20 struct {
	iv []byte
}

const (
	salsaName = "SALSA20-I"
	salsaCode = 0x1238
)

func (g *salsa_20) Compress(src []byte, key []byte, out []byte) (int, error) {
	if g == nil {
		return 0, errors.New("privacy: nil receiver")
	}

	n := len(src)
	if n == 0 {
		return 0, errors.New("compressed failed")
	}

	if len(out) < n {
		return 0, errors.New("out of buffer")
	}

	if len(key) != 32 {
		return 0, errors.New("key length must be 32 bytes")
	}
	if len(g.iv) != 24 {
		return 0, errors.New("IV length must be 24 bytes")
	}
	var keyArr [32]byte

	copy(keyArr[:], key)

	salsa20.XORKeyStream(out, src, g.iv[:], &keyArr)

	return n, nil
}

func (g *salsa_20) Uncompress(src []byte, key []byte, out []byte) (int, error) {
	if g == nil {
		return 0, errors.New("privacy: nil receiver")
	}

	return g.Compress(src, key, out)
}

func (g *salsa_20) MakeSalt() []byte {
	if g == nil {
		return nil
	}
	return randomBytes(24)
}

func (g *salsa_20) GetIV() []byte {
	if g == nil {
		return nil
	}
	return cloneBytes(g.iv)
}

func (g *salsa_20) SetIV(iv []byte) {
	if g == nil {
		return
	}
	g.iv = cloneBytes(iv)
}

func (g *salsa_20) GetSize() int {
	if g == nil {
		return 0
	}
	return 2 + 1 + len(g.iv)
}

func (g *salsa_20) ToBytes() []byte {
	if g == nil {
		return nil
	}
	return methodToBytes(salsaCode, g.iv)
}

// From bytes
func (g *salsa_20) FromBytes(v []byte) error {
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
	register(salsaName, salsaCode, &salsa_20{})
}
