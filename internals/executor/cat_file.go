package executor

import (
	"errors"
	"os"
	"path"

	"github.com/md-talim/codecrafters-git-go/internals/blob"
)

const (
	FlagPrettyPrint string = "-p"
)

type CatFileCommand struct{}

func (c *CatFileCommand) Execute() error {
	if len(os.Args) < 4 {
		return errors.New("usage: git cat-file -p <commit-sha>")
	}

	flag := os.Args[2]
	commitSHA := os.Args[3]

	switch flag {
	case FlagPrettyPrint:
		return c.prettyPrint(commitSHA)
	}
	return nil
}

func (c *CatFileCommand) prettyPrint(commitSHA string) error {
	objectFilePath := getGitObjectPath(commitSHA)
	blobContent, err := blob.Read(objectFilePath)
	if err != nil {
		return err
	}
	os.Stdout.Write(blobContent)
	return err
}

func getGitObjectPath(commitSHA string) string {
	folder := commitSHA[0:2]
	file := commitSHA[2:]

	gitObjectPath := path.Join(".git", "objects", folder, file)
	return gitObjectPath
}
