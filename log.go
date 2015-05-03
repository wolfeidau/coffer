package coffer

import (
	"fmt"
	"log"
	"os"
)

var (
	// Quiet suppresses all logging output below ERROR
	Quiet = false

	// Verbose enables logging output below INFO
	Verbose = false

	// Logger is the log.Logger object that backs this logger.
	Logger = log.New(os.Stdout, "", log.LstdFlags)
)

var Fatalf = log.Fatalf

func fatalf(format string, args ...interface{}) {
	Logger.Printf("FATAL "+format, args...)
	fmt.Scanln()
	os.Exit(1)
}

func Infof(format string, args ...interface{}) {
	if !Quiet {
		Logger.Printf("INFO "+format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if Verbose && !Quiet {
		Logger.Printf("DEBUG "+format, args...)
	}
}
