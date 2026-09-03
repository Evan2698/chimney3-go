package privacy

import (
	"chimney3-go/utils"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"log"
	"reflect"
	"strings"
)

type Encryptor interface {
	// RecommendedBufferSize returns the recommended buffer size for the given source length.
	RecommendedBufferSize(srcLen int) int

	// Compress encrypts the source data using the provided key and writes the result to the output buffer.
	// It returns the number of bytes written to the output buffer and any error encountered.
	Compress(src []byte, key []byte, out []byte) (int, error)

	// Uncompress decrypts the source data using the provided key and writes the result to the output buffer.
	// It returns the number of bytes written to the output buffer and any error encountered.
	Uncompress(src []byte, key []byte, out []byte) (int, error)
}

type IVManager interface {
	// GenerateIV generates a new initialization vector (IV) for encryption.
	// It returns the generated IV and any error encountered.
	GenerateIV() ([]byte, error)

	// SetIV sets the initialization vector (IV) for encryption.
	// It returns any error encountered.
	SetIV(iv []byte) error

	// GetIV returns the current initialization vector (IV) used for encryption.
	// It returns the IV and any error encountered.
	GetIV() []byte
}

type SaltManager interface {
	// GenerateSalt generates a new salt value for encryption.
	// It returns the generated salt and any error encountered.

	GenerateSalt() ([]byte, error)
	// GetSalt returns the current salt value used for encryption.
	// It returns the salt and any error encountered.
	GetSalt() []byte
}

type Serializable interface {
	// ToBytes serializes the object to a byte slice.
	// It returns the serialized byte slice and any error encountered.
	ToBytes() ([]byte, error)

	// FromBytes deserializes the object from a byte slice.
	// It returns any error encountered during deserialization.
	FromBytes(data []byte) error

	// PartialSerializeSize returns the size of the serialized data without actually serializing it.
	PartialSerializeSize() int
}

// EncryptThings for everything protecting
type EncryptThings interface {
	Encryptor
	IVManager
	SaltManager
	Serializable
}

var globalTable map[string]interface{} = make(map[string]interface{})
var globalTablei map[uint16]interface{} = make(map[uint16]interface{})

func normalizeMethodName(name string) string {
	return strings.ToUpper(strings.TrimSpace(name))
}

func register(name string, mask uint16, i interface{}) {
	key := normalizeMethodName(name)
	globalTable[key] = i
	globalTablei[mask] = i
}

// FromBytes ...
func FromBytes(buf []byte) (EncryptThings, error) {
	if buf == nil {
		return nil, errors.New("invalid paramter")
	}

	code := utils.Bytes2Uint16(buf[:2])
	i := newMethodWithCode(code)
	if i == nil {
		return nil, errors.New("use code create method failed")
	}
	if err := i.FromBytes(buf[2:]); err != nil {
		return nil, err
	}
	return i, nil
}

// NewMethodWithName create encrypt method for caller with a name
func NewMethodWithName(name string) EncryptThings {
	if target, ok := globalTable[normalizeMethodName(name)]; ok {
		return createObject(target)
	}
	return nil
}

func createObject(target interface{}) EncryptThings {
	t := reflect.New(reflect.TypeOf(target).Elem()).Elem().Addr().Interface()
	if i, ok := t.(EncryptThings); ok {
		salt, err := i.GenerateSalt()
		if err != nil {
			log.Println("create encrypt method failed!!")
			return nil
		}

		i.SetIV(salt)
		return i
	}
	log.Println("create encrypt method failed!!")
	return nil
}

func newMethodWithCode(code uint16) EncryptThings {
	if target, ok := globalTablei[code]; ok {
		return createObject(target)
	}
	return nil
}

const (
	// CompressionKeySalt is a constant string used as a salt for deriving compression keys.
	CompressionKeySalt = "d7722deb18976aa66e5eb70cb804b0ee"
)

// MakeCompressKey ..
func DeriveCompressionKey(srcKey string) []byte {
	return ComputeHMACSHA256([]byte(CompressionKeySalt), srcKey)
}

// BuildMacHash ..
func ComputeHMACSHA256(key []byte, message string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return h.Sum(nil)
}
