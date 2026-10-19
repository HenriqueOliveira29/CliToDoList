/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/HenriqueOliveira29/cliToDoList/internal/api"
)

var port *string
// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a API server",
	Long: `Start a API server`,
	Run: func(cmd *cobra.Command, args []string) {
		server := api.NewServer(todoService)
		log.Fatal(server.Start(*port))
	},
}

func init() {
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to listen on")
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
