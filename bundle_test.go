package coffer

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/davecgh/go-spew/spew"
)

var bundleText = []byte(`files:
  /home/user/myfile2:
    mode: 755
    owner: root
    group: root
    content: |
      # this is my file
      # with content
`)

var bundleExpected = &Bundle{
	Files: map[string]*FileData{
		"/home/user/myfile2": &FileData{
			Mode:    755,
			Owner:   "root",
			Group:   "root",
			Content: "# this is my file\n# with content\n",
		},
	},
}

func TestEncodeBundle(t *testing.T) {

	data := mustEncodeBundle(bundleExpected)

	spew.Printf("bundle %s\n", string(data))

	assert.Equal(t, string(bundleText), string(data))
}

func TestDecodeBundle(t *testing.T) {

	bundle := mustDecodeBundle(bundleText)

	spew.Printf("bundle %+v\n", bundle)

	bundle.MustValidate()

	assert.Equal(t, bundle, bundleExpected)
}
