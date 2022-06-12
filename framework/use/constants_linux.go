package use

import "golang.org/x/sys/unix"

const IoctlReadTermIos = unix.TCGETS
const IoctlWriteTermIos = unix.TCSETS
