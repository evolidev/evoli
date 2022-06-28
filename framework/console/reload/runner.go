package reload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/evolidev/evoli/framework/filesystem"
	"github.com/evolidev/evoli/framework/use"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func (m *Manager) runner() {

	go m.build()

	for {
		<-m.Restart
		m.Logger.Log("Restarting...")

		pid := filesystem.Read(use.StoragePath("tmp/serve.pid"))
		if pid != "" {
			m.Logger.Log("Killing process with PID: %s", pid)
			// convert pid to int
			pidInt, _ := strconv.Atoi(pid)
			syscall.Kill(pidInt, syscall.SIGTERM)

		} else {
			m.Logger.Print("No process running")
		}

		go m.build()
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
