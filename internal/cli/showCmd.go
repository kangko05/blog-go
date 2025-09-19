package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a blog post",
	Long:  "Display the full content of a blog post by its ID.",
	Args:  cobra.ExactArgs(1),
	Run:   showPost,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func showPost(cmd *cobra.Command, args []string) {
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

	fmt.Println(string(post.Content))
}
