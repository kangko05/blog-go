package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import <filepath>",
	Short: "Import a markdown file as a blog post",
	Long:  "Import an existing markdown file and save it as a new blog post in the database.",
	Args:  cobra.ExactArgs(1),
	Run:   importPost,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func importPost(cmd *cobra.Command, args []string) {
	rb, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := postService.CreatePost(extractTitle(rb), string(rb)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("create post success")
}
