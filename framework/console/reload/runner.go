package reload

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func (m *Manager) runner() {
	var cmd *exec.Cmd

	for {
		<-m.Restart

		if cmd != nil && cmd.Process != nil {
			// kill the previous command
			pid := cmd.Process.Pid
			m.Logger.Success("Stopping: PID %d", pid)

			process, err := os.FindProcess(int(pid))
			if err != nil {
				fmt.Printf("Failed to find process: %s\n", err)
			} else {
				log.Println(process)
				cmd.Process.Kill()
				//err := cmd.Process.Signal(syscall.Signal(0))
				fmt.Printf("process.Signal on pid %d returned: %v\n", pid, err)
			}
		} else {
			m.Logger.Print("No process running")
		}

		cmd = m.getCommand()

		go func() {
			err := m.runAndListen(cmd)
			if err != nil {
				m.Logger.Error(err)
			}
		}()
	}
}

func (m *Manager) getCommandArguments() []string {
	bp := m.FullBuildPath()
	args := []string{"run", bp}
	return append(args, m.CommandFlags...)
}

func (m *Manager) getCommand() *exec.Cmd {
	return exec.Command("go", m.getCommandArguments()...)
}

func (m *Manager) runAndListen(cmd *exec.Cmd) error {
	cmd.Stderr = m.Stderr
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}

	cmd.Stdin = m.Stdin
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = m.Stdout
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}

	var stderr bytes.Buffer

	cmd.Stderr = io.MultiWriter(&stderr, cmd.Stderr)

	// Set the environment variables from config
	if len(m.CommandEnv) != 0 {
		cmd.Env = append(m.CommandEnv, os.Environ()...)
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stderr.String())
	}

	m.Logger.Success("Running: %s (PID: %d)", strings.Join(cmd.Args, " "), cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stderr.String())
	}
	return nil
}
