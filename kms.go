package coffer

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// KeyManagement is a sub-set of the capabilities of the KMS client.
type KeyManagement interface {
	GenerateDataKey(*kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error)
	Decrypt(*kms.DecryptInput) (*kms.DecryptOutput, error)
}

var kmsSvc KeyManagement

func init() {
	kmsSvc = kms.New(nil)
}

// DataKey which contains the details of the KMS key
type DataKey struct {
	CiphertextBlob []byte
	Plaintext      []byte
}

func mustDecrypt(ciphertext []byte) *DataKey {

	params := &kms.DecryptInput{
		CiphertextBlob:    ciphertext,
		EncryptionContext: map[string]*string{},
		GrantTokens:       []*string{},
	}
	resp, err := kmsSvc.Decrypt(params)

	if err != nil {
		log.Fatalf("kms error: %s", err.Error())
	}

	return &DataKey{
		CiphertextBlob: ciphertext,
		Plaintext:      resp.Plaintext, // transfer the plain text key after decryption
	}
}

func mustGenerateDataKey(alias string) *DataKey {

	params := &kms.GenerateDataKeyInput{
		KeyId:             aws.String(alias),
		EncryptionContext: map[string]*string{},
		GrantTokens:       []*string{},
		NumberOfBytes:     aws.Int64(64),
	}

	resp, err := kmsSvc.GenerateDataKey(params)

	if err != nil {
		log.Fatalf("kms error: %s", err.Error())
	}

	return &DataKey{
		CiphertextBlob: resp.CiphertextBlob,
		Plaintext:      resp.Plaintext[:32], // just return 32 bytes for the key
	}
}
