package cli

import (
	"blog-go/internal/post"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all blog posts",
	Long:  "Display a list of all blog posts with their ID, title, and creation date.",
	Run:   listPost,
}

func init() {
	listCmd.Flags().Bool("notes", false, "Show only notes")
	listCmd.Flags().Bool("proj", false, "Show only projects")

	rootCmd.AddCommand(listCmd)
}

func listPost(cmd *cobra.Command, args []string) {
	catNotes, _ := cmd.Flags().GetBool("notes")
	catProj, _ := cmd.Flags().GetBool("proj")

	var cat post.Category = post.ALL

	if catNotes {
		cat = post.NOTES
	} else if catProj {
		cat = post.PROJECTS
	}

	posts, err := postService.ListCategory(cat)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(posts) == 0 {
		fmt.Println("no post has been uploaded yet")
		return
	}

	for _, post := range posts {
		title := post.Title

		if len(title) > 40 {
			title = title[:37] + "..."
		}

		fmt.Printf("%-4d %-7s %-40s %-20s %s\n", post.Id, post.Category, title, post.CreatedAt.Format("2006-01-02 15:04"), strings.Join(post.Tags, ","))
	}
}
