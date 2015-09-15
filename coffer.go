package coffer

import (
	"encoding/base64"
	"log"
	"os"
	"path"

	"github.com/wolfeidau/coffer/nacl"
)

var (
	// CofferBlockSize size of the key
	CofferBlockSize = 32
	// OwnerRead is the default mode set for new coffer files, note the octal number
	OwnerRead = os.FileMode(0600)
	// Version the version of the coffer file
	Version = "2.0.0"
)

// MustEncrypt encrypt the supplied file
func MustEncrypt(cofferFile, secret string) {

	data := mustReadFile(cofferFile)

	payload := mustEncryptPayload(data, secret)

	mustWriteFile(cofferFile, payload, OwnerRead)
}

// MustDecrypt decrypt the supplied file
func MustDecrypt(cofferFile, secret string) []byte {

	data := mustReadFile(cofferFile)

	payload := mustDecryptPayload(data, secret)

	return mustWriteFile(cofferFile, payload, OwnerRead)

}

// MustSync sync the file up to S3
func MustSync(cofferFile, secret, base string) {

	// base not set
	if base == "" {
		base = "/" // set the base directory to '/'
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

// MustUpload upload the file to the supplied s3 bucket
func MustUpload(cofferFile, secret, bucket string) {

	payload := mustReadFile(cofferFile)

	// if the coffer is not encrypted
	if !isEncrypted(payload) {
		payload = mustEncryptPayload(payload, secret)
	}

	filename := path.Base(cofferFile)

	mustUpload(bucket, filename, payload)
}

// MustDownload download the file from the supplied s3 bucket
func MustDownload(cofferFile, secret, bucket string) {

	filename := path.Base(cofferFile)

	payload := mustDownload(bucket, filename)

	mustWriteFile(cofferFile, payload, OwnerRead)
}

// MustDownloadSync download the file from the supplied s3 bucket and sync
// it to the filesystem
func MustDownloadSync(cofferFile, secret, bucket, base string) {

	// base not set
	if base == "" {
		base = "/" // set the base directory to '/'
	}

	filename := path.Base(cofferFile)

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
		log.Printf("secret is longer than block size and will be trucated")
	}

	padded := make([]byte, CofferBlockSize)
	copy(padded, []byte(secret))

	return padded
}

func isEncrypted(data []byte) bool {

	coffer, err := DecodeCoffer(data)

	if err != nil {
		return false
	}

	return coffer.Validate()
}

func mustEncryptPayload(data []byte, secret string) []byte {

	key := buildKey(secret)

	if isEncrypted(data) {
		log.Fatalf("coffer file already encrypted")
	}

	payload := nacl.Encrypt(data, key)
	encoded := base64.StdEncoding.EncodeToString(payload)

	log.Printf("encoded data len: %d", len(encoded))

	return mustEncodeCoffer(&Coffer{
		Version:    Version,
		CipherText: encoded,
	})
}

func mustDecryptPayload(data []byte, secret string) []byte {

	key := buildKey(secret)

	coffer, err := DecodeCoffer(data)

	if err != nil {
		log.Fatalf("unable to decode coffer")
	}

	if coffer.Validate() {

		decoded, err := base64.StdEncoding.DecodeString(coffer.CipherText)

		if err != nil {
			log.Fatalf("coffer file could not be decoded")
		}

		log.Printf("decoded data len: %d", len(decoded))
		return nacl.Decrypt(decoded, key)
	}

	log.Fatalf("unable to decrypt coffer")

	return []byte{}
}
