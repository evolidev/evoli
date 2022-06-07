package filesystem

import (
	"github.com/evolidev/evoli/framework/use"
	"os"
)

func Read(path string) string {
	dat, err := os.ReadFile(path)
	use.AbortUnless(err)

	return string(dat)
}
