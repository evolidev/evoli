package reload2

import (
	"errors"
	"github.com/mitranim/gg"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type Cmd struct {
	sync.Mutex
	Command string
	Args    []string
	Buf     [1]byte
	Cmd     *exec.Cmd
	Stdin   io.WriteCloser
}

func (self *Cmd) Init() {}

func (self *Cmd) Deinit() {
	defer gg.Lock(self).Unlock()
	self.DeinitUnsync()
}

func (self *Cmd) DeinitUnsync() {
	self.BroadcastUnsync(syscall.SIGTERM)
	self.Cmd = nil
	self.Stdin = nil
}

func (s *Cmd) Restart(main *Main) {
	defer gg.Lock(s).Unlock()

	s.DeinitUnsync()

	cmd := s.MakeCmd()
	stdIn, err := cmd.StdinPipe()
	if err != nil {
		log.Printf(`unable to initialize subcommand stdin: %v`, err)
		return
	}

	// Starting the subprocess populates its `.Process`,
	// which allows us to kill the subprocess group on demand.
	err = cmd.Start()
	if err != nil {
		log.Printf(`unable to start subcommand: %v`, err)
		return
	}

	s.Cmd = cmd
	s.Stdin = stdIn
	//go CmdWait(cmd)
	main.CmdWait(cmd)
}

func CmdWait(cmd *exec.Cmd) {
	err := cmd.Wait()
	log.Println(err)

	//if err != nil {
	//	// `go run` reports the program's exit code to stderr.
	//	// In this case we suppress the error message to avoid redundancy.
	//	if !(gg.Head(self.Opt.Args) == `run` && e.As(err, new(*exec.ExitError))) {
	//		log.Printf(`subcommand error: %v`, err)
	//	}
	//} else if self.Opt.Verb {
	//	log.Println(`exit ok`)
	//}
	//
	//self.Opt.Sep.Dump(log.Writer())
}

func (self *Cmd) Has() bool {
	defer gg.Lock(self).Unlock()
	return self.Cmd != nil
}

func (self *Cmd) Broadcast(sig syscall.Signal) {
	defer gg.Lock(self).Unlock()
	self.BroadcastUnsync(sig)
}

/**
Sends the signal to the subprocess group, denoted by the negative sign on the
PID. Requires `syscall.SysProcAttr{Setpgid: true}`.
*/
func (self *Cmd) BroadcastUnsync(sig syscall.Signal) {
	proc := self.ProcUnsync()
	if proc != nil {
		gg.Nop1(syscall.Kill(-proc.Pid, sig))
	}
}

func (self *Cmd) WriteChar(char byte) {
	defer gg.Lock(self).Unlock()

	stdIn := self.Stdin
	if stdIn == nil {
		return
	}

	buf := &self.Buf
	buf[0] = char

	_, err := stdIn.Write(buf[:])
	if err == nil {
		return
	}

	if errors.Is(err, os.ErrClosed) {
		self.Stdin = nil
		return
	}

	panic(err)
}

func (c *Cmd) ProcUnsync() *os.Process {
	cmd := c.Cmd
	if cmd != nil {
		return cmd.Process
	}
	return nil
}

func (c *Cmd) MakeCmd() *exec.Cmd {
	cmd := exec.Command(c.Command, c.Args...)

	// Causes the OS to assign process group ID = `cmd.Process.Pid`.
	// We use this to broadcast signals to the entire subprocess group.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
