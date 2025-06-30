package executor

import "errors"

const (
	Init       = "init"
	CatFile    = "cat-file"
	CommitTree = "commit-tree"
	HashObject = "hash-object"
	LSTree     = "ls-tree"
	WriteTree  = "write-tree"
)

type CommandExecutor interface {
	Execute() error
}

var availableCommands map[string]CommandExecutor = map[string]CommandExecutor{
	Init:       &InitCommand{},
	CatFile:    &CatFileCommand{},
	CommitTree: &CommitTreeCommand{},
	HashObject: &HashObjectCommand{},
	LSTree:     &LSTreeCommand{},
	WriteTree:  &WriteTreeCommand{},
}

func GetCommandExecutor(commandname string) (CommandExecutor, error) {
	command, ok := availableCommands[commandname]
	if !ok {
		return nil, errors.New(`unknown command ${command}\n`)
	}
	return command, nil
}
