package privacy

import (
	"chimney3-go/privacy/chacha20"
	"errors"
)

type cha20 struct {
	iv []byte
}

const (
	chacha20Name = "CHACHA-20"
	chacha20Code = 0x1235
)

func (chacha *cha20) RecommendedBufferSize(srcLen int) int {
	if chacha == nil {
		return srcLen
	}
	return srcLen
}

func (chacha *cha20) Compress(src []byte, key []byte, out []byte) (int, error) {
	if chacha == nil {
		return 0, errors.New("privacy: nil receiver")
	}

	if len(key) != 32 || len(src) == 0 {
		return 0, errors.New("parameter is invalid")
	}
	if len(chacha.iv) != 24 {
		return 0, errors.New("IV length must be 24 bytes")
	}

	a, err := chacha20.NewXChaCha(key, chacha.iv)
	if err != nil {
		return 0, err
	}

	a.XORKeyStream(out, src)

	return len(src), nil

}

func (chacha *cha20) Uncompress(src []byte, key []byte, out []byte) (int, error) {
	if chacha == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	return chacha.Compress(src, key, out)
}

func (chacha *cha20) GenerateIV() ([]byte, error) {
	if chacha == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	iv := randomBytes(24)
	if iv == nil {
		return nil, errors.New("generate IV failed")
	}
	return iv, nil
}

func (chacha *cha20) GenerateSalt() ([]byte, error) {
	return chacha.GenerateIV()
}

func (chacha *cha20) GetIV() []byte {
	if chacha == nil {
		return nil
	}
	return cloneBytes(chacha.iv)
}

func (chacha *cha20) SetIV(iv []byte) error {
	if chacha == nil {
		return errors.New("privacy: nil receiver")
	}
	if len(iv) != 24 {
		return errors.New("IV length must be 24 bytes")
	}
	chacha.iv = cloneBytes(iv)
	return nil
}
func (chacha *cha20) GetSalt() []byte {
	return chacha.GetIV()
}

func (chacha *cha20) PartialSerializeSize() int {
	if chacha == nil {
		return 0
	}
	return 2 + 1 + len(chacha.iv)
}

func (chacha *cha20) ToBytes() ([]byte, error) {
	if chacha == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	return methodToBytes(chacha20Code, chacha.iv), nil
}

// From bytes
func (chacha *cha20) FromBytes(v []byte) error {
	if chacha == nil {
		return errors.New("privacy: nil receiver")
	}
	iv, err := methodFromBytes(v)
	if err != nil {
		return err
	}
	if iv != nil {
		if err := chacha.SetIV(iv); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	register(chacha20Name, chacha20Code, &cha20{})
}
