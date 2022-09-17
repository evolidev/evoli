package use

import (
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

	return res
}

func StoragePath(path ...string) string {
	output := ""
	for _, p := range path {
		output += "/" + p
	}

	return BasePath("storage/" + output)
}

func BasePath(path ...string) string {
	if rootDir == "" {
		if !isBuildAsRun() {
			rootDir = getByExecutable()
		} else {
			rootDir = getByRuntime()
		}

		rootDir = strings.TrimSpace(rootDir + "/")
	}

	output := rootDir
	for _, p := range path {
		if p == "" {
			continue
		}
		output = output + "/" + p
	}

	return strings.TrimSpace(output)
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

		if cnt > 50 {
			break
		}
	}

	return ""
}
