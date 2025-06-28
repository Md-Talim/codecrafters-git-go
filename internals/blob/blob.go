package blob

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
)

func Read(objectFilePath string) ([]byte, error) {
	if _, err := os.Stat(objectFilePath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("not a valid object name %s", objectFilePath)
	}

	fileContents, err := os.ReadFile(objectFilePath)
	if err != nil {
		return nil, err
	}

	r, err := zlib.NewReader(bytes.NewReader(fileContents))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	outputBuffer, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parts := bytes.SplitN(outputBuffer, []byte{0}, 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("malformed git object")
	}

	return parts[1], nil
}
