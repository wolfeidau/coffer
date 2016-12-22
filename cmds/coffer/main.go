package main

import (
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/wolfeidau/coffer"
)

const S3BucketEnv = "S3_BUCKET"

var (
	app        = kingpin.New("coffer", "A command line tool to manage encrypted coffer files.")
	debug      = app.Flag("debug", "Enable debug mode.").Bool()
	cofferFile = app.Flag("coffer-file", "Coffer file.").Required().String()

	// either we supply a alias
	alias = app.Flag("alias", "KMS key alias.").OverrideDefaultFromEnvar("COFFER_KEY_ALIAS").Default("alias/coffer").String()

	// or we are using KMS
	key = app.Flag("key", "The KMS key-id of the master key to use.").OverrideDefaultFromEnvar("COFFER_KEY").String()

	// commands
	encrypt        = app.Command("encrypt", "Encrypt the coffer file.")
	decrypt        = app.Command("decrypt", "Decrypt the coffer file.")
	sync           = app.Command("sync", "Sync the coffer file with the local filesystem.")
	base           = sync.Flag("base", "Set a base path for testing.").String()
	upload         = app.Command("upload", "Upload a bundle to s3.")
	uploadBucket   = upload.Flag("bucket", "Name of the s3 bucket").String()
	download       = app.Command("download", "Download a bundle from s3.")
	downloadBucket = download.Flag("bucket", "Name of the s3 bucket").String()

	downloadSync       = app.Command("download-sync", "Download a bundle from s3 and sync with the local filesystem.")
	downloadSyncBucket = downloadSync.Flag("bucket", "Name of the s3 bucket").String()
	downloadSyncBase   = downloadSync.Flag("base", "Set a base path for testing.").String()

	// Version app version updated by build script
	Version = ""
)

func main() {
	app.Version(Version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// encrypt the coffer file
	case encrypt.FullCommand():
		log.Printf("encrypt")
		coffer.MustEncrypt(*cofferFile, *alias)
	// decrypt the coffer file
	case decrypt.FullCommand():
		log.Printf("decrypt")
		coffer.MustDecrypt(*cofferFile, *alias)
	// sync the file system to the coffer file
	case sync.FullCommand():
		log.Printf("sync")
		coffer.MustSync(*cofferFile, *alias, *base)
	case upload.FullCommand():
		log.Printf("upload")
		s3Bucket := validateBucketName(*uploadBucket)
		coffer.MustUpload(*cofferFile, *alias, s3Bucket)
	case download.FullCommand():
		log.Printf("download")
		s3Bucket := validateBucketName(*uploadBucket)
		coffer.MustDownload(*cofferFile, *alias, s3Bucket)
	case downloadSync.FullCommand():
		log.Printf("download-sync")
		s3Bucket := validateBucketName(*uploadBucket)
		coffer.MustDownloadSync(*cofferFile, *alias, s3Bucket, *downloadSyncBase)
	}
}

func validateBucketName(bucket string) string {
	if bucket != "" {
		return bucket
	}

	s3Bucket := os.Getenv(S3BucketEnv)

	if s3Bucket == "" {
		log.Fatalf("bucket name is required")
	}

	return s3Bucket
}
