package use

import (
	"github.com/Code-Hex/dd"
	"github.com/Code-Hex/dd/p"
	"log"
)

func AbortUnless(e interface{}) {
	if e != nil && e != false {
		panic(e)
	}
}

func D(e ...any) {
	log.Println(dd.Dump(e))
}

func P(e ...any) {
	log.Println(p.P(e))
}
