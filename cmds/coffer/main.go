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
	encrypt = app.Command("encrypt", "Encrypt the coffer file.")
	decrypt = app.Command("decrypt", "Decrypt the coffer file.")
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
	}
}
