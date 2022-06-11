package reload2

import (
	"errors"
	"github.com/evolidev/evoli/framework/use"
	"github.com/mitranim/gg"
	"github.com/rjeczalik/notify"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var cwd = gg.Cwd()

var (
	FD_TERM      = syscall.Stdin
	KILL_SIGS    = []syscall.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	KILL_SIGS_OS = gg.Map(KILL_SIGS, toOsSignal[syscall.Signal])
	RE_WORD      = regexp.MustCompile(`^\w+$`)
	PATH_SEP     = string([]rune{os.PathSeparator})

	REP_SINGLE_MULTI = strings.NewReplacer(
		`\r\n`, gg.Newline,
		`\r`, gg.Newline,
		`\n`, gg.Newline,
	).Replace

	REP_MULTI_SINGLE = strings.NewReplacer(
		"\r\n", `\n`,
		"\r", `\n`,
		"\n", `\n`,
	).Replace
)

type Main struct {
	Cmd         Cmd
	TermState   TermState
	Watcher     *WatchNotify
	LastChar    byte
	LastInst    time.Time
	ChanSignals gg.Chan[os.Signal]
	ChanRestart gg.Chan[struct{}]
	ChanKill    gg.Chan[syscall.Signal]
}

type Watcher interface {
	Init(*Main)
	Deinit()
	Run(*Main)
}

func Watch() {
	//var watcher WatchNotify
	//
	//watcher.Init()
	//watcher.Run()
	var main Main
	main.Init()
}

func (m *Main) Init() {
	//m.Opt.Init(os.Args[1:])

	m.ChanRestart.Init()
	m.ChanKill.Init()

	m.Cmd.Init()
	//m.StdInInit()
	m.SigInit()
	m.WatchInit()
	m.TermInit()
}

func (m *Main) TermInit() {
	if true {
		m.TermState.Init()
	}
}

func (m *Main) TermDeinit() {
	m.TermState.DeInit()
}

/**
We override Go's default signal handling to ensure cleanup before exit.
Cleanup is necessary to restore the previous terminal state and kill any
sub-sub-processes.

The set of signals registered here MUST match the set of signals explicitly
handled by this program; see below.
*/
func (m *Main) SigInit() {
	m.ChanSignals.InitCap(1)
	signal.Notify(m.ChanSignals, KILL_SIGS_OS...)
}

func (m *Main) SigDeinit() {
	if m.ChanSignals != nil {
		signal.Stop(m.ChanSignals)
	}
}

func (m *Main) SigRun() {
	for val := range m.ChanSignals {
		// Should work on all Unix systems. At the time of writing,
		// we're not prepared to support other systems.
		sig := val.(syscall.Signal)

		if gg.Has(KILL_SIGS, sig) {
			m.Kill(sig)
			return
		}

		log.Println(`received signal:`, sig)
	}
}

func (m *Main) WatchInit() {
	wat := new(WatchNotify)
	wat.Init(m)
	m.Watcher = wat
}

func (m *Main) WatchDeinit() {
	if m.Watcher != nil {
		m.Watcher.Deinit()
		m.Watcher = nil
	}
}

func (m *Main) WatchRun() {
	if m.Watcher != nil {
		m.Watcher.Run(m)
	}
}

func (m *Main) CmdRun() {
	for {
		m.Cmd.Restart(m)

		select {
		case <-m.ChanRestart:
			gg.TermClearSoft()
			continue

		case val := <-m.ChanKill:
			m.Cmd.Broadcast(val)
			m.DeInit()
			gg.Nop1(syscall.Kill(os.Getpid(), val))
			return
		}
	}
}

/**
We MUST call this before exiting because:
	* We modify global OS state: terminal, subprocs.
	* OS will NOT auto-cleanup after us.

Otherwise:
	* Terminal is left in unusable state.
	* Subprocs become orphan daemons.

We MUST call this manually before using `syscall.Kill` or `syscall.Exit` on the
current process. Syscalls terminate the process bypassing Go `defer`.
*/
func (m *Main) DeInit() {
	m.TermDeinit()
	m.WatchDeinit()
	m.SigDeinit()
	m.Cmd.Deinit()
}

func (m *Main) CmdWait(cmd *exec.Cmd) {
	err := cmd.Wait()

	if err != nil {
		// `go run` reports the program's exit code to stderr.
		// In this case we suppress the error message to avoid redundancy.
		if !(errors.As(err, new(*exec.ExitError))) {
			log.Printf(`subcommand error: %v`, err)
		}
	} else {
		log.Println(`exit ok`)
	}

	//m.Opt.Sep.Dump(log.Writer())
}

// Must be deferred.
func (m *Main) Exit() {
	err := gg.AnyErrTraced(recover())
	if err != nil {
		log.Println(err)
		//m.Opt.LogErr(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func (m *Main) Restart() { m.ChanRestart.SendZeroOpt() }

func (m *Main) Kill(val syscall.Signal) { m.ChanKill.SendOpt(val) }

type WatchNotify struct {
	Done         gg.Chan[struct{}]
	Events       gg.Chan[notify.EventInfo]
	IgnoredPaths FlagIgnoredPaths `flag:"-i"               desc:"Ignored paths, relative to CWD; multi."`
	Extensions   FlagExtensions   `flag:"-e" init:"go,mod" desc:"Extensions to watch; multi."`
}

func (s *WatchNotify) Init(main *Main) {
	s.Done.Init()
	s.Events.InitCap(1)

	paths := []string{"/Users/omer/Code/evoli/examples"}
	verbose := true

	for _, path := range paths {
		path = filepath.Join(path, `...`)
		if verbose {
			log.Printf(`watching %q`, path)
		}
		gg.Try(notify.Watch(path, s.Events, notify.All))
	}
}

func (s *WatchNotify) Deinit() {
	s.Done.SendZeroOpt()
	if s.Events != nil {
		notify.Stop(s.Events)
	}
}

func (s WatchNotify) Run(main *Main) {
	for {
		select {
		case <-s.Done:
			return

		case event := <-s.Events:
			log.Println(event)
			if s.ShouldRestart(event) {
				log.Println(`restarting on FS event:`, event)
				main.Restart()
			}
		}
	}
}

type FlagExtensions []string

func (s *FlagExtensions) Parse(src string) (err error) {
	defer gg.Rec(&err)
	values := commaSplit(src)
	gg.Each(values, validateExtension)
	gg.AppendVals(s, values...)
	return
}

func (s FlagExtensions) Allow(path string) bool {
	return gg.IsEmpty(s) || gg.Has(s, cleanExtension(path))
}

type FlagWatch []string

func (s *FlagWatch) Parse(src string) error {
	gg.AppendVals(s, commaSplit(src)...)
	return nil
}

type FlagIgnoredPaths []string

func (s *FlagIgnoredPaths) Parse(src string) error {
	values := FlagIgnoredPaths(commaSplit(src))
	values.Norm()
	//gg.AppendVals(s, values...)
	return nil
}

func (s FlagIgnoredPaths) Norm() {
	gg.MapMut(s, toAbsDirPath)
}

func (s FlagIgnoredPaths) Allow(path string) bool {
	return !s.Ignore(path)
}

// Ignore Assumes that the input is an absolute path.
func (s FlagIgnoredPaths) Ignore(path string) bool {
	return gg.Some(s, func(val string) bool {
		return strings.HasPrefix(path, val)
	})
}

func (s *WatchNotify) ShouldRestart(event notify.EventInfo) bool {
	if event == nil {
		return false
	}
	path := event.Path()
	return s.IgnoredPaths.Allow(path) && s.Extensions.Allow(path)
}

func commaSplit(val string) []string {
	if len(val) == 0 {
		return nil
	}
	return strings.Split(val, `,`)
}

func commaJoin(val []string) string {
	return strings.Join(val, `,`)
}

func cleanExtension(val string) string {
	ext := filepath.Ext(val)
	if len(ext) > 0 && ext[0] == '.' {
		return ext[1:]
	}
	return ext
}

func validateExtension(val string) {
	if !RE_WORD.MatchString(val) {
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
	if val == `` || strings.HasSuffix(val, PATH_SEP) {
		return val
	}
	return val + PATH_SEP
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

/**
By default, any regular terminal uses what's known as "cooked mode". It buffers
lines before sending them to the foreground process, and interprets some ASCII
control codes on stdin by sending the corresponding OS signals to the process.
We switch it into "raw mode", where it immediately forwards inputs to our
process's stdin, and doesn't interpret special ASCII codes. This allows to
support special key combinations such as ^R for restarting a subprocess.

The terminal state is shared between all super- and sub-processes. Changes
persist even after our process terminates. We endeavor to restore the previous
state before exiting.

References:
- https://en.wikibooks.org/wiki/Serial_Programming/termios
- man termios
*/
type TermState struct{ gg.Opt[unix.Termios] }

func (s *TermState) Init() {
	s.DeInit()

	state, err := unix.IoctlGetTermios(FD_TERM, use.IoctlReadTermIos)
	if err != nil {
		log.Printf(`unable to read terminal state: %v`, err)
		return
	}
	prev := *state

	// Don't buffer lines.
	state.Lflag &^= unix.ICANON

	// Don't echo characters or special codes.
	state.Lflag &^= unix.ECHO

	// No signals.
	state.Lflag &^= unix.ISIG

	err = unix.IoctlSetTermios(FD_TERM, use.IoctlReadTermIos, state)
	if err != nil {
		log.Printf(`unable to switch terminal to raw mode: %v`, err)
		return
	}

	s.Set(prev)
}

func (s *TermState) DeInit() {
	if s.IsNull() {
		return
	}

	defer s.Clear()
	gg.Nop1(unix.IoctlSetTermios(FD_TERM, use.IoctlWriteTermIos, &s.Val))
}
