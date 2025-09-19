package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a blog post",
	Long:  "Delete a blog post by its ID. This action cannot be undone.",
	Args:  cobra.ExactArgs(1),
	Run:   deletePost,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deletePost(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Are you sure you want to delete post ID %d? (y/N): ", id)
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Delete cancelled")
		return
	}

	if err := postService.DeletePost(id); err != nil {
		fmt.Println(err)
		return
	}
}
