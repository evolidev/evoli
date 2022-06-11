package reload

import (
	"path/filepath"

	"github.com/mitranim/gg"
	"github.com/rjeczalik/notify"
)

// Implementation of `Watcher` that uses "github.com/rjeczalik/notify".
type WatchNotify struct {
	Done   gg.Chan[struct{}]
	Events gg.Chan[notify.EventInfo]
}

func (w *WatchNotify) Init(main *Main) {
	w.Done.Init()
	w.Events.InitCap(1)

	paths := main.Opt.Watch
	verb := main.Opt.Verb && !gg.Equal(paths, OptDefault().Watch)

	for _, path := range paths {
		path = filepath.Join(path, `...`)
		if verb {
			log.Printf(`watching %q`, path)
		}
		gg.Try(notify.Watch(path, w.Events, notify.All))
	}
}

func (w *WatchNotify) DeInit() {
	w.Done.SendZeroOpt()
	if w.Events != nil {
		notify.Stop(w.Events)
	}
}

func (w *WatchNotify) Run(main *Main) {
	for {
		select {
		case <-w.Done:
			return

		case event := <-w.Events:
			log.Println("Should restart", main.Opt.ShouldRestart(event))
			if main.Opt.ShouldRestart(event) {
				if main.Opt.Verb {
					log.Println(`restarting on FS event:`, event)
				}
				main.Restart()
			}
		}
	}
}
