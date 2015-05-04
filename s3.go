package coffer

import (
	"bytes"
	"io"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/s3"
)

func mustUpload(bucket string, filename string, payload []byte) {

	Infof("putting file=%s into bucket=%s length=%d", filename, bucket, len(payload))

	svc := s3.New(nil)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),   // Required
		Key:           aws.String(filename), // Required
		Body:          bytes.NewReader(payload),
		ContentLength: aws.Long(int64(len(payload))),
	}

	resp, err := svc.PutObject(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		Fatalf("AWS Error code: %v message: %s", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		Fatalf("Error make request to AWS: %v", err)
	}

	// Pretty-print the response data.
	Infof("response %s", awsutil.StringValue(resp))

}

func mustDownload(bucket string, filename string) []byte {

	Infof("getting file=%s from bucket=%s", filename, bucket)

	svc := s3.New(nil)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),   // Required
		Key:    aws.String(filename), // Required
	}

	resp, err := svc.GetObject(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		Fatalf("AWS Error code: %v message: %s", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		Fatalf("Error make request to AWS: %v", err)
	}

	defer resp.Body.Close()

	payload := new(bytes.Buffer)

	io.Copy(payload, resp.Body)

	// Pretty-print the response data.
	Infof("response %s", awsutil.StringValue(resp))

	return payload.Bytes()
}
