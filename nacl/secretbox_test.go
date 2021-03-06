package nacl

import (
	"crypto/rand"
	"testing"

	"github.com/bmizerany/assert"
)

var msg = []byte(`
test
test
test
`)

func TestEncryptions(t *testing.T) {

	var key [32]byte
	var nonce [24]byte

	rand.Reader.Read(key[:])
	rand.Reader.Read(nonce[:])

	enc := Encrypt(msg, key[:])

	dec := Decrypt(enc, key[:])

	assert.Equal(t, msg, dec)
}
