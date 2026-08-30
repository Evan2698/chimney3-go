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

func (raw *rawMethod) MakeSalt() []byte {
	if raw == nil {
		return nil
	}
	return randomBytes(24)
}

func (raw *rawMethod) GetIV() []byte {
	if raw == nil {
		return nil
	}
	return cloneBytes(raw.iv)
}

func (raw *rawMethod) SetIV(iv []byte) {
	if raw == nil {
		return
	}
	raw.iv = cloneBytes(iv)
}

func (raw *rawMethod) GetSize() int {
	if raw == nil {
		return 0
	}
	return 2 + 1 + len(raw.iv)
}

func (raw *rawMethod) ToBytes() []byte {
	if raw == nil {
		return nil
	}
	return methodToBytes(rawCode, raw.iv)
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
		raw.SetIV(iv)
	}
	return nil
}

func init() {
	register(rawName, rawCode, &rawMethod{})
}
