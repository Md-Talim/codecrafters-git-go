package executor

import (
	"errors"
	"os"

	"github.com/md-talim/codecrafters-git-go/internals/object"
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
	blobContent, err := object.ReadContent(commitSHA)
	if err != nil {
		return err
	}
	os.Stdout.Write(blobContent)
	return nil
}
