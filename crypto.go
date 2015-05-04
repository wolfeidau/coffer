package coffer

import (
	"crypto/aes"
	"crypto/cipher"
)

// need to work out what the issues are with hard coding this, nonce isn't really that
// helpful in my case I don't believe.
var nonce = []byte("3c819d9a9bed")

func decrypt(ct, key, ad []byte) []byte {

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		Fatalf("Cannot load cipher: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		Fatalf("New GCM failed: %v", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ct, ad)
	if err != nil {
		Fatalf("Open payload failed: %v", err)
	}

	return plaintext
}

func encrypt(plaintext, key, ad []byte) []byte {

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		Fatalf("Cant load cipher: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		Fatalf("New GCM failed: %v", err)
	}

	ct := aesgcm.Seal(nil, nonce, plaintext, ad)

	return ct
}
