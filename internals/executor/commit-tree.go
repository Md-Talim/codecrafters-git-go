package executor

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/md-talim/codecrafters-git-go/internals/object"
)

type CommitTreeCommand struct{}

func (c *CommitTreeCommand) Execute() error {
	if len(os.Args) < 6 {
		return errors.New("usage: git commit-tree <tree-sha> -p <parent-sha> -m <message>")
	}

	treeSHA := os.Args[2]
	var parentSHA, message string

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-p":
			if i+1 >= len(os.Args) {
				return errors.New("missing parent SHA after -p")
			}
			parentSHA = os.Args[i+1]
			i++
		case "-m":
			if i+1 >= len(os.Args) {
				return errors.New("missing commit message after -m")
			}
			message = os.Args[i+1]
			i++
		}
	}

	if parentSHA == "" {
		return errors.New("parent commit SHA is required")
	}
	if message == "" {
		return errors.New("commit message is required")
	}

	commitContent := c.createCommitContent(treeSHA, parentSHA, message)

	hash, err := object.Write(commitContent, "commit")
	if err != nil {
		return fmt.Errorf("failed to write commit object: %w", err)
	}

	fmt.Print(hash)

	return nil
}

func (c *CommitTreeCommand) createCommitContent(treeSHA, parentSHA, message string) []byte {
	// Get current timestamp
	now := time.Now()
	timestamp := now.Unix()
	timezone := now.Format("-7000") // Format timezone as +/-HHMM

	name := "Md Talim"
	email := "the.mohd.talim@gmail.com"

	content := fmt.Sprintf("tree %s\n", treeSHA)
	content += fmt.Sprintf("parent %s\n", parentSHA)
	content += fmt.Sprintf("author %s <%s> %d %s\n", name, email, timestamp, timezone)
	content += fmt.Sprintf("committer %s <%s> %d %s\n", name, email, timestamp, timezone)
	content += fmt.Sprintf("\n%s\n", message)

	return []byte(content)
}
