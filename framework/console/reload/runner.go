package reload

import (
	"bytes"
	"errors"
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

func (m *Manager) getCommandArguments() (string, []string) {
	//bp := m.FullBuildPath()
	parsed, err := parseCommandLine(m.Command)

	if err != nil {
		m.Logger.Error("Failed to start command:", m.Command)
		panic(err)
	}

	return parsed[0], parsed[1:]
	//args := []string{"run", m.Command}
	//return append(args, m.CommandFlags...)
}

func (m *Manager) getCommand() *exec.Cmd {
	command, args := m.getCommandArguments()
	return exec.Command(command, args...)
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

func parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(fmt.Sprintf("Unclosed quote in command line: %s", command))
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}
