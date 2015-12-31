package coffer

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

func mustExtractBundle(bundle *Bundle, base string) {

	log.Printf("extracting bundle to %s", base)

	// if the base path is set then use this as the base for all creation
	if base != "" {
		for k, f := range bundle.Files {
			rpath := path.Join(base, k)

			mustMkDirAll(path.Dir(rpath))
			mustWriteFile(rpath, []byte(f.Content), os.FileMode(f.Mode))
		}
	}
}

func mustReadFile(path string) []byte {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	return data
}

func mustWriteFile(path string, data []byte, mode os.FileMode) []byte {

	log.Printf("WriteFile %s", path)

	err := ioutil.WriteFile(path, data, mode)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}

	// ensure the mode is set even when we haven't modified the file.
	os.Chmod(path, mode)

	return data
}

func mustMkDirAll(rpath string) {
	log.Printf("MkDirAll %s", rpath)

	err := os.MkdirAll(rpath, 0755)
	if err != nil {
		log.Fatalf("failed to create any necessary parents: %v", err)
	}
}
