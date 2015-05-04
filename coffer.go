package coffer

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
)

var (
	CofferFilePrefix = []byte(`COFFER:AES256:`)
	CofferBlockSize  = 32 // AES256
)

func MustEncrypt(cofferFile string, secret string) {

	key := buildKey(secret)

	data := mustReadFile(cofferFile)

	payload := mustEncryptPayload(data, key)

	mustWriteFile(cofferFile, payload)
}

func MustDecrypt(cofferFile string, secret string) {

	key := buildKey(secret)

	data := mustReadFile(cofferFile)

	payload := mustDecryptPayload(data, key)

	mustWriteFile(cofferFile, payload)

}

func mustReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		Fatalf("Unable to open coffer-file: %v", err)
	}
	return data
}

func mustWriteFile(path string, data []byte) []byte {
	err := ioutil.WriteFile(path, data, 600)

	if err != nil {
		Fatalf("Unable to open coffer-file: %v", err)
	}
	return data
}

func buildKey(secret string) []byte {

	if len(secret) > CofferBlockSize {
		Infof("secret is longer than block size and will be trucated")
	}

	padded := make([]byte, CofferBlockSize)
	copy(padded, []byte(secret))

	return padded
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
