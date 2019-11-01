package sb

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestFuzzCorpus(t *testing.T) {
	files, err := filepath.Glob("corpus/*")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		Fuzz(data)
	}
}
