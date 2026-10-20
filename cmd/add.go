/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the list",
	Long: `Add a task to the list.
	Write the title on the task that you want to create`,
	Run: func(cmd *cobra.Command, args []string) {
	  id, err := svc.Add(cmd.Context(), strings.Join(args, " "))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Task added with ID: %d\n", id)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
