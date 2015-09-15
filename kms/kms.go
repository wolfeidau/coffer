package kms

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/service/kms"
)

// DataKey which contains the details of the KMS key
type DataKey struct {
	CiphertextBlob []byte
	Plaintext      []byte
}

// Decrypt decrypt the cipher text using the provided key.
func Decrypt(ciphertext []byte) (*DataKey, error) {

	// setup the creds chain to configure from environment and ec2 IAM role.
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.SharedCredentialsProvider{},
			&ec2rolecreds.EC2RoleProvider{},
		})

	svc := kms.New(&aws.Config{Credentials: creds})

	params := &kms.DecryptInput{
		CiphertextBlob:    ciphertext,
		EncryptionContext: map[string]*string{},
		GrantTokens:       []*string{},
	}
	resp, err := svc.Decrypt(params)

	if err != nil {
		return nil, err
	}

	return &DataKey{
		CiphertextBlob: ciphertext,
		Plaintext:      resp.Plaintext, // transfer the plain text key after decryption
	}, nil
}

// GenerateDataKey generate a key for use saving a coffer file
func GenerateDataKey(alias string) (*DataKey, error) {

	svc := kms.New(nil)

	params := &kms.GenerateDataKeyInput{
		KeyId:             aws.String(alias),
		EncryptionContext: map[string]*string{},
		GrantTokens:       []*string{},
		NumberOfBytes:     aws.Int64(64),
	}

	resp, err := svc.GenerateDataKey(params)

	if err != nil {
		return nil, err
	}

	return &DataKey{
		CiphertextBlob: resp.CiphertextBlob,
		Plaintext:      resp.Plaintext[:32], // just return 32 bytes for the key
	}, nil
}
