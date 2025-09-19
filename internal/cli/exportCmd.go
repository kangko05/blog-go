package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export <id> <filepath>",
	Short: "Export a blog post to a markdown file",
	Long:  "Export an existing blog post to a markdown file.",
	Args:  cobra.ExactArgs(2),
	Run:   exportPost,
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func exportPost(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	outpath := args[1]

	post, err := postService.GetPost(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := os.WriteFile(outpath, []byte(post.Content), 0644); err != nil {
		fmt.Println(err)
		return
	}
}
