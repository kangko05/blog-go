package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all blog posts",
	Long:  "Display a list of all blog posts with their ID, title, and creation date.",
	Run:   listPost,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listPost(cmd *cobra.Command, args []string) {
	posts, err := postService.ListAllPosts()
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

		fmt.Printf("%-4d\t%-40s\t%v\n", post.Id, title, post.CreatedAt.Format("2006-01-02 15:04"))
	}
}
