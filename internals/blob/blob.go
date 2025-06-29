package blob

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Write writes blob content to git objects directory and returns the hash
func Write(content []byte) (string, error) {
	hash := CalculateHash(content)

	blobFolder := hash[:2]
	blobFile := hash[2:]
	blobFolderPath := filepath.Join(".git", "objects", blobFolder)
	blobFilePath := filepath.Join(blobFolderPath, blobFile)

	if err := os.MkdirAll(blobFolderPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create folder %s: %w", blobFolderPath, err)
	}

	f, err := os.Create(blobFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", blobFilePath, err)
	}
	defer f.Close()

	// Create the git object format: "blob <size>\0<content>"
	header := fmt.Sprintf("blob %d\x00", len(content))
	fullContent := append([]byte(header), content...)

	w := zlib.NewWriter(f)
	if _, err := w.Write(fullContent); err != nil {
		return "", fmt.Errorf("compression error: %w", err)
	}
	w.Close()

	return hash, nil
}

// CalculateHash calculates the SHA1 hash of blob content in git format
func CalculateHash(content []byte) string {
	header := fmt.Sprintf("blob %d\x00", len(content))
	fullContent := append([]byte(header), content...)
	hash := sha1.Sum(fullContent)
	return fmt.Sprintf("%x", hash)
}

// Read reads the content of a git object file and returns the content
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
