package nacl

import (
	"crypto/rand"
	"log"

	"golang.org/x/crypto/nacl/secretbox"
)

// Decrypt basic wrapper around secretbox which will decrypt a box
// using ct which is comprised of a nonce followed by the box.
func Decrypt(ct, key []byte) []byte {

	mustValidateKey(key)

	var k [32]byte
	var nonce [24]byte
	var opened []byte

	// extract the nonce from the start of the message
	copy(nonce[:], ct[:24])

	copy(k[:], key[:32])

	// out, box, nonce, key
	var ok bool
	opened, ok = secretbox.Open(opened[:0], ct[24:], &nonce, &k)

	if !ok {
		log.Fatalln("Failed to decrypt data")
	}

	return opened
}

// Encrypt basic wrapper around secretbox which will encrypt a plain text
// and return a message comprised of the nonce followed by the encrypted box.
func Encrypt(plaintext, key []byte) []byte {

	mustValidateKey(key)

	var k [32]byte
	var nonce [24]byte
	var box []byte

	// must be unique for each encrypt
	rand.Reader.Read(nonce[:])
	copy(k[:], key[:32])

	// out, message, nonce, key
	box = secretbox.Seal(box[:0], plaintext, &nonce, &k)

	// add the nonce to the start of the message
	box = append(nonce[:], box...)

	return box
}

func mustValidateKey(key []byte) {
	if len(key) < 32 {
		log.Fatalln("Key validatation failed, must be 32 bytes long")
	}
}
