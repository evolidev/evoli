package reload

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/mitranim/gg"
)

const (
	AsciiEndOfText      = 3  // ^C
	AsciiFileSeparator  = 28 // ^\
	AsciiDeviceControl2 = 18 // ^R
	AsciiDeviceControl4 = 20 // ^T
	AsciiUnitSeparator  = 31 // ^- or ^?

	CodeInterrupt    = AsciiEndOfText
	CodeQuit         = AsciiFileSeparator
	CodeRestart      = AsciiDeviceControl2
	CodeStop         = AsciiDeviceControl4
	CodePrintCommand = AsciiUnitSeparator
)

var (
	FdTerm        = syscall.Stdin
	KillSignals   = []syscall.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	KillSignalsOs = gg.Map(KillSignals, toOsSignal[syscall.Signal])
	ReWord        = regexp.MustCompile(`^\w+$`)
	PathSep       = string([]rune{os.PathSeparator})

	RepSingleMulti = strings.NewReplacer(
		`\r\n`, gg.Newline,
		`\r`, gg.Newline,
		`\n`, gg.Newline,
	).Replace

	RepMultiSingle = strings.NewReplacer(
		"\r\n", `\n`,
		"\r", `\n`,
		"\n", `\n`,
	).Replace
)

/**
Implemented by `notify.EventInfo`.
Path must be an absolute filesystem path.
*/
type FsEvent interface{ Path() string }

// Implemented by `WatchNotify`.
type Watcher interface {
	Init(*Main)
	DeInit()
	Run(*Main)
}

func commaSplit(val string) []string {
	if len(val) == 0 {
		return nil
	}
	return strings.Split(val, `,`)
}

func commaJoin(val []string) string { return strings.Join(val, `,`) }

func cleanExtension(val string) string {
	ext := filepath.Ext(val)
	if len(ext) > 0 && ext[0] == '.' {
		return ext[1:]
	}
	return ext
}

func validateExtension(val string) {
	if !ReWord.MatchString(val) {
		panic(gg.Errf(`invalid extension %q`, val))
	}
}

func toAbsPath(val string) string {
	if !filepath.IsAbs(val) {
		val = filepath.Join(cwd, val)
	}
	return filepath.Clean(val)
}

func toDirPath(val string) string {
	if val == `` || strings.HasSuffix(val, PathSep) {
		return val
	}
	return val + PathSep
}

func toAbsDirPath(val string) string { return toDirPath(toAbsPath(val)) }

func toOsSignal[A os.Signal](src A) os.Signal { return src }

func recLog() {
	val := recover()
	if val != nil {
		log.Println(val)
	}
}

func withNewline[A ~string](val A) A {
	if gg.HasNewlineSuffix(val) {
		return val
	}
	return val + A(gg.Newline)
}
