package cmd

import (
	"log"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/HenriqueOliveira29/cliToDoList/internal/api"
)

var (
	addr     string
	noDocker bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Postgres (via Docker) and the API server",

	// Replaces the root hook: start Docker first, then connect to the database.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !noDocker {
			if err := composeUp(cmd); err != nil {
				return err
			}
		}
		return initApp(cmd)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		log.Printf("listening on %s", addr)
		return api.NewServer(svc).Start(ctx, addr)
	},
}

func init() {
	serveCmd.Flags().StringVar(&addr, "addr", ":8080", "address to listen on")
	serveCmd.Flags().BoolVar(&noDocker, "no-docker", false, "skip starting the Postgres container")
	rootCmd.AddCommand(serveCmd)
}