package coffer

import (
	"bytes"
	"encoding/hex"
	"os"
	"path"
)

var (
	CofferFilePrefix = []byte(`COFFER:AES256:`)
	CofferBlockSize  = 32                // AES256
	OwnerRead        = os.FileMode(0600) // os.FileMode, note the octal number
)

func MustEncrypt(cofferFile, secret string) {

	data := mustReadFile(cofferFile)

	payload := mustEncryptPayload(data, secret)

	mustWriteFile(cofferFile, payload, OwnerRead)
}

func MustDecrypt(cofferFile, secret string) []byte {

	data := mustReadFile(cofferFile)

	payload := mustDecryptPayload(data, secret)

	return mustWriteFile(cofferFile, payload, OwnerRead)

}

func MustSync(cofferFile, secret, base string) {

	// base not set
	if base == "" {
		base, _ = os.Getwd() // assign the current working directory
	}

	payload := mustReadFile(cofferFile)

	// if the coffer file is encrypted, decrypt it
	if isEncrypted(payload) {
		payload = mustDecryptPayload(payload, secret)
	}

	bundle := mustDecodeBundle(payload)

	bundle.MustValidate()

	mustExtractBundle(bundle, base)
}

func MustUpload(cofferFile, secret, bucket string) {

	payload := mustReadFile(cofferFile)

	// if the coffer is not encrypted
	if !isEncrypted(payload) {
		payload = mustEncryptPayload(payload, secret)
	}

	filename := path.Base(cofferFile)

	mustUpload(bucket, filename, payload)
}

func MustDownload(cofferFile, secret, bucket string) {

	filename := path.Base(cofferFile)

	payload := mustDownload(bucket, filename)

	mustWriteFile(cofferFile, payload, OwnerRead)
}

func MustDownloadSync(cofferFile, secret, bucket, base string) {

	// base not set
	if base == "" {
		base, _ = os.Getwd() // assign the current working directory
	}

	payload := mustDownload(bucket, filename)

	// if the coffer file is encrypted, decrypt it
	if isEncrypted(payload) {
		payload = mustDecryptPayload(payload, secret)
	}

	bundle := mustDecodeBundle(payload)

	bundle.MustValidate()

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

func mustEncryptPayload(data []byte, secret string) []byte {

	key := buildKey(secret)

	if bytes.HasPrefix(data, CofferFilePrefix) {
		Fatalf("coffer file alread encrypted")
	}

	payload := encrypt(data, key)
	encoded := make([]byte, hex.EncodedLen(len(payload)))

	n := hex.Encode(encoded, payload)

	Infof("encoded data len: %d", n)

	return bytes.Join([][]byte{CofferFilePrefix, encoded, []byte("\n")}, []byte{})
}

func mustDecryptPayload(data []byte, secret string) []byte {

	key := buildKey(secret)

	if bytes.HasPrefix(data, CofferFilePrefix) {

		payload := bytes.TrimPrefix(data, CofferFilePrefix)
		payload = bytes.TrimSpace(payload) // remove any trailing whitespace
		decoded := make([]byte, hex.DecodedLen(len(payload)))

		n, err := hex.Decode(decoded, payload)
		if err != nil {
			Fatalf("coffer file could not be decoded")
		}

		Infof("decoded data len: %d", n)
		return decrypt(decoded, key)
	}

	Fatalf("unable to decrypt coffer")

	return []byte{}
}
