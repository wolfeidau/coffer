package coffer

import "gopkg.in/yaml.v2"

type Bundle struct {
	Files map[string]*FileData `yaml:"files"`
}

func (b *Bundle) MustValidate() {

	for k, f := range b.Files {
		if k == "" {
			Fatalf("Validation failed: file name must be set for %v", f)
		}

		f.MustValidate(k)
	}
}

type FileData struct {
	Mode    uint32 `yaml:"mode"`
	Owner   string `yaml:"owner"`
	Group   string `yaml:"group"`
	Content string `yaml:"content"`
}

func (f *FileData) MustValidate(name string) {
	if f.Mode == 0 {
		Fatalf("Validation failed: file mode must be set for %s", name)
	}
}

func mustDecodeBundle(data []byte) *Bundle {

	var bundle *Bundle

	err := yaml.Unmarshal(data, &bundle)

	if err != nil {
		Fatalf("Unable to decode YAML data: %v", err)
	}

	return bundle
}

func mustEncodeBundle(bundle *Bundle) []byte {

	data, err := yaml.Marshal(bundle)
	if err != nil {
		Fatalf("Unable to encode YAML data: %v", err)
	}

	return data
}
