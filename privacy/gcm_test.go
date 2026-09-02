package privacy

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncryptDecryptGCM(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("Hello, World! This is a test message for GCM encryption.")
	s := &gcm{}
	salt, _ := s.GenerateSalt()
	s.SetIV(salt)

	ciphertext := make([]byte, len(plaintext)+16)
	n, err := s.Compress(plaintext, key, ciphertext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	if n != len(plaintext)+16 {
		t.Fatalf("unexpected ciphertext length: got %d want %d", n, len(plaintext)+16)
	}

	decrypted := make([]byte, len(plaintext))
	n, err = s.Uncompress(ciphertext[:n], key, decrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	if n != len(plaintext) {
		t.Fatalf("unexpected plaintext length: got %d want %d", n, len(plaintext))
	}
	if !bytes.Equal(plaintext, decrypted) {
		t.Fatal("Decrypted text does not match original")
	}

	t.Logf("Original: %s", plaintext)
	t.Logf("Ciphertext: %x", ciphertext[:n])
	t.Logf("Decrypted: %s", decrypted)
}
