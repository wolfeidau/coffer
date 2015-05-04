package coffer

import (
	"bytes"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"github.com/davecgh/go-spew/spew"
)

func mustUpload(bucket string, filename string, payload []byte) {

	Infof("putting file=%s into bucket=%s length=%d", filename, bucket, len(payload))

	config := aws.DefaultConfig
	creds, err := config.Credentials.Credentials()
	if err != nil {
		Fatalf("Unable to load creds: %v", err)
	}

	svc := s3.New(config)

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

	spew.Dump(resp)
	spew.Dump(config)
	spew.Dump(creds)

	// Pretty-print the response data.
	Infof("response %s", awsutil.StringValue(resp))

}
