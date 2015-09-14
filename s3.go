package coffer

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

func mustUpload(bucket string, filename string, payload []byte) {

	Infof("putting file=%s into bucket=%s length=%d", filename, bucket, len(payload))

	svc := newS3Svc()

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),   // Required
		Key:           aws.String(filename), // Required
		Body:          bytes.NewReader(payload),
		ContentLength: aws.Int64(int64(len(payload))),
	}

	resp, err := svc.PutObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// A service error occurred.
			Fatalf("AWS Error code: %v message: %s", awsErr.Code(), awsErr.Message())
		}
		// A non-service error occurred.
		Fatalf("Error make request to AWS: %v", err)
	}

	// Pretty-print the response data.
	Infof("response %vs", resp)

}

func mustDownload(bucket string, filename string) []byte {

	Infof("getting file=%s from bucket=%s", filename, bucket)

	svc := newS3Svc()

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),   // Required
		Key:    aws.String(filename), // Required
	}

	resp, err := svc.GetObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// A service error occurred.
			Fatalf("AWS Error code: %v message: %s", awsErr.Code(), awsErr.Message())
		}
		// A non-service error occurred.
		Fatalf("Error make request to AWS: %v", err)
	}

	defer resp.Body.Close()

	payload := new(bytes.Buffer)

	io.Copy(payload, resp.Body)

	// Pretty-print the response data.
	Infof("response %v", resp)

	return payload.Bytes()
}

func newS3Svc() *s3.S3 {

	// setup the creds chain to configure from environment and ec2 IAM role.
	creds := credentials.NewEnvCredentials()

	return s3.New(&aws.Config{Credentials: creds})
}
