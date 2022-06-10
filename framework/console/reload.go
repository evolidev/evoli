package console

import (
	"log"
	"path/filepath"

	"github.com/mitranim/gg"
	"github.com/rjeczalik/notify"
)

func Watch() {
	var watcher WatchNotify

	watcher.Init()
	watcher.Run()
}

type WatchNotify struct {
	Done   gg.Chan[struct{}]
	Events gg.Chan[notify.EventInfo]
}

func (self *WatchNotify) Init() {
	self.Done.Init()
	self.Events.InitCap(1)

	paths := []string{"/Users/omer/Code/evoli/examples"}
	verbose := true

	for _, path := range paths {
		path = filepath.Join(path, `...`)
		if verbose {
			log.Printf(`watching %q`, path)
		}
		gg.Try(notify.Watch(path, self.Events, notify.All))
	}
}

func (self *WatchNotify) Deinit() {
	self.Done.SendZeroOpt()
	if self.Events != nil {
		notify.Stop(self.Events)
	}
}

func (self WatchNotify) Run() {
	for {
		select {
		case <-self.Done:
			return

		case event := <-self.Events:
			log.Println(event)
			//if main.Opt.ShouldRestart(event) {
			//	if main.Opt.Verb {
			//		log.Println(`restarting on FS event:`, event)
			//	}
			//	main.Restart()
			//}
		}
	}
}
