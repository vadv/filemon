package server

import (
	"fmt"
)

type aviableCommands map[string]command

type command interface {
	New(string, []string) (command, error)
	Process(string)
	Result() float64
}

var ArgsError = fmt.Errorf("COMMAND_ARGUMENTS_ERROR")
var aviable = make(aviableCommands, 0)

func registerCommand(name string, cmd command) {
	aviable[name] = cmd
}

func newCommand(name, expr string, args []string) (command, error) {
	cmd, ok := aviable[name]
	if !ok {
		return nil, fmt.Errorf("COMMAND_NOT_FOUND")
	}
	return cmd.New(expr, args)
}
