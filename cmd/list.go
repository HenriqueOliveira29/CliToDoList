/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all task that you have",
	Long: `Return all the tasks that you have saved on to do`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := todoService.GetAll()
		if err != nil {
			fmt.Println(err)
		}
		for _, element := range tasks {
			fmt.Println("tarefa com o id %d com o titulo %s", element.Id, element.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
