package coffer

import (
	"bytes"
	"io"
	"net/url"

	s3 "github.com/rlmcpherson/s3gof3r"
)

func mustUpload(bucket string, filename string, payload []byte) {

	u, err := url.Parse(bucket)
	if err != nil {
		Fatalf("Failed to parse bucket URL: %v", err)
	}

	k, err := s3.EnvKeys() // get S3 keys from environment
	if err != nil {
		Fatalf("Unable to load AWS credentials from environment")
	}

	// Open bucket to put file into
	s3 := s3.New("s3.amazonaws.com", k)
	b := s3.Bucket(u.Host)

	// Open a PutWriter for upload
	w, err := b.PutWriter(filename, nil, nil)
	if err != nil {
		Fatalf("Put to bucket failed: %v", err)
	}

	if _, err = io.Copy(w, bytes.NewBuffer(payload)); err != nil { // Copy into S3
		Fatalf("Copy content to bucket failed: %v", err)
	}

	if err = w.Close(); err != nil {
		Fatalf("Closing file failed: %v", err)
	}
}
