package executor

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/md-talim/codecrafters-git-go/internals/blob"
)

const FlagNameOnly string = "--name-only"

type LSTreeCommand struct{}

type TreeEntry struct {
	Mode string
	Type string
	Hash string
	Name string
}

func (l *LSTreeCommand) Execute() error {
	if len(os.Args) < 3 {
		return errors.New("usage: git ls-tree <tree-sha>")
	}

	var nameOnly bool
	var treeSHA string

	if len(os.Args) == 4 && os.Args[2] == FlagNameOnly {
		nameOnly = true
		treeSHA = os.Args[3]
	} else {
		nameOnly = false
		treeSHA = os.Args[2]
	}

	entries, err := l.parseTreeObject(treeSHA)
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	if nameOnly {
		for _, entry := range entries {
			fmt.Println(entry.Name)
		}
	} else {
		for _, entry := range entries {
			fmt.Printf("%s %s %s\t%s\n", entry.Mode, entry.Type, entry.Hash, entry.Name)
		}
	}
	return nil
}

func (l *LSTreeCommand) parseTreeObject(hash string) ([]TreeEntry, error) {
	rawContent, err := blob.ReadRaw(hash)
	if err != nil {
		return nil, err
	}

	content, err := l.extractTreeContent(rawContent)
	if err != nil {
		return nil, err
	}

	entries, err := l.parseTreeEntries(content)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (l *LSTreeCommand) extractTreeContent(data []byte) ([]byte, error) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, errors.New("malformed tree object: no null byte found")
	}

	header := string(data[:nullIndex])
	if !strings.HasPrefix(header, "tree ") {
		return nil, fmt.Errorf("not a tree object")
	}

	return data[nullIndex+1:], nil
}

func (l *LSTreeCommand) parseTreeEntries(content []byte) ([]TreeEntry, error) {
	var entries []TreeEntry
	offset := 0

	for offset < len(content) {
		// Find the space that separates mode from name
		spaceIndex := bytes.IndexByte(content[offset:], ' ')
		if spaceIndex == -1 {
			break
		}
		spaceIndex += offset
		mode := string(content[offset:spaceIndex])

		// Find the null byte that separates name from hash
		nullIndex := bytes.IndexByte(content[spaceIndex+1:], 0)
		if nullIndex == -1 {
			return nil, fmt.Errorf("malformed tree entry: no null byte found")
		}
		nullIndex += spaceIndex + 1
		name := string(content[spaceIndex+1 : nullIndex])

		if nullIndex+21 > len(content) {
			return nil, fmt.Errorf("malformed tree: incomplete hash")
		}
		hashBytes := content[nullIndex+1 : nullIndex+21]
		hashHex := hex.EncodeToString(hashBytes)

		// Determine type based on mode
		var objectType string
		modeInt, err := strconv.Atoi(mode)
		if err != nil {
			return nil, fmt.Errorf("invalid mode %s: %v", mode, err)
		}

		switch modeInt {
		case 40000:
			objectType = "tree"
		case 100644, 100755, 120000:
			objectType = "blob"
		default:
			objectType = "unknown"
		}

		entries = append(entries, TreeEntry{
			Mode: mode,
			Type: objectType,
			Hash: hashHex,
			Name: name,
		})

		offset = nullIndex + 21
	}

	return entries, nil
}
