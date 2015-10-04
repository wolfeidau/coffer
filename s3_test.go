package coffer

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bmizerany/assert"
)

func TestUpload(t *testing.T) {
	fakeS3 := &FakeS3{
		PutOutputs: []s3.PutObjectOutput{
			{},
			{},
		},
	}

	s3Svc = fakeS3

	mustUpload("mybucket", "test.coffer", make([]byte, 32))
}

func TestDownload(t *testing.T) {

	ciphertext := []byte{0xa1, 0xb2, 0xc3, 0xd4}

	fakeS3 := &FakeS3{
		GetOutputs: []s3.GetObjectOutput{
			{
				Body: ioutil.NopCloser(bytes.NewReader(ciphertext)),
			},
		},
	}

	s3Svc = fakeS3

	payload := mustDownload("mybucket", "test.coffer")

	assert.Equal(t, []byte{0xa1, 0xb2, 0xc3, 0xd4}, payload)
}
