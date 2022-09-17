package filesystem

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Read(path string) string {
	dat, _ := os.ReadFile(path)
	//use.AbortUnless(err)

	return string(dat)
}

func MakeDirectory(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		//use.AbortUnless(err)
	}
}

func Write(path string, data string) {
	MakeDirectory(filepath.Dir(path))

	ioutil.WriteFile(path, []byte(data), 0644)
	//use.AbortUnless(err)
}

func Delete(path string) {
	os.Remove(path)
	//use.AbortUnless(err)
}

// Copy copies a file and returns the bytes transferred
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
