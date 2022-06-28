package reload

import (
	"context"
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/use"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Manager struct {
	*Configuration
	ID         string
	Logger     *logging.Logger
	Restart    chan bool
	cancelFunc context.CancelFunc
	context    context.Context
	gil        *sync.Once
	cmd        *exec.Cmd
}

func New(c *Configuration) *Manager {
	return NewWithContext(c, context.Background())
}

func NewWithContext(c *Configuration, ctx context.Context) *Manager {
	ctx, cancelFunc := context.WithCancel(ctx)
	m := &Manager{
		Configuration: c,
		ID:            ID(),
		Logger: logging.NewLogger(&logging.Config{
			Name:         "reload",
			EnableColors: true,
			PrefixColor:  148,
		}),
		Restart:    make(chan bool),
		cancelFunc: cancelFunc,
		context:    ctx,
		gil:        &sync.Once{},
	}
	return m
}

func (m *Manager) Start() error {
	w := NewWatcher(m)
	w.Start()

	restart := func() {
		m.Restart <- true
	}

	debounced := use.Debounce(100 * time.Millisecond)

	if !m.Debug {
		go func() {
			for {
				select {
				case event := <-w.Events():
					if !w.isFileEligibleForChange(event.Name) {
						continue
					}

					if event.Op != fsnotify.Chmod {
						debounced(restart)
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
				m.Logger.Error("Manager error", err)
			case <-m.context.Done():
				break
			}
		}
	}()

	m.runner()
	return nil
}

func (m *Manager) build() *exec.Cmd {

	timer := use.TimeRecord()
	//m.Logger.Print("Rebuild on: %s", event.Name)

	command, args := m.getCommandArguments()
	cmd := exec.CommandContext(m.context, command, args...)
	cmd.Dir = m.AppRoot

	err := m.runAndListen(cmd)
	m.cmd = cmd
	if err != nil {
		if strings.Contains(err.Error(), "no buildable Go source files") {
			m.cancelFunc()
			log.Fatal(err)
		}
		return nil
	}

	m.Logger.Success("Buildings Completed (PID: %d) %s",
		cmd.Process.Pid,
		timer.ElapsedColored(),
	)

	return cmd
}
