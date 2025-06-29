package executor

import "errors"

const (
	Init       = "init"
	CatFile    = "cat-file"
	HashObject = "hash-object"
)

type CommandExecutor interface {
	Execute() error
}

var availableCommands map[string]CommandExecutor = map[string]CommandExecutor{
	Init:       &InitCommand{},
	CatFile:    &CatFileCommand{},
	HashObject: &HashObjectCommand{},
}

func GetCommandExecutor(commandname string) (CommandExecutor, error) {
	command, ok := availableCommands[commandname]
	if !ok {
		return nil, errors.New(`unknown command ${command}\n`)
	}
	return command, nil
}
