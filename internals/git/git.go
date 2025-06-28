package git

import (
	"fmt"
	"os"
)

type Command interface {
	execute() error
}

type Git struct{}

func (g *Git) Run(command Command) {
	if err := command.execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
}
