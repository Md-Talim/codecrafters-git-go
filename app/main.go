package main

import (
	"github.com/md-talim/codecrafters-git-go/pkg/git"
)

func main() {
	gitClient := git.NewGitClient()
	gitClient.Run()
}
