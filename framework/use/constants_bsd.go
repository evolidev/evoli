//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package use

import "golang.org/x/sys/unix"

const IoctlReadTermIos = unix.TIOCGETA
const IoctlWriteTermIos = unix.TIOCSETA
