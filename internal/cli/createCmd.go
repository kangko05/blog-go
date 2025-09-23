package cli

import (
	"blog-go/internal/post"
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
	createCmd.Flags().Bool("proj", false, "set category to projects")
	createCmd.Flags().Bool("notes", false, "set category to notes")
	createCmd.Flags().StringSliceP("tags", "t", nil, "set tags for the post")

	rootCmd.AddCommand(createCmd)
}

func createPost(cmd *cobra.Command, args []string) {
	catNotes, _ := cmd.Flags().GetBool("notes")
	catProj, _ := cmd.Flags().GetBool("proj")
	tags, _ := cmd.Flags().GetStringSlice("tags")

	if !(catNotes || catProj) {
		fmt.Println("please specify a category: --notes or --proj")
		return
	}

	if catNotes && catProj {
		fmt.Println("cannot specify both categories")
		return
	}

	var cat post.Category
	if catNotes {
		cat = post.NOTES
	}
	if catProj {
		cat = post.PROJECTS
	}

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

	if _, err := postService.CreatePost(cat, extractTitle(content), string(content), tags); err != nil {
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
