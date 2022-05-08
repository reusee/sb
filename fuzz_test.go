package sb

import (
	"io/ioutil"
	"path/filepath"
	"strings"
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
	files, err = filepath.Glob("crashers/*")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file, ".output") {
			continue
		}
		if strings.HasSuffix(file, ".quoted") {
			continue
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		Fuzz(data)
	}
}

func FuzzAllCorpus(f *testing.F) {

	files, err := filepath.Glob("corpus/*")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		f.Add(data)
	}

	files, err = filepath.Glob("crashers/*")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file, ".output") {
			continue
		}
		if strings.HasSuffix(file, ".quoted") {
			continue
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		f.Add(data)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		switch Fuzz(data) {
		case 0:
			t.Skip()
		}
	})

}
