package object

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
func Write(content []byte, objectType string) (string, error) {
	hash := CalculateHash(content, objectType)
	objectFolder, objectPath, err := getObjectPaths(hash)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(objectFolder, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create folder %s: %w", objectFolder, err)
	}

	f, err := os.Create(objectPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", objectPath, err)
	}
	defer f.Close()

	fullContent := createGitObjectFormat(content, objectType)

	w := zlib.NewWriter(f)
	if _, err := w.Write(fullContent); err != nil {
		return "", fmt.Errorf("compression error: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("compression close error: %w", err)
	}

	return hash, nil
}

// CalculateHash calculates the SHA1 hash of blob content in git format
func CalculateHash(content []byte, objectType string) string {
	fullContent := createGitObjectFormat(content, objectType)
	hash := sha1.Sum(fullContent)
	return fmt.Sprintf("%x", hash)
}

// ReadContent reads git object and returns just the content (without header)
func ReadContent(hash string) ([]byte, error) {
	rawData, err := Read(hash)
	if err != nil {
		return nil, err
	}
	return ExtractContent(rawData)
}

// Read reads the content of a git object by its hash
func Read(hash string) ([]byte, error) {
	_, objectPath, err := getObjectPaths(hash)
	if err != nil {
		return nil, err
	}
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

// ExtractContent extracts content from git object format (removes header)
func ExtractContent(data []byte) ([]byte, error) {
	parts := bytes.SplitN(data, []byte{0}, 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("malformed git object")
	}
	return parts[1], nil
}

func getObjectPaths(hash string) (string, string, error) {
	if len(hash) < 3 {
		return "", "", fmt.Errorf("invalid hash length: %d", len(hash))
	}

	objectDir := hash[:2]
	objectFile := hash[2:]
	objectDirPath := filepath.Join(".git", "objects", objectDir)
	objectFilePath := filepath.Join(objectDirPath, objectFile)
	return objectDirPath, objectFilePath, nil
}

func createGitObjectFormat(content []byte, objectType string) []byte {
	header := fmt.Sprintf("%s %d\x00", objectType, len(content))
	return append([]byte(header), content...)
}
