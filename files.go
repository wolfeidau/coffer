package coffer

import (
	"io/ioutil"
	"os"
	"path"
)

func mustExtractBundle(bundle *Bundle, base string) {

	// if the base path is set then use this as the base for all creation
	if base != "" {
		for k, f := range bundle.Files {
			rpath := path.Join(base, k)

			mustMkDirAll(path.Base(rpath))
			mustWriteFile(rpath, []byte(f.Content), os.FileMode(f.Mode))
		}
	}
}

func mustReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		Fatalf("Unable to open file: %v", err)
	}
	return data
}

func mustWriteFile(path string, data []byte, mode os.FileMode) []byte {
	Infof("WriteFile %s", path)
	err := ioutil.WriteFile(path, data, mode)

	if err != nil {
		Fatalf("Unable to open file: %v", err)
	}
	return data
}

func mustMkDirAll(rpath string) {
	Infof("MkDirAll %s", path.Base(rpath))
	err := os.MkdirAll(path.Base(rpath), 755)
	if err != nil {
		Fatalf("failed to create any necessary parents: %v", err)
	}
}