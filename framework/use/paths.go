package use

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var rootDir string

func isBuildAsRun() bool {
	s := os.Args[0]
	p := "go-build\\d+"

	res, _ := regexp.MatchString(p, s)
	fmt.Println("matching", res)
	return res
}

func BasePath(path ...string) string {
	if rootDir != "" {
		return rootDir
	}

	if !isBuildAsRun() {
		rootDir = getByExecutable()
	} else {
		rootDir = getByRuntime()
	}

	output := strings.TrimSpace(rootDir)
	for _, p := range path {
		output = output + "/" + p
	}

	return output
}

func getByExecutable() string {
	filePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)

	}
	rootDir = filepath.Dir(filePath)

	return rootDir
}

func getByRuntime() string {
	cnt := 0

	for true {
		_, b, _, ok := runtime.Caller(cnt)

		if !ok {
			break
		}

		tmp := path.Join(path.Dir(b))
		res1 := strings.HasSuffix(tmp, "src/runtime")

		if res1 {
			cnt--

			_, c, _, _ := runtime.Caller(cnt)
			return path.Join(path.Dir(c))
		}
		cnt++
	}

	return ""
}
