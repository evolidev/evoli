/**
Go Watch: missing watch mode for the "go" command. Invoked exactly like the
"go" command, but also watches Go files and reruns on changes.
*/
package reload

import (
	e "errors"
	l "log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/mitranim/gg"
)

var (
	log = l.New(os.Stderr, `[gow] `, 0)
	//cwd = gg.Cwd()
)

var cwd string

func Init() {
	cwd = "/Users/omer/Code/evoli/examples"
	log.Println(cwd)
	start()
}

func start() {
	var main Main
	defer main.Exit()
	defer main.Deinit()
	main.Init()
	main.Run()
}

type Main struct {
	Opt         Opt
	Cmd         Cmd
	TermState   TermState
	Watcher     Watcher
	LastChar    byte
	LastInst    time.Time
	ChanSignals gg.Chan[os.Signal]
	ChanRestart gg.Chan[struct{}]
	ChanKill    gg.Chan[syscall.Signal]
}

func (m *Main) Init() {
	m.Opt.Init(os.Args[1:])

	m.ChanRestart.Init()
	m.ChanKill.Init()

	m.Cmd.Init()
	m.StdInInit()
	m.SigInit()
	m.WatchInit()
	m.TermInit()
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
func (m *Main) Deinit() {
	m.TermDeInit()
	m.WatchDeinit()
	m.SigDeinit()
	m.Cmd.Deinit()
}

func (m *Main) Run() {
	go m.StdInRun()
	go m.SigRun()
	go m.WatchRun()
	m.CmdRun()
}

func (m *Main) TermInit() {
	if m.Opt.Raw {
		m.TermState.Init()
	}
}

func (m *Main) TermDeInit() { m.TermState.Deinit() }

func (m *Main) StdInInit() { m.AfterByte(0) }

/**
 * See `Main.InitTerm`. "Raw mode" allows us to support our own control codes,
 * but we're also responsible for interpreting common ASCII codes into OS signals.
 */
func (m *Main) StdInRun() {
	buf := make([]byte, 1, 1)

	for {
		size, err := os.Stdin.Read(buf)
		if err != nil || size == 0 {
			return
		}
		m.OnByte(buf[0])
	}
}

/**
 * Interpret known ASCII codes as OS signals.
 * Otherwise forward the input to the subprocess.
 */
func (m *Main) OnByte(val byte) {
	defer recLog()
	defer m.AfterByte(val)

	switch val {
	case CodeInterrupt:
		m.OnCodeInterrupt()

	case CodeQuit:
		m.OnCodeQuit()

	case CodePrintCommand:
		m.OnCodePrintCommand()

	case CodeRestart:
		m.OnCodeRestart()

	case CodeStop:
		m.OnCodeStop()

	default:
		m.OnByteAny(val)
	}
}

func (m *Main) AfterByte(val byte) {
	m.LastChar = val
	m.LastInst = time.Now()
}

func (m *Main) OnCodeInterrupt() {
	if m.Cmd.Has() {
		if m.LastChar == CodeInterrupt &&
			time.Now().Sub(m.LastInst) < time.Second {
			if m.Opt.Verb {
				log.Println(`received ^C^C, shutting down`)
			}
			m.Kill(syscall.SIGINT)
			return
		}

		if m.Opt.Verb {
			log.Println(`received ^C, stopping subprocess`)
		}
		m.Cmd.Broadcast(syscall.SIGINT)
		return
	}

	if m.Opt.Verb {
		log.Println(`received ^C, shutting down`)
	}
	m.Kill(syscall.SIGINT)
}

func (m *Main) OnCodeQuit() {
	if m.Cmd.Has() {
		if m.Opt.Verb {
			log.Println(`received ^\, stopping subprocess`)
		}
		m.Cmd.Broadcast(syscall.SIGQUIT)
		return
	}

	if m.Opt.Verb {
		log.Println(`received ^\, shutting down`)
	}
	m.Kill(syscall.SIGQUIT)
}

func (m *Main) OnCodePrintCommand() {
	log.Printf(`current command: %q`, os.Args)
}

func (m *Main) OnCodeRestart() {
	if m.Opt.Verb {
		log.Println(`received ^R, restarting`)
	}
	m.Restart()
}

func (m *Main) OnCodeStop() {
	if m.Cmd.Has() {
		if m.Opt.Verb {
			log.Println(`received ^T, stopping`)
		}
		m.Cmd.Broadcast(syscall.SIGTERM)
		return
	}

	if m.Opt.Verb {
		log.Println(`received ^T, nothing to stop`)
	}
}

func (m *Main) OnByteAny(char byte) { m.Cmd.WriteChar(char) }

/**
We override Go's default signal handling to ensure cleanup before exit.
Cleanup is necessary to restore the previous terminal state and kill any
sub-sub-processes.

The set of signals registered here MUST match the set of signals explicitly
handled by this program; see below.
*/
func (m *Main) SigInit() {
	m.ChanSignals.InitCap(1)
	signal.Notify(m.ChanSignals, KillSignalsOs...)
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

		if gg.Has(KillSignals, sig) {
			m.Kill(sig)
			return
		}

		if m.Opt.Verb {
			log.Println(`received signal:`, sig)
		}
	}
}

func (m *Main) WatchInit() {
	wat := new(WatchNotify)
	wat.Init(m)
	m.Watcher = wat
}

func (m *Main) WatchDeinit() {
	if m.Watcher != nil {
		m.Watcher.DeInit()
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
			m.Opt.TermClear()
			continue

		case val := <-m.ChanKill:
			m.Cmd.Broadcast(val)
			m.Deinit()
			gg.Nop1(syscall.Kill(os.Getpid(), val))
			return
		}
	}
}

func (m *Main) CmdWait(cmd *exec.Cmd) {
	err := cmd.Wait()

	if err != nil {
		// `go run` reports the program's exit code to stderr.
		// In this case we suppress the error message to avoid redundancy.
		if !(gg.Head(m.Opt.Args) == `run` && e.As(err, new(*exec.ExitError))) {
			log.Printf(`subcommand error: %v`, err)
		}
	} else if m.Opt.Verb {
		log.Println(`exit ok`)
	}

	m.Opt.Sep.Dump(log.Writer())
}

// Must be deferred.
func (m *Main) Exit() {
	err := gg.AnyErrTraced(recover())
	if err != nil {
		m.Opt.LogErr(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func (m *Main) Restart() {
	m.ChanRestart.SendZeroOpt()
}

func (m *Main) Kill(val syscall.Signal) {
	m.ChanKill.SendOpt(val)
}
