package coffer

import (
	"encoding/base64"
	"log"
	"os"
	"path"

	"github.com/wolfeidau/coffer/kms"
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
func MustEncrypt(cofferFile, alias string) {

	data := mustReadFile(cofferFile)

	payload := mustEncryptPayload(data, alias)

	mustWriteFile(cofferFile, payload, OwnerRead)
}

// MustDecrypt decrypt the supplied file
func MustDecrypt(cofferFile, alias string) []byte {

	data := mustReadFile(cofferFile)

	payload := mustDecryptPayload(data, alias)

	return mustWriteFile(cofferFile, payload, OwnerRead)

}

// MustSync sync the file up to S3
func MustSync(cofferFile, alias, base string) {

	// base not set
	if base == "" {
		base = "/" // set the base directory to '/'
	}

	payload := mustReadFile(cofferFile)

	// if the coffer file is encrypted, decrypt it
	if isEncrypted(payload) {
		payload = mustDecryptPayload(payload, alias)
	}

	bundle := mustDecodeBundle(payload)

	bundle.MustValidate()

	mustExtractBundle(bundle, base)
}

// MustUpload upload the file to the supplied s3 bucket
func MustUpload(cofferFile, alias, bucket string) {

	payload := mustReadFile(cofferFile)

	// if the coffer is not encrypted
	if !isEncrypted(payload) {
		payload = mustEncryptPayload(payload, alias)
	}

	filename := path.Base(cofferFile)

	mustUpload(bucket, filename, payload)
}

// MustDownload download the file from the supplied s3 bucket
func MustDownload(cofferFile, alias, bucket string) {

	filename := path.Base(cofferFile)

	payload := mustDownload(bucket, filename)

	mustWriteFile(cofferFile, payload, OwnerRead)
}

// MustDownloadSync download the file from the supplied s3 bucket and sync
// it to the filesystem
func MustDownloadSync(cofferFile, alias, bucket, base string) {

	// base not set
	if base == "" {
		base = "/" // set the base directory to '/'
	}

	filename := path.Base(cofferFile)

	payload := mustDownload(bucket, filename)

	// if the coffer file is encrypted, decrypt it
	if isEncrypted(payload) {
		payload = mustDecryptPayload(payload, alias)
	}

	bundle := mustDecodeBundle(payload)

	bundle.MustValidate()

	mustExtractBundle(bundle, base)
}

func isEncrypted(data []byte) bool {

	coffer, err := DecodeCoffer(data)

	if err != nil {
		return false
	}

	return coffer.Validate()
}

func mustEncryptPayload(data []byte, alias string) []byte {

	if isEncrypted(data) {
		log.Fatalf("coffer file already encrypted")
	}

	key, err := kms.GenerateDataKey(alias)

	if err != nil {
		log.Fatalf("kms error: %s", err.Error())
	}

	payload := nacl.Encrypt(data, key.Plaintext)
	cipherText := base64.StdEncoding.EncodeToString(payload)
	keyCipherText := base64.StdEncoding.EncodeToString(key.CiphertextBlob)

	log.Printf("encoded data len: %d", len(cipherText))

	return mustEncodeCoffer(&Coffer{
		Key:        keyCipherText,
		Version:    Version,
		CipherText: cipherText,
	})
}

func mustDecryptPayload(data []byte, alias string) []byte {

	coffer, err := DecodeCoffer(data)

	if err != nil {
		log.Fatalf("unable to decode coffer")
	}

	if coffer.Validate() {

		keyData, err := base64.StdEncoding.DecodeString(coffer.Key)

		if err != nil {
			log.Fatalf("unable to decode key")
		}

		key, err := kms.Decrypt(keyData)

		if err != nil {
			log.Fatalf("kms error: %s", err.Error())
		}

		decoded, err := base64.StdEncoding.DecodeString(coffer.CipherText)

		if err != nil {
			log.Fatalf("coffer file could not be decoded")
		}

		log.Printf("decoded data len: %d", len(decoded))
		return nacl.Decrypt(decoded, key.Plaintext)
	}

	log.Fatalf("invalid coffer file, unable to decrypt")

	return []byte{}
}
