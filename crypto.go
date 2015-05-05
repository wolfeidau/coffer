package coffer

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/secretbox"
)

func decrypt(ct, key []byte) []byte {

	var k [32]byte
	var nonce [24]byte
	var opened []byte

	// extract the nonce from the start of the message
	copy(nonce[:], ct[:24])

	copy(k[:], key[:32])

	Debugf("nonce=%x key=%x box=%x", nonce, k, ct[24:])

	// out, box, nonce, key
	var ok bool
	opened, ok = secretbox.Open(opened[:0], ct[24:], &nonce, &k)

	if !ok {
		Fatalf("Failed to decrypt data")
	}

	return opened
}

func encrypt(plaintext, key []byte) []byte {

	var k [32]byte
	var nonce [24]byte
	var box []byte

	// must be unique for each encrypt
	rand.Reader.Read(nonce[:])
	copy(k[:], key[:32])

	// out, message, nonce, key
	box = secretbox.Seal(box[:0], plaintext, &nonce, &k)

	Debugf("nonce=%x key=%x box=%x", nonce, k, box)

	// add the nonce to the start of the message
	box = append(nonce[:], box...)

	return box
}

func generateNonce(ln int) []byte {

	b := make([]byte, ln)
	_, err := rand.Read(b)
	if err != nil {
		Fatalf("Failed to generate random nonce: %v", err)
	}
	return b
}
