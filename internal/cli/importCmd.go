package cli

import (
	"blog-go/internal/post"
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
	importCmd.Flags().Bool("proj", false, "set category to projects")
	importCmd.Flags().Bool("notes", false, "set category to notes")
	rootCmd.AddCommand(importCmd)
}

func importPost(cmd *cobra.Command, args []string) {
	catNotes, _ := cmd.Flags().GetBool("notes")
	catProj, _ := cmd.Flags().GetBool("proj")

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

	rb, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := postService.CreatePost(cat, extractTitle(rb), string(rb)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("create post success")
}
