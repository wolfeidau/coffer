package coffer

import (
	"log"

	"gopkg.in/yaml.v2"
)

// Coffer used as the container for an encrypted coffer
type Coffer struct {
	Name       string `yaml:"name,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Key        string `yaml:"key,omitempty"`
	CipherText string `yaml:"ct,omitempty"`
}

// Validate checks the validity of the coffer
func (c *Coffer) Validate() bool {

	if len(c.Version) == 0 {
		return false
	}

	if len(c.CipherText) == 0 {
		return false
	}

	return true
}

// DecodeCoffer decode the coffer file
func DecodeCoffer(data []byte) (coffer *Coffer, err error) {

	err = yaml.Unmarshal(data, &coffer)

	if err != nil {
		return
	}

	return
}

func mustEncodeCoffer(coffer *Coffer) []byte {

	data, err := yaml.Marshal(coffer)
	if err != nil {
		log.Fatalf("Unable to encode YAML data: %v", err)
	}

	return data
}
