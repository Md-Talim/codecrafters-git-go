package git

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

const (
	FlagPrettyPrint = "-p"
)

type CatFileCommand struct {
	flag      string
	commitSHA string
}

func NewCatFileCommand(flag, commitSHA string) *CatFileCommand {
	return &CatFileCommand{flag, commitSHA}
}

func (c *CatFileCommand) execute() error {
	switch c.flag {
	case FlagPrettyPrint:
		return c.handlePrettyPrint()
	}

	return nil
}

func (c *CatFileCommand) handlePrettyPrint() error {
	objectFilePath := getGitObjectPath(c.commitSHA)

	if _, err := os.Stat(objectFilePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("not a valid object name %s", objectFilePath)
	}

	fileContents, err := os.ReadFile(objectFilePath)
	if err != nil {
		return err
	}

	r, err := zlib.NewReader(bytes.NewReader(fileContents))
	if err != nil {
		return err
	}
	defer r.Close()

	outputBuffer, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	parts := bytes.SplitN(outputBuffer, []byte{0}, 2)
	if len(parts) < 2 {
		return fmt.Errorf("malformed git object")
	}

	_, err = os.Stdout.Write(parts[1])

	return err
}

func getGitObjectPath(commitSHA string) string {
	folder := commitSHA[0:2]
	file := commitSHA[2:]

	gitObjectPath := path.Join(".git", "objects", folder, file)
	return gitObjectPath
}
