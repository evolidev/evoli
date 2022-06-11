package reload

import (
	e "errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/mitranim/gg"
)

type Cmd struct {
	sync.Mutex
	Buf   [1]byte
	Cmd   *exec.Cmd
	Stdin io.WriteCloser
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

func (self *Cmd) Restart(main *Main) {
	defer gg.Lock(self).Unlock()

	self.DeinitUnsync()

	cmd := main.Opt.MakeCmd()
	stdin, err := cmd.StdinPipe()
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

	self.Cmd = cmd
	self.Stdin = stdin
	go main.CmdWait(cmd)
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

	stdin := self.Stdin
	if stdin == nil {
		return
	}

	buf := &self.Buf
	buf[0] = char

	_, err := stdin.Write(buf[:])
	if err == nil {
		return
	}

	if e.Is(err, os.ErrClosed) {
		self.Stdin = nil
		return
	}

	panic(err)
}

func (self *Cmd) ProcUnsync() *os.Process {
	cmd := self.Cmd
	if cmd != nil {
		return cmd.Process
	}
	return nil
}
