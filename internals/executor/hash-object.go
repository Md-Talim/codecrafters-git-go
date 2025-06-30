package executor

import (
	"errors"
	"fmt"
	"os"

	"github.com/md-talim/codecrafters-git-go/internals/object"
)

const (
	OptionWrite string = "-w"
)

type HashObjectCommand struct{}

func (h *HashObjectCommand) Execute() error {
	if len(os.Args) < 3 {
		return errors.New("usage: git hash-object [-w] <file>")
	}

	var option string = ""
	var file string = ""
	if len(os.Args) >= 4 {
		option = os.Args[2]
		file = os.Args[3]
	} else {
		file = os.Args[2]
	}

	fileContents, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("could not open %s for reading: %w", file, err)
	}

	switch option {
	case OptionWrite:
		hash, err := object.Write(fileContents, "blob")
		if err != nil {
			return err
		}
		fmt.Print(hash)
	default:
		hash := object.CalculateHash(fileContents, "blob")
		fmt.Print(hash)
	}

	return nil
}
