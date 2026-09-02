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

func (g *salsa_20) RecommendedBufferSize(srcLen int) int {
	if g == nil {
		return srcLen
	}
	return srcLen
}

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

func (g *salsa_20) GenerateIV() ([]byte, error) {
	if g == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	iv := randomBytes(24)
	if iv == nil {
		return nil, errors.New("generate IV failed")
	}
	return iv, nil
}

func (g *salsa_20) GenerateSalt() ([]byte, error) {
	return g.GenerateIV()
}

func (g *salsa_20) GetIV() []byte {
	if g == nil {
		return nil
	}
	return cloneBytes(g.iv)
}

func (g *salsa_20) GetSalt() []byte {
	return g.GetIV()
}

func (g *salsa_20) SetIV(iv []byte) error {
	if g == nil {
		return errors.New("privacy: nil receiver")
	}
	if len(iv) != 24 {
		return errors.New("IV length must be 24 bytes")
	}
	g.iv = cloneBytes(iv)
	return nil
}

func (g *salsa_20) PartialSerializeSize() int {
	if g == nil {
		return 0
	}
	return 2 + 1 + len(g.iv)
}

func (g *salsa_20) ToBytes() ([]byte, error) {
	if g == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	return methodToBytes(salsaCode, g.iv), nil
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
		if err := g.SetIV(iv); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	register(salsaName, salsaCode, &salsa_20{})
}
