package coffer

import (
	"log"

	"gopkg.in/yaml.v2"
)

// Bundle bundle of files and their related information
type Bundle struct {
	Files map[string]*FileData `yaml:"files"`
}

// MustValidate checks the validity of the bundle
func (b *Bundle) MustValidate() {

	for k, f := range b.Files {
		if k == "" {
			log.Fatalf("Validation failed: file name must be set for %v", f)
		}

		f.MustValidate(k)
	}
}

// FileData an encoded file with it's permissions
type FileData struct {
	Mode    uint32 `yaml:"mode"`
	Owner   string `yaml:"owner"`
	Group   string `yaml:"group"`
	Content string `yaml:"content"`
}

// FileData checks the validity of the file data structure.
func (f *FileData) MustValidate(name string) {
	if f.Mode == 0 {
		log.Fatalf("Validation failed: file mode must be set for %s", name)
	}
}

func mustDecodeBundle(data []byte) *Bundle {

	var bundle *Bundle

	err := yaml.Unmarshal(data, &bundle)

	if err != nil {
		log.Fatalf("Unable to decode YAML data: %v", err)
	}

	return bundle
}

func mustEncodeBundle(bundle *Bundle) []byte {

	data, err := yaml.Marshal(bundle)
	if err != nil {
		log.Fatalf("Unable to encode YAML data: %v", err)
	}

	return data
}
