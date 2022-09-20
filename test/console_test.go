package test

import (
	"github.com/evolidev/evoli/framework/console"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSimpleCommand(t *testing.T) {
	t.Parallel()
	t.Run("Parse simple command with required parameter", func(t *testing.T) {
		command := "mail:send foo"
		definition := "mail:send {user}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "foo", cmd.GetArgument("user"))
	})

	t.Run("Parse simple command with optional parameter", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user?}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "", cmd.GetArgument("user"))
	})

	t.Run("Parse simple command with optional parameter", func(t *testing.T) {
		command := "mail:send foo"
		definition := "mail:send {user?}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "foo", cmd.GetArgument("user"))
	})

	t.Run("Parse simple command with optional parameter and default value", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user=foo}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "foo", cmd.GetArgument("user"))
	})

	t.Run("Parse command and pass options", func(t *testing.T) {
		command := "mail:send foo --queue"
		definition := "mail:send {user} {--queue}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, true, cmd.GetOption("queue").Bool())
	})

	t.Run("Parse command and pass options", func(t *testing.T) {
		command := "serve"
		definition := "serve {--port}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, 1010, cmd.GetOptionWithDefault("port", 1010).Integer())
	})

	t.Run("Parse command and pass required option", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user} {--queue=}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "", cmd.GetOption("queue").String())
	})

	t.Run("Parse command and pass option and alias", func(t *testing.T) {
		command := "mail:send -Q"
		definition := "mail:send {user} {--Q|queue}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, true, cmd.GetOption("Q").Bool())
	})

	t.Run("Get name of command", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user} {--Q|queue}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "mail:send", cmd.GetName())
	})

	t.Run("Get prefix of command", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user} {--Q|queue}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "mail", cmd.GetPrefix())
	})

	t.Run("Get subcommand of command", func(t *testing.T) {
		command := "mail:send"
		definition := "mail:send {user} {--Q|queue}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "send", cmd.GetSubCommand())
	})

	t.Run("Get empty subcommand of command", func(t *testing.T) {
		command := "mail"
		definition := "mail:send {user} {--Q|queue}"

		cmd := console.Parse(definition, command)

		assert.Equal(t, "", cmd.GetSubCommand())
	})

}
