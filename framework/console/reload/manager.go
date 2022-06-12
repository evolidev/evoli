package reload

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Manager struct {
	*Configuration
	ID         string
	Logger     *Logger
	Restart    chan bool
	cancelFunc context.CancelFunc
	context    context.Context
	gil        *sync.Once
}

func New(c *Configuration) *Manager {
	return NewWithContext(c, context.Background())
}

func NewWithContext(c *Configuration, ctx context.Context) *Manager {
	ctx, cancelFunc := context.WithCancel(ctx)
	m := &Manager{
		Configuration: c,
		ID:            ID(),
		Logger:        NewLogger(c),
		Restart:       make(chan bool),
		cancelFunc:    cancelFunc,
		context:       ctx,
		gil:           &sync.Once{},
	}
	return m
}

func (m *Manager) Start() error {
	w := NewWatcher(m)
	w.Start()
	go m.build(fsnotify.Event{Name: ":start:"})

	if !m.Debug {
		go func() {
			for {
				select {
				case event := <-w.Events():
					//log.Println("received event", event)
					if !w.isFileEligibleForChange(event.Name) {
						continue
					}

					if event.Op != fsnotify.Chmod {
						go m.build(event)
					}

					if w.ForcePolling {
						//w.Logger.Print("Removing file from watchlist: %s", event.Name)
						w.Remove(event.Name)
						w.Add(event.Name)
					}

				case <-m.context.Done():
					m.Logger.Print("Shutting down")
					break
				}
			}
		}()
	}
	go func() {
		for {
			select {
			case err := <-w.Errors():
				m.Logger.Error(err)
			case <-m.context.Done():
				break
			}
		}
	}()
	m.runner()
	return nil
}

func (m *Manager) build(event fsnotify.Event) {
	m.gil.Do(func() {
		defer func() {
			m.gil = &sync.Once{}
		}()

		m.buildTransaction(func() error {
			// time.Sleep(r.BuildDelay * time.Millisecond)

			now := time.Now()
			m.Logger.Print("Rebuild on: %s", event.Name)

			command, args := m.getCommandArguments()
			cmd := exec.CommandContext(m.context, command, args...)
			cmd.Dir = m.AppRoot

			err := m.runAndListen(cmd)
			if err != nil {
				if strings.Contains(err.Error(), "no buildable Go source files") {
					m.cancelFunc()
					log.Fatal(err)
				}
				return err
			}

			tt := time.Since(now)
			m.Logger.Success("Buildings Completed (PID: %d) (Time: %s)", cmd.Process.Pid, tt)
			//m.Restart <- true
			return nil
		})
	})
}

func (m *Manager) buildTransaction(fn func() error) {
	logPath := ErrorLogPath()
	err := fn()
	if err != nil {
		f, _ := os.Create(logPath)
		fmt.Fprint(f, err)
		m.Logger.Error("Error!")
		m.Logger.Error(err)
	} else {
		os.Remove(logPath)
	}
}
