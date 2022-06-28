package filesystem

import (
	"github.com/evolidev/evoli/framework/use"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Read(path string) string {
	dat, err := os.ReadFile(path)
	use.AbortUnless(err)

	return string(dat)
}

func MakeDirectory(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		use.AbortUnless(err)
	}
}

func Write(path string, data string) {
	MakeDirectory(filepath.Dir(path))

	err := ioutil.WriteFile(path, []byte(data), 0644)
	use.AbortUnless(err)
}

func Delete(path string) {
	err := os.Remove(path)
	use.AbortUnless(err)
}
