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
	objectFolder, objectPath := getObjectPaths(hash)

	if err := os.MkdirAll(objectFolder, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create folder %s: %w", objectFolder, err)
	}

	f, err := os.Create(objectPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", objectPath, err)
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
	fullContent := createGitObjectFormat(content)
	hash := sha1.Sum(fullContent)
	return fmt.Sprintf("%x", hash)
}

// ReadRaw reads the raw content of a git object file (without extracting blob content)
func ReadRaw(hash string) ([]byte, error) {
	_, objectPath := getObjectPaths(hash)
	if _, err := os.Stat(objectPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("not a valid object name %s", hash)
	}

	fileContents, err := os.ReadFile(objectPath)
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

	return outputBuffer, nil
}

// Read reads the content of a git object file and returns the content
func Read(hash string) ([]byte, error) {
	_, objectPath := getObjectPaths(hash)
	if _, err := os.Stat(objectPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("not a valid object name %s", hash)
	}

	fileContents, err := os.ReadFile(objectPath)
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

	return extractContentFromGitObject(outputBuffer)
}

func extractContentFromGitObject(data []byte) ([]byte, error) {
	parts := bytes.SplitN(data, []byte{0}, 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("malformed git object")
	}
	return parts[1], nil
}

func getObjectPaths(hash string) (string, string) {
	if len(hash) < 3 {
		return "", ""
	}

	objectDir := hash[:2]
	objectFile := hash[2:]
	objectDirPath := filepath.Join(".git", "objects", objectDir)
	objectFilePath := filepath.Join(objectDirPath, objectFile)
	return objectDirPath, objectFilePath
}

func createGitObjectFormat(content []byte) []byte {
	header := fmt.Sprintf("blob %d\x00", len(content))
	return append([]byte(header), content...)
}
