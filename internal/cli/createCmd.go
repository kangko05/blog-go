package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new blog post",
	Long:  "Create a new blog post by opening your preferred editor. The post will be saved to the database after editing.",
	Run:   createPost,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createPost(cmd *cobra.Command, args []string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	tmpFile, err := os.CreateTemp("", "blog-post-*.md")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(tmpFile.Name())

	command := exec.Command(editor, tmpFile.Name())

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Println(err)
		return
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(content) == 0 {
		fmt.Println("empty content, post not created")
		return
	}

	if _, err := postService.CreatePost(extractTitle(content), string(content)); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("create post success")
	}
}

func extractTitle(content []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}

	return "Untitled"
}
