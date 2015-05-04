package main

import (
	"os"

	"github.com/wolfeidau/coffer"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	app        = kingpin.New("coffer", "A command line tool to manage encrypted coffer files.")
	debug      = app.Flag("debug", "Enable debug mode.").Bool()
	cofferFile = app.Flag("coffer-file", "Coffer file.").Required().String()
	secret     = app.Flag("secret", "Coffer secret.").OverrideDefaultFromEnvar("COFFER_SECRET").Required().String()

	// commands
	encrypt        = app.Command("encrypt", "Encrypt the coffer file.")
	decrypt        = app.Command("decrypt", "Decrypt the coffer file.")
	sync           = app.Command("sync", "Sync the coffer file with the local filesystem.")
	base           = sync.Flag("base", "Set a base path for testing.").String()
	upload         = app.Command("upload", "Upload a bundle to s3.")
	uploadBucket   = upload.Flag("bucket", "Name of the s3 bucket").Required().String()
	download       = app.Command("download", "Download a bundle from s3.")
	downloadBucket = download.Flag("bucket", "Name of the s3 bucket").Required().String()
)

func main() {
	kingpin.Version("0.0.1")

	if *debug {
		coffer.Verbose = true
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// encrypt the coffer file
	case encrypt.FullCommand():
		coffer.Debugf("encrypt")
		coffer.MustEncrypt(*cofferFile, *secret)
	// decrypt the coffer file
	case decrypt.FullCommand():
		coffer.Debugf("decrypt")
		coffer.MustDecrypt(*cofferFile, *secret)
	// sync the file system to the coffer file
	case sync.FullCommand():
		coffer.Debugf("sync")
		coffer.MustSync(*cofferFile, *secret, *base)
	case upload.FullCommand():
		coffer.Debugf("upload")
		coffer.MustUpload(*cofferFile, *secret, *uploadBucket)
	case download.FullCommand():
		coffer.Debugf("download")
		coffer.MustDownload(*cofferFile, *secret, *downloadBucket)
	}
}
