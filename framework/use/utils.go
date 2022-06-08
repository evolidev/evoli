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

func D(e interface{}) {
	log.Println(dd.Dump(e))
}

func P(e interface{}) {
	log.Println(p.P(e))
}
