package coffer

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/bmizerany/assert"
	"github.com/kr/pretty"
)

func TestGenerateDataKey(t *testing.T) {

	fakeKMS := &FakeKMS{
		GenerateOutputs: []kms.GenerateDataKeyOutput{
			{
				CiphertextBlob: []byte("yay"),
				KeyId:          aws.String("coffer"),
				Plaintext:      make([]byte, 32),
			},
		},
		DecryptOutputs: []kms.DecryptOutput{
			{
				KeyId:     aws.String("key1"),
				Plaintext: make([]byte, 32),
			},
		},
	}

	kmsSvc = fakeKMS

	dataKey := mustGenerateDataKey("coffer")

	pretty.Printf("dataKey %+v\n", dataKey)

	assert.Equal(t, []byte("yay"), dataKey.CiphertextBlob)
	assert.Equal(t, make([]byte, 32), dataKey.Plaintext)
}

func TestEncrypt(t *testing.T) {

	fakeKMS := &FakeKMS{
		GenerateOutputs: []kms.GenerateDataKeyOutput{
			{
				CiphertextBlob: []byte("yay"),
				KeyId:          aws.String("coffer"),
				Plaintext:      make([]byte, 32),
			},
		},
		DecryptOutputs: []kms.DecryptOutput{
			{
				KeyId:     aws.String("coffer"),
				Plaintext: make([]byte, 32),
			},
		},
	}

	kmsSvc = fakeKMS

	ciphertext := make([]byte, 24)

	dataKey := mustDecrypt(ciphertext)

	assert.Equal(t, make([]byte, 32), dataKey.Plaintext)

}
