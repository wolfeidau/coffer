package coffer

import (
	"bytes"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3svc = s3.New(nil)

func mustUpload(bucket string, filename string, payload []byte) {

	log.Printf("putting file=%s into bucket=%s length=%d", filename, bucket, len(payload))

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),   // Required
		Key:           aws.String(filename), // Required
		Body:          bytes.NewReader(payload),
		ContentLength: aws.Int64(int64(len(payload))),
	}

	resp, err := s3svc.PutObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// A service error occurred.
			log.Fatalf("AWS Error code: %v message: %s", awsErr.Code(), awsErr.Message())
		}
		// A non-service error occurred.
		log.Fatalf("Error make request to AWS: %v", err)
	}

	// Pretty-print the response data.
	log.Printf("response %vs", resp)

}

func mustDownload(bucket string, filename string) []byte {

	log.Printf("getting file=%s from bucket=%s", filename, bucket)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),   // Required
		Key:    aws.String(filename), // Required
	}

	resp, err := s3svc.GetObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// A service error occurred.
			log.Fatalf("AWS Error code: %v message: %s", awsErr.Code(), awsErr.Message())
		}
		// A non-service error occurred.
		log.Fatalf("Error make request to AWS: %v", err)
	}

	defer resp.Body.Close()

	payload := new(bytes.Buffer)

	io.Copy(payload, resp.Body)

	// Pretty-print the response data.
	log.Printf("response %v", resp)

	return payload.Bytes()
}
