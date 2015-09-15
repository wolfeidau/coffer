package main

import (
	"log"
	"os"

	"github.com/wolfeidau/coffer"
	"gopkg.in/alecthomas/kingpin.v2"
)

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
	uploadBucket   = upload.Flag("bucket", "Name of the s3 bucket").OverrideDefaultFromEnvar("S3_BUCKET").Required().String()
	download       = app.Command("download", "Download a bundle from s3.")
	downloadBucket = download.Flag("bucket", "Name of the s3 bucket").OverrideDefaultFromEnvar("S3_BUCKET").Required().String()

	downloadSync       = app.Command("download-sync", "Download a bundle from s3 and sync with the local filesystem.")
	downloadSyncBucket = downloadSync.Flag("bucket", "Name of the s3 bucket").OverrideDefaultFromEnvar("S3_BUCKET").Required().String()
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
		coffer.MustUpload(*cofferFile, *alias, *uploadBucket)
	case download.FullCommand():
		log.Printf("download")
		coffer.MustDownload(*cofferFile, *alias, *downloadBucket)
	case downloadSync.FullCommand():
		log.Printf("download-sync")
		coffer.MustDownloadSync(*cofferFile, *alias, *downloadSyncBucket, *downloadSyncBase)
	}
}
