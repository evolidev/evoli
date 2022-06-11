package reload

import (
	"github.com/evolidev/evoli/framework/use"
	"github.com/mitranim/gg"
	"golang.org/x/sys/unix"
)

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

	https://en.wikibooks.org/wiki/Serial_Programming/termios

	man termios
*/
type TermState struct{ gg.Opt[unix.Termios] }

func (t *TermState) Init() {
	t.Deinit()

	state, err := unix.IoctlGetTermios(FdTerm, use.IoctlReadTermIos)
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

	// Seems unnecessary on my system. Might be needed elsewhere.
	// state.Cflag |= unix.CS8
	// state.Cc[unix.VMIN] = 1
	// state.Cc[unix.VTIME] = 0

	err = unix.IoctlSetTermios(FdTerm, use.IoctlWriteTermIos, state)
	if err != nil {
		log.Printf(`unable to switch terminal to raw mode: %v`, err)
		return
	}

	t.Set(prev)
}

func (t *TermState) Deinit() {
	if !t.IsNull() {
		defer t.Clear()
		gg.Nop1(unix.IoctlSetTermios(FdTerm, use.IoctlWriteTermIos, &t.Val))
	}
}
