package git

import (
	"fmt"
	"os"

	"github.com/md-talim/codecrafters-git-go/internals/executor"
)

type GitClient struct{}

func NewGitClient() *GitClient {
	return &GitClient{}
}

func (g *GitClient) Run() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	commandExecutor, err := executor.GetCommandExecutor(os.Args[1])
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	if err = commandExecutor.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
}
