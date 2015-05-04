package coffer

import (
	"bytes"
	"encoding/hex"
)

var (
	CofferFilePrefix = []byte(`COFFER:AES256:`)
	CofferBlockSize  = 32 // AES256
)

func MustEncrypt(cofferFile, secret string) {

	key := buildKey(secret)

	data := mustReadFile(cofferFile)

	payload := mustEncryptPayload(data, key)

	mustWriteFile(cofferFile, payload, 0600)
}

func MustDecrypt(cofferFile, secret string) []byte {

	key := buildKey(secret)

	data := mustReadFile(cofferFile)

	payload := mustDecryptPayload(data, key)

	return mustWriteFile(cofferFile, payload, 0600)

}

func MustSync(cofferFile, secret, base string) {

	var payload []byte

	payload = mustReadFile(cofferFile)

	// if the coffer file is encrypted, decrypt it
	if isEncrypted(payload) {
		payload = MustDecrypt(cofferFile, secret)
	}

	bundle := mustDecodeBundle(payload)

	mustExtractBundle(bundle, base)

}

func buildKey(secret string) []byte {

	if len(secret) > CofferBlockSize {
		Infof("secret is longer than block size and will be trucated")
	}

	padded := make([]byte, CofferBlockSize)
	copy(padded, []byte(secret))

	return padded
}

func isEncrypted(data []byte) bool {
	return bytes.HasPrefix(data, CofferFilePrefix)
}

func mustEncryptPayload(data []byte, key []byte) []byte {

	if bytes.HasPrefix(data, CofferFilePrefix) {
		Fatalf("coffer file alread encrypted")
	}

	payload := encrypt(data, key, CofferFilePrefix)
	encoded := make([]byte, hex.EncodedLen(len(payload)))

	n := hex.Encode(encoded, payload)

	Infof("encoded data len: %d", n)

	return bytes.Join([][]byte{CofferFilePrefix, encoded, []byte("\n")}, []byte{})
}

func mustDecryptPayload(data []byte, key []byte) []byte {
	if bytes.HasPrefix(data, CofferFilePrefix) {

		payload := bytes.TrimPrefix(data, CofferFilePrefix)
		payload = bytes.TrimSpace(payload) // remove any trailing whitespace
		decoded := make([]byte, hex.DecodedLen(len(payload)))

		n, err := hex.Decode(decoded, payload)
		if err != nil {
			Fatalf("coffer file could not be decoded")
		}

		Infof("decoded data len: %d", n)
		return decrypt(decoded, key, CofferFilePrefix)
	}

	Fatalf("unable to decrypt coffer")

	return []byte{}
}
