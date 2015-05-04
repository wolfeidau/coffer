package coffer

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

const nonceLen = 12

func decrypt(ct, key, ad []byte) []byte {

	nonce := ct[:nonceLen]
	ct = ct[nonceLen:]

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

	// must be unique for each encrypt
	nonce := generateNonce()

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

	return bytes.Join([][]byte{nonce, ct}, []byte{})
}

func generateNonce() []byte {

	b := make([]byte, nonceLen)
	_, err := rand.Read(b)
	if err != nil {
		Fatalf("Failed to generate random nonce: %v", err)
	}
	return b
}
