package privacy

import (
	"chimney3-go/utils"
	"errors"
)

type rawMethod struct {
	iv []byte
}

const (
	rawName = "RAW"
	rawCode = 0x1237
)

func (raw *rawMethod) RecommendedBufferSize(srcLen int) int {
	if raw == nil {
		return srcLen
	}
	return srcLen
}

func (raw *rawMethod) Compress(src []byte, key []byte, out []byte) (int, error) {
	if raw == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	defer utils.Trace("Compress")()
	n := copy(out, src)
	return n, nil
}

func (raw *rawMethod) Uncompress(src []byte, key []byte, out []byte) (int, error) {
	if raw == nil {
		return 0, errors.New("privacy: nil receiver")
	}
	return raw.Compress(src, key, out)
}

func (raw *rawMethod) GenerateIV() ([]byte, error) {
	if raw == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	iv := randomBytes(24)
	if iv == nil {
		return nil, errors.New("generate IV failed")
	}
	return iv, nil
}

func (raw *rawMethod) GenerateSalt() ([]byte, error) {
	return raw.GenerateIV()
}

func (raw *rawMethod) GetIV() []byte {
	if raw == nil {
		return nil
	}
	return cloneBytes(raw.iv)
}

func (raw *rawMethod) SetIV(iv []byte) error {
	if raw == nil {
		return errors.New("privacy: nil receiver")
	}
	raw.iv = cloneBytes(iv)
	return nil
}

func (raw *rawMethod) PartialSerializeSize() int {
	if raw == nil {
		return 0
	}
	return 2 + 1 + len(raw.iv)
}

func (raw *rawMethod) ToBytes() ([]byte, error) {
	if raw == nil {
		return nil, errors.New("privacy: nil receiver")
	}
	return methodToBytes(rawCode, raw.iv), nil
}

// From bytes
func (raw *rawMethod) FromBytes(v []byte) error {
	if raw == nil {
		return errors.New("privacy: nil receiver")
	}
	iv, err := methodFromBytes(v)
	if err != nil {
		return err
	}
	if iv != nil {
		if err := raw.SetIV(iv); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	register(rawName, rawCode, &rawMethod{})
}
