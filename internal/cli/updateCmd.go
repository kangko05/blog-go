package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing blog post",
	Long:  "Update an existing blog post by opening it in your preferred editor. The post content will be loaded for editing.",
	Args:  cobra.ExactArgs(1),
	Run:   updatePost,
}

func init() {
	updateCmd.Flags().StringSliceP("tags", "t", nil, "update tags for the post")
	rootCmd.AddCommand(updateCmd)
}

func updatePost(cmd *cobra.Command, args []string) {
	tags, _ := cmd.Flags().GetStringSlice("tags")

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	post, err := postService.GetPost(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	tmpFile, err := os.CreateTemp("", "blog-post-*.md")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(post.Content); err != nil {
		fmt.Println(err)
		return
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	command := exec.Command(editor, tmpFile.Name())

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Println(err)
		return
	}

	rb, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Println(err)
		return
	}

	title := extractTitle(rb)

	if err := postService.UpdatePost(id, title, string(rb), tags); err != nil {
		fmt.Printf("update cancelled: %v\n", err)
		return
	}

	fmt.Printf("update post id %d success\n", id)
}
