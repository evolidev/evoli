package reload

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evolidev/evoli/framework/console/filenotify"
)

type Watcher struct {
	filenotify.FileWatcher
	*Manager
	context context.Context
}

func NewWatcher(r *Manager) *Watcher {
	var watcher filenotify.FileWatcher

	if r.ForcePolling {
		watcher = filenotify.NewPollingWatcher()
	} else {
		watcher, _ = filenotify.NewEventWatcher()
	}

	return &Watcher{
		FileWatcher: watcher,
		Manager:     r,
		context:     r.context,
	}
}

func (w *Watcher) Start() {
	if w.ForcePolling {
		w.watchWithPolling()
	} else {
		w.watchWithFsNotify()
	}
}

func (w *Watcher) watchWithFsNotify() {
	watchDir := w.AppRoot

	if err := filepath.WalkDir(watchDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return w.FileWatcher.Add(path)
		}

		return nil
	}); err != nil {
		w.Logger.Error("Watch FS", err)
	}
}

func (w *Watcher) watchWithPolling() {
	go func() {
		for {
			err := filepath.Walk(w.AppRoot, func(path string, info os.FileInfo, err error) error {
				//w.Logger.Print(fmt.Sprintf("Check file: %s", path))

				if info == nil {
					w.cancelFunc()
					return errors.New("nil directory")
				}
				if info.IsDir() {
					if strings.HasPrefix(filepath.Base(path), "_") {
						return filepath.SkipDir
					}
					if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") || w.isIgnoredFolder(path) {
						return filepath.SkipDir
					}
				}
				if w.isWatchedFile(path) {
					//w.Logger.Print(fmt.Sprintf("Add file: %s", path))
					w.Add(path)
				}
				return nil
			})

			if err != nil {
				w.context.Done()
				break
			}
			// sweep for new files every 1 second
			time.Sleep(1 * time.Second)
		}
	}()
}

func (w *Watcher) isIgnoredFolder(path string) bool {
	paths := strings.Split(path, "/")
	if len(paths) <= 0 {
		return false
	}

	for _, e := range w.IgnoredFolders {
		if strings.TrimSpace(e) == paths[0] {
			return true
		}
	}
	return false
}

func (w *Watcher) isWatchedFile(path string) bool {
	ext := filepath.Ext(path)

	for _, e := range w.IncludedExtensions {
		if strings.TrimSpace(e) == ext {
			return true
		}
	}

	return false
}

func (w *Watcher) isFileEligibleForChange(path string) bool {
	//w.Logger.Print(fmt.Sprintf("isFileEligibleForChange: %s", path))
	// check if the last character of path is tilde and replace it
	if strings.HasSuffix(path, "~") {
		path = strings.Replace(path, "~", "", -1)
	}

	info, err := os.Stat(path)
	if err != nil {
		w.Logger.Error("Is file eligible: %s", err)
		return false
	}

	if info == nil {
		//w.cancelFunc() //??
		w.Logger.Error("info not found")
		return false
	}

	if info.IsDir() {
		basePath := filepath.Base(path)
		if strings.HasPrefix(basePath, "_") {
			return false
		}
		if len(path) > 1 && strings.HasPrefix(basePath, ".") || w.isIgnoredFolder(path) {
			return false
		}
	}

	if !w.isWatchedFile(path) {
		return false
	}

	return true
}
