package executor

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/md-talim/codecrafters-git-go/internals/object"
)

type WriteTreeCommand struct{}

func (l *WriteTreeCommand) Execute() error {
	hash, err := l.writeTree(".")
	if err != nil {
		return err
	}

	fmt.Print(hash)
	return nil
}

func (w *WriteTreeCommand) writeTree(dirPath string) (string, error) {
	entries, err := w.collectEntries(dirPath)
	if err != nil {
		return "", err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	treeContent := w.createTreeContent(entries)

	hash, err := object.Write(treeContent, "tree")
	if err != nil {
		return "", fmt.Errorf("failed to write tree object: %w", err)
	}

	return hash, nil
}

func (w *WriteTreeCommand) collectEntries(dirPath string) ([]TreeEntry, error) {
	var entries []TreeEntry

	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, entry := range dirEntries {
		if entry.Name() == ".git" {
			continue
		}

		entryPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			subTreeHash, err := w.writeTree(entryPath)
			if err != nil {
				return nil, err
			}

			entries = append(entries, TreeEntry{
				Mode: "40000",
				Type: "tree",
				Name: entry.Name(),
				Hash: subTreeHash,
			})
		} else {
			fileContent, err := os.ReadFile(entryPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", entryPath, err)
			}

			blobHash, err := object.Write(fileContent, "blob")
			if err != nil {
				return nil, fmt.Errorf("failed to write blob for file %s: %w", entryPath, err)
			}

			mode, err := w.getFileMode(entryPath)
			if err != nil {
				return nil, err
			}

			entries = append(entries, TreeEntry{
				Mode: mode,
				Type: "blob",
				Name: entry.Name(),
				Hash: blobHash,
			})
		}
	}

	return entries, nil
}

func (w *WriteTreeCommand) getFileMode(filePath string) (string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	mode := info.Mode()

	if mode&fs.ModeSymlink != 0 {
		return "120000", nil // Symbolic link mode
	}

	if mode&0111 != 0 {
		return "100755", nil // Executable file mode
	}

	return "100644", nil // Regular file mode
}

func (w *WriteTreeCommand) createTreeContent(entries []TreeEntry) []byte {
	var content []byte

	for _, entry := range entries {
		entryLine := fmt.Sprintf("%s %s", entry.Mode, entry.Name)
		content = append(content, []byte(entryLine)...)
		content = append(content, 0)

		hashBytes, err := hex.DecodeString(entry.Hash)
		if err != nil {
			panic(fmt.Sprintf("Invalid hash: %s", entry.Hash))
		}
		content = append(content, hashBytes...)
	}

	return content
}
